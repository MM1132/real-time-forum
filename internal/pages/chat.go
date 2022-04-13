package pages

import (
	"encoding/json"
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Chat struct {
	forumEnv.Env
	OnlineUsers *OnlineUsers
}

type ChatData struct {
	forumEnv.GenericData
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// MESSAGE FORMAT STUFF
type ReceiveFormat struct {
	MessageType string          `json:"messageType"`
	Content     json.RawMessage `json:"content"`
}

type SendFormat struct {
	MessageType string      `json:"messageType"`
	Content     interface{} `json:"content"`
}

func (env Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := ChatData{}
	data.InitData(env.Env, r)

	// Check if the user is online
	if data.User.UserID == 0 {
		return
	}

	// Create and add the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	// If the connection with the ID already exists, remove it and add the new one
	env.OnlineUsers.Remove(data.User.UserID)
	env.OnlineUsers.Add(conn, data.User.UserID, data.User.NickName, data.User.Image)

	/* recentUsers, err := env.Message.GetRecentIDs(data.User.UserID)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%-v", recentUsers) */

	// Create a go routine
	go func(conn *websocket.Conn) {

		err := env.UpdateUserLists(data)
		if err != nil {
			log.Println(err)
			return
		}
		// ########

		running := true
		for running {
			if data.User.UserID == 0 {
				break
			}

			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(0, err)
				break
			}

			// Get the message type
			messageFormat := ReceiveFormat{}
			err = json.Unmarshal(message, &messageFormat)
			if err != nil {
				log.Println(1, err)
				break
			}

			// Do things
			switch messageFormat.MessageType {
			case "message-client-server": // Client is sending a message to the server
				// Decode the message
				decodedMsg := forumDB.Message{}
				err := json.Unmarshal(messageFormat.Content, &decodedMsg)
				if err != nil {
					log.Println(2, err)
					break
				}

				// If the user is trying to send a message to themselves, don't allow it
				if decodedMsg.ToID == data.User.UserID {
					log.Println("User tried sending a message to themselves")
					break
				}

				// If the "to" ID exists
				if _, err := env.Users.Get(decodedMsg.ToID); err != nil {
					log.Println(3, err)
					break
				}

				// Save the message to DB
				newMessageID, err := env.Message.Insert(data.User.UserID, decodedMsg.ToID, decodedMsg.Body)
				if err != nil {
					log.Println(4, err)
					break
				}

				// Once the message is saved, send the new recent list to both the users
				err = env.UpdateRecentList(data.User.UserID)
				if err != nil {
					log.Println(err)
				}
				// The second user
				err = env.UpdateRecentList(decodedMsg.ToID)
				if err != nil {
					log.Println(err)
				}

				// Get the saved message from the DB
				newMessage, err := env.Message.GetByID(newMessageID)
				if err != nil {
					log.Println(5, err)
					break
				}

				// Get the message ready for the sending
				messageFormat := SendFormat{
					MessageType: "message-server-client",
				}
				messageFormat.Content = newMessage

				// Try sending the message to both connections
				errFrom := env.OnlineUsers.SendMessage(newMessage.FromID, messageFormat)
				errTo := env.OnlineUsers.SendMessage(newMessage.ToID, messageFormat)

				// And handle errors for both connections
				if errFrom != nil {
					log.Println(7, err)
				}
				if errTo != nil {
					log.Println(7, err)
				}

			case "request-history": // Client is requesting 10 messages from the server
				// Read the request from the client
				requestHistory := forumDB.RequestHistory{}
				if err = json.Unmarshal(messageFormat.Content, &requestHistory); err != nil {
					log.Println(err)
					break
				}

				// Get the messages
				messages, err := env.Message.GetChat(data.User.UserID, requestHistory.UserID, requestHistory.Index)
				if err != nil {
					log.Println(err)
					break
				}

				// Send the messages back to the client
				sendMessages := SendFormat{
					MessageType: "history",
					Content:     messages,
				}
				if err := env.OnlineUsers.SendMessage(data.User.UserID, sendMessages); err != nil {
					log.Println(err)
				}

			case "logout": // Logging out of the chat
				running = false
			}
		}

		// User disconnected from the chat
		// Remove them from the list of connected users
		env.OnlineUsers.Remove(data.User.UserID)

		// Also, once the user has disconnected, send the updated list of all online users to everyone
		err = env.UpdateUserLists(data)
		if err != nil {
			log.Println(err)
			return
		}

	}(conn)

	if err != nil {
		log.Println(err)
		return
	}
}

func (env Chat) UpdateUserLists(data ChatData) error {
	// Get all the online users for everyone
	onlineUsers := env.OnlineUsers.GetConnectedUsers([]ConnectedUser{})

	// And also, update the recent messages for everyone separately
	for _, v := range env.OnlineUsers.ConnectionList {
		// Here are the recent users for the UserID
		recentUsers, err := env.OnlineUsers.GetRecentUsers(env, v.ConnectedUser.UserID)
		if err != nil {
			log.Println(err)
			return err
		}

		// Get the recent and online uesrs into the sending format
		allUsers := struct {
			RecentUsers []ConnectedUser `json:"recent"`
			OnlineUsers []ConnectedUser `json:"online"`
		}{
			RecentUsers: recentUsers,
			OnlineUsers: onlineUsers,
		}
		sendUsers := SendFormat{MessageType: "update-user-list", Content: allUsers}

		// Send
		err = v.Conn.WriteJSON(sendUsers)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (env Chat) UpdateRecentList(UserID int) error {
	// Get all the online users. Same for everyone
	onlineUsers := env.OnlineUsers.GetConnectedUsers([]ConnectedUser{})

	// Get the recent users for UserID passed
	recentUsers, err := env.OnlineUsers.GetRecentUsers(env, UserID)
	if err != nil {
		return err
	}

	// Get it into the sending format
	allUsers := struct {
		RecentUsers []ConnectedUser `json:"recent"`
		OnlineUsers []ConnectedUser `json:"online"`
	}{
		RecentUsers: recentUsers,
		OnlineUsers: onlineUsers,
	}
	sendUsers := SendFormat{MessageType: "update-user-list", Content: allUsers}

	// Send
	connection, err := env.OnlineUsers.GetById(UserID)
	if err != nil {
		return err
	}
	err = connection.Conn.WriteJSON(sendUsers)
	if err != nil {
		return err
	}

	return nil
}
