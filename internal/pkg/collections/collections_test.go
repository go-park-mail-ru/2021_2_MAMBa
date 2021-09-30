package collections

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type testRow struct {
	inQuery string
	out     string
	status  int
	name    string
}

var testTableSuccess = [...]testRow{
	{
		inQuery: "skip=0&limit=1",
		out:     `{"collections_list":[{"id":1,"title":"Для ценителей Хогвартса","picture_url":"server/images/collections1.png"}],"more_available":true,"collection_total":12,"current_sort":"","current_limit":1,"current_skip":1}` + "\n",
		status:  http.StatusOK,
		name:    `limit works`,
	},
	{
		inQuery: "skip=10&limit=1",
		out:     `{"collections_list":[{"id":11,"title":"Про петлю времени","picture_url":"server/images/collections11.png"}],"more_available":true,"collection_total":12,"current_sort":"","current_limit":1,"current_skip":11}` + "\n",
		status:  http.StatusOK,
		name:    `skip works`,
	},
	{
		inQuery: "skip=11&limit=10",
		out:     `{"collections_list":[{"id":12,"title":"Классика на века","picture_url":"server/images/collections12.jpg"}],"more_available":false,"collection_total":12,"current_sort":"","current_limit":10,"current_skip":21}` + "\n",
		status:  http.StatusOK,
		name:    `does not overflow`,
	},
}
var testTableFailure = [...]testRow{
	{
		inQuery: "skip=-1&limit=10",
		out:     errSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
	},
	{
		inQuery: "skip=11&limit=-2",
		out:     errLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
	},
	{
		inQuery: "skip=14&limit=1",
		out:     errSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
	},
}

func TestGetCollectionsSuccess(t *testing.T) {
	apiPath := "/api/collections/getCollections?"
	for _, test := range testTableSuccess {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath + test.inQuery, bodyReader)
		GetCollections(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: " + test.name)
		assert.Equal(t, test.status, w.Code, "Test: " + test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}
func TestGetCollectionsFailure(t *testing.T) {
	apiPath := "/api/collections/getCollections?"
	for _, test := range testTableFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath + test.inQuery, bodyReader)
		GetCollections(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: " + test.name)
		assert.Equal(t, test.status, w.Code, "Test: " + test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}
