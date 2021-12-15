package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

var testTableGetCollsSuccess = [...]testRow{
	{
		inQuery: "id=8&skip_reviews=0&limit_reviews=10&skip_recommend=0&limit_recommend=10",
		out:     `{"body":{"collections_list":[{"id":1,"title":"Для ценителей Хогвартса","picture_url":"/static/media/img/collections/1.webp"}],"more_available":true,"collection_total":12,"current_sort":"","current_limit":1,"current_skip":1},"status":200}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
	},
}

var testTableGetCollsFailure = [...]testRow{
	{
		inQuery: "id=8&skip=-1&limits=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "id=8&skip=11&limit=-2",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
}

func TestGetCollsFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/collections/getCollections?"
	for _, test := range testTableGetCollsFailure {
		var cl domain.Collections
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		bodyReader := strings.NewReader("")
		handler:=CollectionsHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetCollections(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
