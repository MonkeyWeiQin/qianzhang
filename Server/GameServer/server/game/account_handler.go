package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/network"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/request"
	"battle_rabbit/protocol/response"
	v1 "battle_rabbit/server/login/api/v1"
	"battle_rabbit/service/log"
	"battle_rabbit/service/redisDB"
	"battle_rabbit/utils"
	"battle_rabbit/utils/http_client"
	"battle_rabbit/utils/serialize"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

const (
	HTTPCheckTokenUri = "/v1/checkToken"
)

func (g *Game) CheckTokenTest(sess iface.ISession, msg *codec.Message) {
	req := new(request.BindSessionReq)
	err := jsoniter.Unmarshal(msg.Data, req)
	if err != nil {
		log.Error("登录解析参数错误: ", err)
		sess.Send(protocol.ErrCode(msg.Id,define.MsgCode400))
		return
	}

	user := new(model.UserModel)
	devId := req.Url
	err = model.GetUserCollection().FindOne(nil, bson.M{"devId": devId}, user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			user, err = v1.CreateUsersByDevId(nil, devId)
			if err != nil {
				log.Error(err)
				sess.Send(protocol.Err(msg.Id))
				return
			}
		} else {
			log.Error(err)
			sess.Send(protocol.Err(msg.Id))
			return
		}
	}
	// 单点登录检测
	err = g.singleLoginCheck(sess)
	if err != nil {
		sess.Send(protocol.ErrCode(msg.Id,define.MsgCode402))
		return
	}
	// TODO  检查用户邮件数据并推送

	// 绑定用户
	err = sess.Bind(user.Uid, g.NodeId)
	if err != nil {
		log.Error("绑定session 失败!! err : ", err)
		sess.Send(protocol.ErrCode(msg.Id,define.MsgCode402))
		return
	}
	log.Debug("用户绑定成功!!! uid : %d,sessionId:%s", sess.GetUid(), sess.GetSessionId())
	sess.Send(protocol.Success(msg.Id))
}

// 用户绑定session
func (g *Game) CheckToken(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id,code))
			time.AfterFunc(time.Second * 2, func() {
				_ = sess.Close() // 关闭连接
			})
		}
	}()

	req := new(request.BindSessionReq)
	err := jsoniter.Unmarshal(msg.Data, req)
	if err != nil {
		log.Error("登录解析参数错误: ", err)
		code = define.MsgCode400
		return
	}
	// 1 使用http校验token
	client, err := http_client.NewClientBuilder().Build()
	if err != nil {
		log.Error("登录解析参数错误: ", err)
		code = define.MsgCode402
		return
	}
	client.SetHeader(map[string]string{"x-token": req.Token})
	httpResp := client.Get(req.Url + HTTPCheckTokenUri)
	if httpResp.StatusCode() != http.StatusOK || httpResp.Error() != nil {
		log.Error("登录校验失败: http status: %d , error: ", httpResp.StatusCode(), httpResp.Error())
		code = define.MsgCode402
		return
	}

	//2 登录成功, 解析用户ID,绑定session
	httpRespBody := new(response.ValidateUserResp)
	err = jsoniter.Unmarshal(httpResp.Content(), httpRespBody)
	if err != nil || httpRespBody.Code != http.StatusOK { //resp.Code 自定义的code 与http code区分
		log.Error("登录解析远程服务器响应失败: ", err, httpRespBody.Code)
		code = define.MsgCode402
		return
	}

	// 处理的绑定成功后的事件
	// 单点登录检测
	err = g.singleLoginCheck(sess)
	if err != nil {
		code = define.MsgCode402
		return
	}
	// TODO  检查用户邮件数据并推送

	// 绑定用户
	uid, ok := httpRespBody.Data.(float64)
	if !ok {
		code = define.MsgCode500
		return
	}
	err = sess.Bind( int(uid), g.NodeId)
	if err != nil {
		log.Error("绑定session 失败!! err : ", err)
		code = define.MsgCode402
		return
	}
	log.Debug("用户绑定成功!!! uid : %d,sessionId:%s", sess.GetUid(), sess.GetSessionId())
	// 加载用数据
	//if _,ok := g.GetPlayer(int(uid));!ok {
	//	player := NewPlayer(int(uid))
	//	err = player.LoadDataFromDB()
	//	if err != nil {
	//		return reply.ErrCode(args.Msg.GetMsgId(), define.MsgCode500)
	//	}
	//	g.AddPlayer(int(uid),player)
	//}

	//

	// 返回登录成功的数据
	sess.Send(protocol.Success(msg.Id))
}

// 单点登录,断开前面的连接
func (g *Game) singleLoginCheck(session iface.ISession) (err error) {
	data, err := redis.Bytes(redisDB.Client.GET(fmt.Sprintf("session_%d", session.GetUid())))
	if err != nil {
		if err == redis.ErrNil {
			return nil
		}
		return
	}
	oldSess := new(network.Session)
	err = serialize.Decode(data, oldSess)
	if err != nil {
		return
	}
	connection := oldSess.IsConnect()
	if connection {
		oldSess.Send(protocol.MakePushMsg(define.PushMsgId1001, nil))
		time.AfterFunc(time.Second * 3, func() {
			oldSess.Close()
		})
	}
	return
}

// 获取用户信息
func (g *Game) GetPlayerInfo(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		// 重新登录
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	// 每次登录会拉取一次用户信息,在这里做登录天数统计
	if utils.GetDay(utils.GetNowTimeStamp()) != utils.GetDay(player.Account.LoginTime) && player.Task != nil { // 不是同一天
		UpdateTask(player.Task,sess,[]*model.TaskProgress{{
				Type:      define.TaskContinuousLogin,
				Num:       1,
			}})
	}
	model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"loginTime": utils.GetNowTimeStamp()})
	sess.Send(protocol.SuccessData(msg.Id,player.Account))
}

// 更新新手引导步骤
func (g *Game) UpdateGuideStep(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		// 重新登录
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	req := make(map[string]string)
	err := jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		sess.Send(protocol.ErrCode(msg.Id,define.MsgCode400))
		return
	}
	player.Account.GuideStep = req["guideStep"]

	_, err = model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"guideStep": req["guideStep"]})
	if err != nil {
		log.Error(err)
	}
	sess.Send(protocol.Success(msg.Id))
}


