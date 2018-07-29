package impl

import (
	"github.com/gorilla/websocket"
	"sync"
	"errors"
)

type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte

	mutex sync.Mutex
	isClosed bool
}

func InitConn(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:wsConn,
		inChan:make(chan []byte, 1000),
		outChan:make(chan []byte, 1000),
		closeChan:make(chan byte, 1),
	}
	//启动读协程
	go conn.readLoop()

	//启动写协程
	go conn.writeLoop()

	return
}

//API
func (conn *Connection)ReadMessage() (data []byte, err error) {
	select {
	case data = <- conn.inChan:
	case <- conn.closeChan:
		err = errors.New("ws connection is closed")
	}

	return
}

func (conn *Connection)WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <- conn.closeChan:
		err = errors.New("ws connection is closed")
	}

	return 
}

func (conn *Connection) Close() {
	//线程安全,可重入的Close
	conn.wsConn.Close()

	//关闭closeChan, 需确保只关闭1次
	conn.mutex.Lock()
	if !conn.isClosed{
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

//内部实现
func (conn *Connection)readLoop() {
	var(
		data []byte
		err error
	)
	for{
		if _,data,err = conn.wsConn.ReadMessage();err!=nil{
			goto ERR
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			//closeChan关闭了
			goto ERR
		}
		//read ok
	}
ERR:
	conn.Close()
}

func (conn *Connection)writeLoop(){
	var(
		data []byte
		err error
	)
	for{
		select {
		case data = <-conn.outChan:
		case <- conn.closeChan:
			//closeChan关闭了
			goto ERR
		}

		if conn.wsConn.WriteMessage(websocket.TextMessage,data);err!=nil{
			goto ERR
		}
		//write ok
	}
ERR:
	conn.Close()
}