package main

import (
	server "2021_2_MAMBa/internal/app"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"time"
)

const serverPort = ":8080"
const collPort = "50040"
const authPort = "50041"

func main() {
	go func() {
		opt := option.WithCredentialsFile("firebasePrivateKey.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Warn(fmt.Sprintf("error initializing app: %v\n", err))
		}

		ctx := context.Background()
		client, err := app.Messaging(ctx)
		if err != nil {
			log.Warn(fmt.Sprintf("error getting Messaging client: %v\n", err))
		}

		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: "Test title",
				Body:  "Test description",
			},
			Topic: "all",
		}

		ticker := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-ticker.C:
				response, err := client.Send(ctx, message)
				if err != nil {
					log.Error(err)
				}
				log.Info("Successfully sent message: " + response)
			}
		}

	}()
	server.RunServer(serverPort, collPort, authPort)
}
