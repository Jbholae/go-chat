package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

var ctx = context.Background()

const welcomeMessage = "%s joined the room"

type Rooms struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	Private    bool `json:"private"`
}

// NewRoom creates a new Room
func NewRoom(name string, private bool) *Rooms {
	return &Rooms{
		ID:         uuid.New(),
		Name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		Private:    private,
	}
}

// RunRoom runs our room, accepting various requests
func (room *Rooms) RunRoom() {
	// go room.subscribeToRoomMessages()
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message.encode())
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
		client.send <- message
	}
}

func (room *Rooms) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClientsInRoom(message.encode())
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
