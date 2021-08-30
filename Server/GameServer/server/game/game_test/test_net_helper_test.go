package game_test

import (
	"battle_rabbit/codec"
	"battle_rabbit/service/log"
	"fmt"
	"net"
)

const (
	url = "127.0.0.1:19010"
)

var (
	RespChan = make(chan *codec.Message,10)
	ReqChan = make(chan *codec.Message)
)

type RespStatus struct {
	Code int `json:"code"`
}


type Cli struct {
	conn *net.TCPConn
	stop bool
}


func NewTcpCli() *Cli {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		panic(err)
	}

	go Reader(conn)
	go Writer(conn)

	return &Cli{
		conn: conn.(*net.TCPConn),
		stop: false,
	}
}

func Writer(conn net.Conn) {
	for {
		msg ,ok := <-ReqChan
		if !ok {
			return
		}
		// --------------
		log.Debug(" \r\n 客户端发送数据: id : %d,  data : %s ",msg.Id, string(msg.Data))
		// -----------
		err := msg.WriterPack(conn)
		if err != nil {
			fmt.Println("WriterPack===", err)
			return
		}
	}
	fmt.Println(" 写  协程退出 !!")
}

func Reader(conn net.Conn) {
	for {
		msg := new(codec.Message)
		err := msg.ReadPack(conn)
		if err != nil {
			fmt.Println("ReadPack===", err)
			return
		}
		RespChan <- msg
	}
	fmt.Println(" 读协程退出 !!")
}

func (c *Cli) Stop()  {
	c.stop = true
	c.conn.Close()
	close(ReqChan)
	close(RespChan)
}

