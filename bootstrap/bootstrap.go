package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"golang-chat/config"
	"golang-chat/controller"
	"golang-chat/routes"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var Module = fx.Options(
	config.Module,
	controller.Module,
	routes.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifeCycle fx.Lifecycle,
	db *gorm.DB,
	routes routes.Routes,
) {
	lifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			routes.InitializeRouter()

			go func() {
				log.Println("Starting http server on port::8080")
				err := http.ListenAndServe(":8000", nil)
				if err != nil {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			dbConnection, err := db.DB()
			dbConnection.Close()
			return err
		},
	})

}
