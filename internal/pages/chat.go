package pages

import (
	"forum/internal/forumEnv"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Chat struct {
	forumEnv.Env
}

// type chatData struct {
// 	forumEnv.GenericData
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (env Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	connection.Close()

	data := forumEnv.GenericData{}
	data.InitData(env.Env, r)
}
