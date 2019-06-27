package websocket

import (
	"net/http"
	"strconv"

	"alopex/app"

	"github.com/gorilla/websocket"
)

type WebSocketController struct{}

var upgrader websocket.Upgrader

func init() {
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	app.CJoin("websocket", WebSocketController{})
}

func (ctrl WebSocketController) Todo(h *app.Http) {
	conn, err := upgrader.Upgrade(h.Rep, h.Req, nil)
	if err != nil {
		if app.IsDeveloper {
			app.Dump("red", "WebSocket连接异常，"+err.Error())
		}
	} else if conn != nil {
		conn.SetCloseHandler(func(code int, text string) error {
			app.Dump("yellow", strconv.Itoa(code)+"-"+text)
			return nil
		})
		defer conn.Close()
		for {
			// 接收客户端消息
			_, bs, err := conn.ReadMessage()
			if (err != nil) || (len(bs) < 1) {
				continue
			}
			result := make(map[string]interface{})
			switch string(bs) {
			case "refresh":
				result = ctrl.GetData()
			}
			// 返回数据给客户端
			conn.WriteJSON(result)
		}
	}
}

func (ctrl WebSocketController) GetData() map[string]interface{} {
	// do more ....
	return map[string]interface{}{"sdfs": 2342, "sss": 2342, "sdsfs": 2342, "sdafs": 2342, "sdaafs": 2342, "sdfsd": 231231}
}
