package cmd

import (
	"log"
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/app"
	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/http/routes"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the backend server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		config.InitDB()

		application := app.InitApp(config.AppConfig)

		router := routes.RegisterRoutes(application)

		addr := ":8080"
		log.Printf("Server running on %s", addr)
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
