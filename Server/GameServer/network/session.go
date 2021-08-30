package network

import (
	"battle_rabbit/iface"
	"battle_rabbit/service/log"
	"battle_rabbit/service/redisDB"
	"battle_rabbit/utils/serialize"
	"battle_rabbit/utils/xid"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
)

type (
	Session struct {
		NodeId string
		agent  *Agent
		SessId string
		Uid    int                    // 用户ID
		Cache  map[string]interface{} // 缓存少量数据
		mtx    sync.RWMutex
	}
)

func NewSession(agent *Agent) *Session {
	return &Session{
		agent:  agent,
		SessId: xid.New().String(),
	}
}

func (sess *Session) SetNodeId(nodeId string) error {
	sess.NodeId = nodeId
	return sess.agent.Storage(sess)
}

func (sess *Session) GetNodeId() string {
	return sess.NodeId
}

func (sess *Session) GetSessionId() string {
	return sess.SessId
}

func (sess *Session) SetSessionId(sessId string) error {
	sess.SessId = sessId
	return sess.agent.Storage(sess)
}

func (sess *Session) GetUid() int {
	return sess.Uid
}

func (sess *Session) Bind(uid int,nodeId string) (err error) {
	defer func() {
		if e := recover(); e != nil || err != nil {
			log.Error(err, e)
			err = nil
		}
	}()

	sess.Uid = uid
	sess.NodeId = nodeId
	// 绑定玩家之前的事件触发
	// 查询 这个uid以前的缓存数据
	var data []byte
	data, err = sess.agent.Query(uid)
	if err != nil {
		return
	}
	// 反序列化数据
	if data != nil {
		old := new(Session)
		err = serialize.Decode(data, old)
		if err != nil {
			return err
		}
		if len(old.Cache) > 0 {
			for k, v := range old.Cache {
				sess.Cache[k] = v
			}
		}
	}
	// 绑定session
	err = sess.agent.Storage(sess)
	if err != nil {
		return err
	}

	log.Debug("bind :::: , uid : %d \r\n ",  sess.GetUid())
	return nil

}

func (sess *Session) Close() error {
	return sess.agent.OnClose()
}

func (sess *Session) IsConnect() bool {
	return !sess.agent.stopped
}

func (sess *Session) Send(msg iface.IMessage) {
	sess.agent.writerChan <- msg
	return
}

func (sess *Session) GetUserCacheAll() map[string]interface{} {
	return sess.Cache
}

func (sess *Session) GetUserCacheByKey(k string) interface{} {
	sess.mtx.RLock()
	i := sess.Cache[k]
	sess.mtx.RUnlock()
	return i
}

func (sess *Session) SetUserCacheByKV(k string, v interface{}) error {
	sess.mtx.Lock()
	sess.Cache[k] = v
	sess.mtx.Unlock()
	return sess.agent.Storage(sess)
}

func (sess *Session) SetUserCacheByMap(data map[string]interface{}) error {
	sess.Cache = data
	return sess.agent.Storage(sess)
}

type (
	SessionStorage struct{}
	SessionLearner struct{}
)

//当连接建立  tcp协议握手成功
func (*SessionLearner) Connect(sess iface.ISession) {
	log.Info("客户端建立连接, session id is : %s", sess.GetSessionId())
}

//当连接关闭	或者客户端主动发送close命令
func (*SessionLearner) DisConnect(sess iface.ISession) {
	log.Info("客户端断开连接, session id is: %s ", sess.GetSessionId())
}

// 存储session到缓存
func (*SessionStorage) Storage(session iface.ISession) (err error) {
	log.Debug("调用存储 session")
	// 1 序列化
	sessStr, err := serialize.Encode(session)
	if err != nil {
		return
	}
	_, err = redisDB.Client.SET(fmt.Sprintf("session_%d", session.GetUid()), sessStr, "EX", 60*60*24)
	return
}

// 删除缓存的session
func (*SessionStorage) Delete(session iface.ISession) (err error) {
	log.Debug("调用删除session")
	_, err = redisDB.Client.DEL(fmt.Sprintf("session_%d", session.GetUid()))
	return
}

// 查询缓存的session
func (*SessionStorage) Query(uid int) (data []byte, err error) {
	log.Debug("查询session 信息")
	data, err = redis.Bytes(redisDB.Client.GET(fmt.Sprintf("session_%d", uid)))
	if err == redis.ErrNil {
		return nil, nil
	}
	return
}
