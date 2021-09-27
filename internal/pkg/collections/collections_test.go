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
	out string
	status int
	name string
}

var testTable = [...]testRow {
	{
		inQuery: "skip=0&limit=1",
		out:`{"collections_list":[{"id":1,"title":"Для ценителей Хогвардса","picture_url":"server/images/collections1.png"}],"more_avaliable":true,"collection_total":12,"current_sort":"","current_limit":1,"currentSkip":1}`+"\n",
		status: http.StatusOK,
		name: `function works`,
	},
}


func TestGetCollectionsSuccess(t *testing.T) {
		for _, test := range testTable {
			fmt.Fprintf(os.Stdout, "Test:" + test.name)
			bodyReader := strings.NewReader("")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/collections/getCollections?"+test.inQuery, bodyReader)
			GetCollections(w, r)
			assert.Equal(t, test.out, w.Body.String(), "Test: " + test.name)
			assert.Equal(t, test.status, w.Code, "Test: " + test.name)
			fmt.Fprintf(os.Stdout, " done\n")
		}
}