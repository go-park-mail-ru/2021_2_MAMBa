package http

import (
	mock2 "2021_2_MAMBa/internal/pkg/collections/usecase/mock"
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/domain/errors"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type testRow struct {
	inQuery    string
	bodyString string
	out        string
	status     int
	name       string
	skip       int
	limit      int
}

var testTableSuccess = [...]testRow{
	{
		inQuery: "skip=0&limit=1",
		out:     `{"collections_list":[{"id":1,"title":"Для ценителей Хогвартса","picture_url":"server/images/collections1.png"}],"more_available":true,"collection_total":12,"current_sort":"","current_limit":1,"current_skip":1}` + "\n",
		status:  http.StatusOK,
		name:    `limit works`,
		skip:    0,
		limit:   1,
	},
	{
		inQuery: "skip=10&limit=1",
		out:     `{"collections_list":[{"id":11,"title":"Про петлю времени","picture_url":"server/images/collections11.png"}],"more_available":true,"collection_total":12,"current_sort":"","current_limit":1,"current_skip":11}` + "\n",
		status:  http.StatusOK,
		name:    `skip works`,
		skip:    10,
		limit:   1,
	},
	{
		inQuery: "skip=11&limit=10",
		out:     `{"collections_list":[{"id":12,"title":"Классика на века","picture_url":"server/images/collections12.jpg"}],"more_available":false,"collection_total":12,"current_sort":"","current_limit":10,"current_skip":21}` + "\n",
		status:  http.StatusOK,
		name:    `does not overflow`,
		skip:    11,
		limit:   10,
	},
}
var testTableFailure = [...]testRow{
	{
		inQuery: "skip=-1&limit=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "skip=11&limit=-2",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
	{
		inQuery: "skip=14&limit=1",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
		skip:    14,
		limit:   1,
	},
}

func TestGetCollectionSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/collections/getCollections?"
	for _, test := range testTableSuccess {
		var cl domain.Collections
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockCollectionsUsecase(ctrl)
		mock.EXPECT().GetCollections(test.skip, test.limit).Times(1).Return(cl, nil)
		handler := CollectionsHandler{CollectionsUsecase: mock}
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetCollections(w, r)
		result:= `{"body":`+test.out[:len(test.out)-1]+`,"status":`+fmt.Sprint(test.status)+"}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func TestGetCollectionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/collections/getCollections?"
	for i, test := range testTableFailure {
		var cl domain.Collections
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockCollectionsUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().GetCollections(test.skip, test.limit).Times(1).Return(domain.Collections{}, customErrors.ErrorSkip)
		}
		handler := CollectionsHandler{CollectionsUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetCollections(w, r)
		result:= `{"body":{"error":"`+test.out[:len(test.out)-1]+`"},"status":`+fmt.Sprint(test.status)+"}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
