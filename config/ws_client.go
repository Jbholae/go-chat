package config

import (
	"encoding/json"
	"golang-chat/constants"
	"golang-chat/models"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// Client represents the websocket client at the server
type Client struct {
	models.Client
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	Send     chan []byte
	rooms    map[*Rooms]bool
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	return &Client{
		Client: models.Client{
			ID:   uuid.New(),
			Name: name,
		},
		conn:     conn,
		wsServer: wsServer,
		Send:     make(chan []byte, 256),
		rooms:    make(map[*Rooms]bool),
	}

}

// This goRoutine read's message
func (client *Client) readPump() {
	defer func() {
		println("stopped runnign")
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}

}

// This goRoutine Send's message
func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.wsServer.Unregister <- client
	for room := range client.rooms {
		room.Unregister <- client
	}
	close(client.Send)
	client.conn.Close()
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, name[0])

	go client.readPump()
	go client.writePump()

	wsServer.Register <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message models.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}
	// Attach the client object as the sender of the messsage.
	message.Sender = client.Client

	switch message.Action {
	case constants.SendMessageAction:
		// The Send-message action, this will Send messages to a specific room now.
		// Which room wil depend on the message Target
		room := message.Target
		// Use the ChatServer method to find the room, and if found, broadcast!
		if room := client.wsServer.FindRoomByName(room.Name); room != nil {
			room.Broadcast <- &message
		}

	case constants.JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case constants.LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case constants.JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message)
	}

}

func (client *Client) handleJoinRoomMessage(message models.Message) {
	roomName := message.Message
	room := client.wsServer.FindRoomByName(roomName)
	if room == nil {
		room = client.wsServer.CreateRoom(roomName, false)
	}

	client.JoinRoom(roomName, client, room)
}

func (client *Client) handleLeaveRoomMessage(message models.Message) {
	room := client.wsServer.FindRoomByName(message.Message)
	if room == nil {
		return
	}

	if ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.Unregister <- client
}

func (client *Client) handleJoinRoomPrivateMessage(message models.Message) {
	target := client.wsServer.FindClientByID(message.Sender.ID)

	if target == nil {
		return
	}

	// create unique room name combined to the two IDs
	roomName := message.Message + client.ID.String()

	// joinedRoom := client.joinRoom(roomName, target)

	// if joinedRoom != nil {
	// 	client.inviteTargetUser(target, joinedRoom)
	// }

	client.JoinRoom(roomName, target, nil)
	target.JoinRoom(roomName, client, nil)

}

func (client *Client) JoinRoom(roomName string, sender *Client, room *Rooms) *Rooms {
	if room == nil {
		room := client.wsServer.FindRoomByName(roomName)
		if room == nil {
			room = client.wsServer.CreateRoom(roomName, sender != nil)
		}
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return nil
	}

	if !client.isInRoom(room) {

		client.rooms[room] = true
		room.Register <- client

		client.notifyRoomJoined(*room, *sender)
	}
	return room
}

func (client *Client) isInRoom(room *Rooms) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}

func (client *Client) notifyRoomJoined(room Rooms, sender Client) {
	message := models.Message{
		Action: constants.RoomJoinedAction,
		Target: room.Room,
		Sender: sender.Client,
	}

	client.Send <- message.Encode()
}

func (client *Client) GetName() string {
	return client.Name
}

func (client *Client) GetId() string {
	return client.ID.String()
}

// func (client *Client) inviteTargetUser(target models.User, room *Room) {
// 	inviteMessage := &Message{
// 		Action:  JoinRoomPrivateAction,
// 		Message: target.GetId(),
// 		Target:  room,
// 		Sender:  client,
// 	}
// 	if err := config.Redis.Publish(PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
// 		log.Println(err)
// 	}
// }
