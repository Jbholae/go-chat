package config

import (
	"fmt"
	"golang-chat/constants"
	"golang-chat/models"
	"gorm.io/gorm"
)

const welcomeMessage = "%s joined the room"

type Rooms struct {
	models.Room
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *models.Message
}

// NewRoom creates a new Room
func NewRoom(name string, private bool) *Rooms {
	return &Rooms{
		Room: models.Room{
			Model: gorm.Model{
				ID: 112,
			},
			Name:    name,
			Private: private,
		},
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *models.Message),
	}
}

// RunRoom runs our room, accepting various requests
func (room *Rooms) RunRoom() {
	// go room.subscribeToRoomMessages()
	for {
		select {

		case client := <-room.Register:
			room.registerClientInRoom(client)

		case client := <-room.Unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.Broadcast:
			room.broadcastToClientsInRoom(message.Encode())
		}

	}
}

func (room *Rooms) registerClientInRoom(client *Client) {
	// if !room.Private {
	// by sending the message first the new user won't see his own message.
	room.notifyClientJoined(client)
	// }
	room.clients[client] = true
}

func (room *Rooms) unregisterClientInRoom(client *Client) {
	if ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Rooms) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.Send <- message
	}
}

func (room *Rooms) notifyClientJoined(client *Client) {
	message := &models.Message{
		Action:  constants.SendMessageAction,
		Target:  room.Room,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClientsInRoom(message.Encode())
	// room.publishRoomMessage(message.encode())
}

// func (room *Room) publishRoomMessage(message []byte) {
// 	// err := config.Redis.Publish(room.GetName(), message).Err()

// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func (room *Room) subscribeToRoomMessages() {
// 	pubsub := config.Redis.Subscribe(room.GetName())

// 	ch := pubsub.Channel()

// 	for msg := range ch {
// 		room.broadcastToClientsInRoom([]byte(msg.Payload))
// 	}
// }
