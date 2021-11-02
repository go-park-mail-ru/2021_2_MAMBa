package queryChecker

import (
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)


type testRow struct {
	inQuery    string
	bodyString string
	outint64 uint64
	outint int
	outFloat float64
	err error
	name string
}

var test1 = testRow{
	inQuery: "skip=1",
	bodyString: "skip",
	name: "test qc",
	err: nil,
	outint: 1,
}

func TestQueryIntSuccess(t *testing.T) {
	bodyReader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/query?"+test1.inQuery, bodyReader)
	actual, err := CheckIsIn(w,r,test1.bodyString,0, customErrors.ErrorBadInput)
	assert.Equal(t, test1.outint, actual, "Test: "+test1.name)
	assert.Equal(t, test1.err, err, "Test: "+test1.name)
}

var test2 = testRow{
	inQuery: "skip=-1",
	bodyString: "skip",
	name: "test qc",
	err: customErrors.ErrorBadInput,
	outint: 0,
}

func TestQueryIntFail(t *testing.T) {
	bodyReader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/query?"+test2.inQuery, bodyReader)
	actual, err := CheckIsIn(w,r,test2.bodyString,0, customErrors.ErrorBadInput)
	assert.Equal(t, test2.outint, actual, "Test: "+test2.name)
	assert.Equal(t, test2.err, err, "Test: "+test2.name)
}

var test3 = testRow{
	inQuery: "id=1",
	bodyString: "id",
	name: "test qc",
	err: nil,
	outint64: 1,
}

func TestQueryInt64Success(t *testing.T) {
	bodyReader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/query?"+test3.inQuery, bodyReader)
	actual, err := CheckIsIn64(w,r,test3.bodyString,0, customErrors.ErrorBadInput)
	assert.Equal(t, test3.outint64, actual, "Test: "+test3.name)
	assert.Equal(t, test3.err, err, "Test: "+test3.name)
}

var test4 = testRow{
	inQuery: "id=-1",
	bodyString: "id",
	name: "test qc",
	err: customErrors.ErrorBadInput,
	outint64: 0,
}

func TestQueryInt64Fail(t *testing.T) {
	bodyReader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/query?"+test4.inQuery, bodyReader)
	actual, err := CheckIsIn64(w,r,test4.bodyString,0, customErrors.ErrorBadInput)
	assert.Equal(t, test4.outint64, actual, "Test: "+test4.name)
	assert.Equal(t, test4.err, err, "Test: "+test4.name)
}

var test5 = testRow{
	inQuery: "stars=4.5",
	bodyString: "stars",
	name: "test qc",
	err: nil,
	outFloat: 4.5,
}

func TestQueryFloat64Success(t *testing.T) {
	bodyReader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/query?"+test5.inQuery, bodyReader)
	actual, err := CheckIsInFloat64(w,r,test5.bodyString,0, customErrors.ErrorBadInput)
	assert.Equal(t, test5.outFloat, actual, "Test: "+test5.name)
	assert.Equal(t, test5.err, err, "Test: "+test5.name)
}

var test6 = testRow{
	inQuery: "stars=-1",
	bodyString: "stars",
	name: "test qc",
	err: customErrors.ErrorBadInput,
	outint64: 0,
}

func TestQueryFloat64Fail(t *testing.T) {
	bodyReader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/query?"+test6.inQuery, bodyReader)
	actual, err := CheckIsInFloat64(w,r,test6.bodyString,0, customErrors.ErrorBadInput)
	assert.Equal(t, test6.outFloat, actual, "Test: "+test6.name)
	assert.Equal(t, test6.err, err, "Test: "+test6.name)
}