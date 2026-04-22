// @title Parkin API
// @version 1.0
// @description This is Parkin backend API
// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"github.com/capstone-orang-jago/parkin-be/cmd"
	"github.com/capstone-orang-jago/parkin-be/config"
)

func main() {
	config.LoadConfig()
	config.InitDB()

	cmd.Execute()
}
