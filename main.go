package main

import (
	"go.uber.org/fx"
	"golang-chat/bootstrap"
)

func main() {
	fx.New(bootstrap.Module).Run()
}
