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
