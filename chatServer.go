package main

import (
	"github.com/google/uuid"
)

const PubSubGeneralChannel = "general"

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Rooms]bool
	//users      []models.User
	//roomRepository repository.RoomRepository
	//userRepository models.UserRepository
}

// NewWebsocketServer creates a new WsServer type
// func NewWebsocketServer(roomRepository models.RoomRepository, userRepository models.UserRepository) *WsServer {
func NewWebsocketServer() *WsServer {
	wsServer := &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[*Rooms]bool),
		broadcast:  make(chan []byte),
		//roomRepository: roomRepository,
		// userRepository: userRepository,
	}

	// Add users from database to server
	// wsServer.users = userRepository.GetAllUsers()

	return wsServer
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	// go server.listenPubSubChannel()
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	// server.userRepository.AddUser(client)
	// server.publishClientJoined(client)
	server.notifyClientJoined(client)
	server.listOnlineClients(client)
	server.clients[client] = true
	// server.users = append(server.users, message.Sender)
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		server.notifyClientLeft(client)

		// for i, user := range server.users {
		// 	if user.GetId() == message.Sender.GetId() {
		// 		server.users[i] = server.users[len(server.users)-1]
		// 		server.users = server.users[:len(server.users)-1]
		// 	}
		// }

		// server.userRepository.RemoveUser(client)
		// server.publishClientLeft(client)
	}
}

func (server *WsServer) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: *client,
	}

	server.broadcastToClients(message.encode())
}

func (server *WsServer) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: *client,
	}

	server.broadcastToClients(message.encode())
}

func (server *WsServer) listOnlineClients(client *Client) {
	for existingClient := range server.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: *existingClient,
		}
		client.send <- message.encode()
	}
	// for _, user := range server.users {
	// 	message := &Message{
	// 		Action: UserJoinedAction,
	// 		Sender: user,
	// 	}
	// 	client.send <- message.encode()
	// }

}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) findRoomByID(ID string) *Rooms {
	var foundRoom *Rooms
	parse, _ := uuid.Parse(ID)
	for room := range server.rooms {
		if room.ID == parse {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(name string, private bool) *Rooms {
	room := NewRoom(name, private)
	// server.roomRepository.AddRoom(room)
	go room.RunRoom()
	server.rooms[room] = true

	return room
}

func (server *WsServer) findClientByID(ID uuid.UUID) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.ID == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}

/*func (server *WsServer) runRoomFromRepository(name string) *Rooms {
	var room *Rooms
	dbRoom := server.roomRepository.FindRoomByName(name)
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetId())

		go room.RunRoom()
		server.rooms[room] = true
	}

	return room
}
*/

func (server *WsServer) findRoomByName(name string) *Rooms {
	var foundRoom *Rooms
	for room := range server.rooms {
		if room.Name == name {
			foundRoom = room
			break
		}
	}

	// NEW: if there is no room, try to create it from the repo
	//if foundRoom == nil {
	// Try to run the room from the repository, if it is found.
	//foundRoom = server.runRoomFromRepository(name)
	//}

	return foundRoom
}

// func (server *WsServer) publishClientJoined(client *Client) {

// 	message :=
// 	 &Message{
// 		Action: UserJoinedAction,
// 		Sender: client,
// 	}

// if err := config.Redis.Publish(PubSubGeneralChannel, message.encode()).Err(); err != nil {
// 	log.Println(err)
// }
// }

// func (server *WsServer) publishClientLeft(client *Client) {

// 	message := &Message{
// 		Action: UserLeftAction,
// 		Sender: client,
// 	}

// if err := config.Redis.Publish(PubSubGeneralChannel, message.encode()).Err(); err != nil {
// 	log.Println(err)
// }
// }

// func (server *WsServer) listenPubSubChannel() {

// 	// pubsub := config.Redis.Subscribe(PubSubGeneralChannel)
// 	// ch := pubsub.Channel()
// 	for msg := range ch {

// 		var message Message
// 		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
// 			log.Printf("Error on unmarshal JSON message %s", err)
// 			return
// 		}

// 		switch message.Action {
// 		case UserJoinedAction:
// 			server.handleUserJoined(message)
// 		case UserLeftAction:
// 			server.handleUserLeft(message)
// 		case JoinRoomPrivateAction:
// 			server.handleUserJoinPrivate(message)
// 		}
// 	}
// }

func (server *WsServer) handleUserJoinPrivate(message Message) {
	targetClient := server.findClientByID(message.Sender.ID)
	if targetClient != nil {
		targetClient.joinRoom(message.Target.Name, &message.Sender, nil)
	}
}

/*func (server *WsServer) FindUserById(ID string) models.User {
	var foundUser models.User
	for _, client := range server.users {
		if client.GetId() == ID {
			foundUser = client
			break
		}
	}
	return foundUser
}*/

/*func (server *WsServer) handleUserJoined(message Message) {
	// Add the user to the slice
	server.users = append(server.users, *message.Sender)
	server.broadcastToClients(message.encode())
}*/

/*func (server *WsServer) handleUserLeft(message Message) {
	// Remove the user from the slice
	for i, user := range server.users {
		if user.GetId() == message.Sender.GetId() {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
		}
	}
	server.broadcastToClients(message.encode())
}*/
