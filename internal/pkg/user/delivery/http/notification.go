package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/mailru/easyjson"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (handler *UserHandler) AddUserToNotificationTopic(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tokenForm := new(domain.UserNotificationToken)
	var p []byte
	p, err := ioutil.ReadAll(r.Body)
	err = tokenForm.UnmarshalJSON(p)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

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

	registrationTokens := []string{tokenForm.Token}

	response, err := client.SubscribeToTopic(ctx, registrationTokens, "all")
	if err != nil {
		log.Error(err)
	}
	log.Info(strconv.Itoa(response.SuccessCount) + " tokens were subscribed successfully")

	us := domain.UserNotificationToken{}
	x, err := easyjson.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}