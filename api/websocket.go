package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	news   = make(map[string]chan interface{})
	client = make(map[string]*websocket.Conn)
	mux    sync.Mutex
)

func (s *Server) subscribe(ctx *gin.Context) {
	id, _ := ctx.GetQuery("id")

	client, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = client.WriteMessage(websocket.TextMessage, []byte("websocket connected"))
	if err != nil {
		logrus.Errorln("websocket write msg error", err)
		return
	}

	s.mu.Lock()
	s.clients[id] = client
	s.mu.Unlock()
	//for {
	//	// 每隔两秒给前端推送一句消息“hello, WebSocket”
	//	err := client.WriteMessage(websocket.TextMessage, []byte("hello, WebSocket"))
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	time.Sleep(time.Second * 2)
	//}
	return
}

func GetPushNews(ctx *gin.Context) {
	id := ctx.Query("userId")

	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Info("websocket connected")

	WsHandler(ctx.Writer, ctx.Request, id)
}

func DeleteClient(ctx *gin.Context) {
	id := ctx.Param("id")

	conn, exist := getClient(id)
	if exist {
		conn.Close()
		deleteClient(id)
	} else {
		ctx.JSON(http.StatusOK, msgResponse("client not found"))
	}

	_, exist = getChannel(id)
	if exist {
		deleteChannel(id)
	}
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request, id string) {
	var (
		conn  *websocket.Conn
		err   error
		exist bool
	)

	pingTicker := time.NewTicker(time.Second * 10)
	conn, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	addClient(id, conn)

	m, exist := getChannel(id)
	if !exist {
		m = make(chan interface{})
		addChannel(id, m)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		deleteClient(id)
		log.Println(code)
		return nil
	})

	for {
		select {
		case content, _ := <-m:
			err = conn.WriteJSON(content)
			if err != nil {
				log.Println(err)
				conn.Close()
				deleteClient(id)
				return
			}
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
			err = conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("send ping err:", err)
				conn.Close()
				deleteClient(id)
				return
			}
		}
	}
}

func addClient(id string, conn *websocket.Conn) {
	mux.Lock()
	client[id] = conn
	mux.Unlock()
}

func getClient(id string) (conn *websocket.Conn, exist bool) {
	mux.Lock()
	conn, exist = client[id]
	mux.Unlock()
	return
}

func deleteClient(id string) {
	mux.Lock()
	delete(client, id)
	log.Println(id + "websocket退出")
	mux.Unlock()
}

func addChannel(id string, m chan interface{}) {
	mux.Lock()
	news[id] = m
	mux.Unlock()
}

func getChannel(id string) (m chan interface{}, exist bool) {
	mux.Lock()
	m, exist = news[id]
	mux.Unlock()
	return
}

func deleteChannel(id string) {
	mux.Lock()
	if m, ok := news[id]; ok {
		close(m)
		delete(news, id)
	}
	mux.Unlock()
}
