package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var listWsConnection = make(map[string]*websocket.Conn)

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsFood(c *gin.Context) {
	wsGetFood(c.Writer, c.Request)
}

func wsGetFood(w http.ResponseWriter, r *http.Request) {

	ip := r.Header.Get("ip")

	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade", err.Error())
		return
	}

	// if ip == "172.16.41.4" || ip == "172.16.41.3" {
	listWsConnection[ip] = conn
	// }
}
