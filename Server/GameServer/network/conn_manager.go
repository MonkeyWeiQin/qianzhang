package network

import (
	"battle_rabbit/iface"
	"sync"
)

type (
	ConnManager struct {
		connections map[string]iface.IAgent
		lock        sync.RWMutex
	}
)

func NewConnManager() *ConnManager {
	mgr := &ConnManager{
		connections: make(map[string]iface.IAgent),
	}
	return mgr
}

//添加链接
func (mgr *ConnManager) Add(agt iface.IAgent) {
	mgr.lock.Lock()
	mgr.connections[agt.GetSession().GetSessionId()] = agt
	mgr.lock.Unlock()

	//mgr.connections.Store(agt.Session.GetSessionId(),agt)
	//mgr.count ++
	// TODO 通知其他网关 有玩家登陆了,并清理已经存在的 登录
}

//删除连接
func (mgr *ConnManager) Remove(sessionId string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	if _, ok := mgr.connections[sessionId]; ok {
		delete(mgr.connections, sessionId)
	}

	//mgr.connections.Delete(sessionId)
	//mgr.count --
}

//利用ConnID获取链接
func (mgr *ConnManager) Get(sessionId string) iface.IAgent {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	return mgr.connections[sessionId]

	//agt,ok := mgr.connections.Load(sessionId)
	//if ok {
	//	return  agt.(*Agent)
	//}
	//return nil
}

// 获取当前服务器总连接数量
func (mgr *ConnManager) Len() int {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	return len(mgr.connections)

	//count := 0
	//mgr.connections.Range(func(key, value interface{}) bool {
	//	count++
	//	return true
	//})
	//return count

	//return mgr.count
}

//删除并停止所有链接
func (mgr *ConnManager) ClearConn() {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	//for _, agt := range mgr.Connections {
	//	_ = agt.OnClose()
	//}
	mgr.connections = make(map[string]iface.IAgent)

	//mgr.connections.Range(func(key, value interface{}) bool {
	//	mgr.connections.Delete(key)
	//	return true
	//})
}

//遍历连接
func (mgr *ConnManager) Range(fn func(agt iface.IAgent)) {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	for _, agt := range mgr.connections {
		fn(agt)
	}

	//mgr.connections.Range(func(key, value interface{}) bool {
	//	fn(value.(*Agent))
	//	return true
	//})
}
