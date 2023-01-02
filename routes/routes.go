package routes

import (
	"fmt"
	"go.uber.org/fx"
	"golang-chat/config"
	"golang-chat/controller"
	"net/http"
)

var Module = fx.Options(
	fx.Provide(NewRoutes),
)

type Routes struct {
	userController controller.UserController
	roomController controller.RoomController
	wsServer       *config.WsServer
}

func NewRoutes(
	userController controller.UserController,
	roomController controller.RoomController,
	wsServer *config.WsServer,
) Routes {
	return Routes{
		userController: userController,
		roomController: roomController,
		wsServer:       wsServer,
	}
}

func (r Routes) InitializeRouter() {
	http.HandleFunc("/api/user", r.userController.CreateUser)

	http.HandleFunc("/api/rooms", r.roomController.GetRooms)

	http.HandleFunc("/api/nothing", controller.GetNothing)

	http.HandleFunc("/chat", func(w http.ResponseWriter, req *http.Request) {
		config.ServeWs(r.wsServer, w, req)
		fmt.Println("running")
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/front", fs)

}
