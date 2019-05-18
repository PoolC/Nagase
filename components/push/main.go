package push

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var ctx context.Context
var client *messaging.Client

func RegisterToken(memberUUID string, pushToken string) error {
	_, err := client.SubscribeToTopic(ctx, []string{pushToken}, memberUUID)
	if err != nil {
		return err
	}
	return nil
}

func DeregisterToken(memberUUID string, pushToken string) error {
	_, err := client.UnsubscribeFromTopic(ctx, []string{pushToken}, memberUUID)
	if err != nil {
		return err
	}
	return nil
}

func SendPush(memberUUID string, title string, body string, data map[string]string) error {
	dataMessage := make(map[string]string)
	dataMessage["title"] = title
	dataMessage["body"] = body

	for key, value := range data {
		dataMessage[key] = value
	}

	message := messaging.Message{
		Topic: memberUUID,
		Webpush: &messaging.WebpushConfig{
			Data: dataMessage,
		},
	}

	_, err := client.Send(ctx, &message)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	credPath := ""
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(dir + "/secrets/service-account.json"); err == nil {
			credPath = dir + "/secrets/service-account.json"
			break
		}

		if dir == "/" {
			panic("service-account.json not found")
		}
		dir = filepath.Dir(dir)
	}

	opt := option.WithCredentialsFile(credPath)
	config := firebase.Config{ProjectID: "poolc-b18fa"}
	app, err := firebase.NewApp(context.Background(), &config, opt)
	if err != nil {
		fmt.Errorf("Failed to initialized firebase application")
		panic(err)
	}

	ctx = context.Background()
	newClient, err := app.Messaging(ctx)
	if err != nil {
		fmt.Errorf("Failed to initialized firebase messaging")
		panic(err)
	}

	client = newClient
}
