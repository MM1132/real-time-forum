package pages

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn          *websocket.Conn
	ConnectedUser ConnectedUser
}

type ConnectedUser struct {
	UserID   int    `json:"id"`
	Nickname string `json:"nickname"`
	Image    string `json:"image"`
}

type OnlineUsers struct {
	ConnectionList []Connection
}

// Get all connected user objects
// Useful for sending all the online user's information to the client
func (self *OnlineUsers) GetConnectedUsers(recentUsers []ConnectedUser) []ConnectedUser {
	connectedUsers := []ConnectedUser{}
	for _, connection := range self.ConnectionList {
		for _, recentUser := range recentUsers {
			if recentUser.UserID == connection.ConnectedUser.UserID {
				goto End
			}
		}
		connectedUsers = append(connectedUsers, connection.ConnectedUser)
	End:
	}
	return connectedUsers
}

func (self *OnlineUsers) GetRecentUsers(env Chat, userID int) ([]ConnectedUser, error) {
	// Once the user is connected, automatically send the information about all online users to everyone
	recentUsers := []ConnectedUser{}
	recentUserIDs, err := env.Message.GetRecentIDs(userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, id := range recentUserIDs {
		user, err := env.Users.Get(id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		connectedUser := ConnectedUser{
			UserID:   user.UserID,
			Nickname: user.NickName,
			Image:    user.Image,
		}

		recentUsers = append(recentUsers, connectedUser)
	}

	return recentUsers, nil
}

// Add a new connection to the list of connections
func (self *OnlineUsers) Add(conn *websocket.Conn, userID int, nickname, image string) {
	newConnectedUser := ConnectedUser{UserID: userID, Nickname: nickname, Image: image}
	newConnection := Connection{Conn: conn, ConnectedUser: newConnectedUser}
	self.ConnectionList = append(self.ConnectionList, newConnection)
	log.Println("Adding", nickname, "Length", len(self.ConnectionList))
}

// Remove the user with the given ID from the list of online users
func (self *OnlineUsers) Remove(UserID int) {
	for i, v := range self.ConnectionList {
		if v.ConnectedUser.UserID == UserID {
			// Remove the user here
			fmt.Println("Removing", v.ConnectedUser.Nickname)
			self.ConnectionList = append(self.ConnectionList[:i], self.ConnectionList[i+1:]...)
			break
		}
	}
}

// Get a connection by the UserID from the list
func (self *OnlineUsers) GetById(userID int) (Connection, error) {
	for _, v := range self.ConnectionList {
		if v.ConnectedUser.UserID == userID {
			return v, nil
		}
	}
	return Connection{}, errors.New("The connection with the id does not exist")
}

func (self *OnlineUsers) SendMessage(toID int, message interface{}) error {
	// Check if the connection with the userID exists
	connection, err := self.GetById(toID)
	if err != nil {
		return err
	}

	// Send the message through the connection
	if err = connection.Conn.WriteJSON(message); err != nil {
		return err
	}

	// Everything passed successfully without any errors to be found over here
	return nil
}

// Send the message to all users online
func (self *OnlineUsers) Broadcast(messageType string, message interface{}) error {
	// Create the message
	newMessage := SendFormat{
		MessageType: messageType,
		Content:     message,
	}

	for _, v := range self.ConnectionList {
		if err := v.Conn.WriteJSON(newMessage); err != nil {
			return err
		}
	}
	return nil
}
