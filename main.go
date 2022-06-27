package main

import (
	"hypefast-url-shortener/app"
	"hypefast-url-shortener/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
