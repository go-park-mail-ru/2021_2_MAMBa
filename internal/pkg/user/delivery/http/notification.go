package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"net/http"
	"strconv"
)

func (handler *UserHandler) AddUserToNotificationTopic(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tokenForm := new(domain.UserNotificationToken)
	err := json.NewDecoder(r.Body).Decode(tokenForm)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	opt := option.WithCredentialsFile("firebasePrivateKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Warn(fmt.Sprintf("error initializing app: %v\n", err))
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Warn(fmt.Sprintf("error getting Messaging client: %v\n", err))
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	registrationTokens := []string{tokenForm.Token}

	response, err := client.SubscribeToTopic(ctx, registrationTokens, "all")
	if err != nil {
		log.Error(err)
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	log.Info(strconv.Itoa(response.SuccessCount) + " tokens were subscribed successfully")

	us := domain.UserNotificationToken{}
	x, err := json.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

// SendPushToAll - Ручка для демонстрации на РК4
func (handler *UserHandler) SendPushToAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	notificationForm := new(domain.UserNotificationToSend)
	err := json.NewDecoder(r.Body).Decode(notificationForm)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	opt := option.WithCredentialsFile("firebasePrivateKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Warn(fmt.Sprintf("error initializing app: %v\n", err))
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Warn(fmt.Sprintf("error getting Messaging client: %v\n", err))
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notificationForm.Title,
			Body:  notificationForm.Description,
		},
		Topic: "all",
	}
	_, err = client.Send(ctx, message)
	if err != nil {
		log.Error(err)
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	log.Info("Successfully sent push to all users")

	us := domain.UserNotificationToken{}
	x, err := json.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
