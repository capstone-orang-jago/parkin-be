package cmd

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/capstone-orang-jago/parkin-be/app"
	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/capstone-orang-jago/parkin-be/http/routes"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the backend server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		config.InitDB()

		container, err := config.NewWAContainer()
		if err != nil {
			log.Fatalf("Failed to initialize WhatsApp container: %v", err)
		}

		deviceStore, err := container.GetFirstDevice(context.Background())
		if err != nil {
			log.Fatalf("Failed to get first device: %v", err)
		}

		client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))

		if client.Store.ID == nil {
			qrChan, err := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				log.Fatalf("Failed to connect to WhatsApp: %v", err)
			}

			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else {
					log.Printf("WhatsApp connection event: %s", evt.Event)
				}
			}
		} else {
			err = client.Connect()
			if err != nil {
				log.Fatalf("Failed to connect to WhatsApp: %v", err)
			}
		}

		waService := helpers.NewWaService(client)

		application := app.InitApp(config.AppConfig, waService)

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
