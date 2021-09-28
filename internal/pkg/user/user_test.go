package user

import (
	"2021_2_MAMBa/internal/pkg/database"
	"fmt"
	"github.com/gorilla/mux"
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
}

/*
	{
		inQuery: "",
		bodyString: ``,
		out: ``,
		status: http.StatusOK,
		name: "register one",
	},
*/
var testTableRegisterSuccess = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan","surname": "Ivanov","email": "ivan1@mail.ru","password": "123456","password_repeat": "123456"}`,
		out:        `{"id":1,"first_name":"Ivan","surname":"Ivanov","email":"ivan1@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusCreated,
		name:       "register one",
	},
}
var testTableRegisterFailure = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan","surname": "Ivanov","email": "ivan1@mail.ru","password": "123456","password_repeat": "123456"}`,
		out:        errorAlreadyIn + "\n",
		status:     http.StatusConflict,
		name:       "already in",
	},
	{
		inQuery:    "",
		bodyString: `{"first_nme": "Ivan",}`,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "bad fields",
	},
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan",}`,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "empty fields",
	},
}

func TestRegisterSuccess(t *testing.T) {
	for _, test := range testTableRegisterSuccess {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/register"+test.inQuery, bodyReader)
		Register(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func TestRegisterFailure(t *testing.T) {
	db.AddUser(&database.User{
		FirstName:  "Ivan",
		Surname:    "Ivanov",
		Password:   "123456",
		Email:      "ivan1@mail.ru",
		ProfilePic: "/pic/1.jpg",
	})
	for _, test := range testTableRegisterFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/register"+test.inQuery, bodyReader)
		Register(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func fillMockDB() {
	db = database.Database{}
	db.AddUser(&database.User{
		FirstName:  "Ivan",
		Surname:    "Ivanov",
		Password:   "123456",
		Email:      "ivan@mail.ru",
		ProfilePic: "/pic/1.jpg",
	})
	db.AddUser(&database.User{
		FirstName:  "Ivan",
		Surname:    "Ivanov",
		Password:   "123456",
		Email:      "iva21@mail.ru",
		ProfilePic: "/pic/1.jpg",
	})
}

var testTableGetSuccess = [...]testRow{
	{
		inQuery:    "2",
		bodyString: ``,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "find user",
	},
}

var testTableGetFailure = [...]testRow{
	{
		inQuery:    "3",
		bodyString: ``,
		out:        errorBadInput + "\n",
		status:     http.StatusNotFound,
		name:       "out of index",
	},
	{
		inQuery:    "a",
		bodyString: ``,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "no uinteger",
	},
	{
		inQuery:    "",
		bodyString: ``,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "empty",
	},
}

func TestGetBasicInfoSuccess(t *testing.T) {
	fillMockDB()
	for _, test := range testTableGetSuccess {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/get/"+test.inQuery, bodyReader)
		// Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"id": test.inQuery,
		}
		r = mux.SetURLVars(r, vars)
		GetBasicInfo(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func TestGetBasicInfoFailure(t *testing.T) {
	fillMockDB()
	for _, test := range testTableGetFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/get/"+test.inQuery, bodyReader)
		// Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"id": test.inQuery,
		}
		r = mux.SetURLVars(r, vars)
		GetBasicInfo(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

var testTableLoginSuccess = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "log in user",
	},
}

var testTableLoginFailure = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"email": "raddom@mail.su","password": "123456"}`,
		out:        errorBadCredentials + "\n",
		status:     http.StatusUnauthorized,
		name:       "user not in base",
	},
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "122456"}`,
		out:        errorBadCredentials + "\n",
		status:     http.StatusUnauthorized,
		name:       "wrong pass",
	},
	{
		inQuery:    "",
		bodyString: `{"password": "122456"}`,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "no email",
	},
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru"}`,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "no pass",
	},
	{
		inQuery:    "",
		bodyString: `{"emal": "iva21@mail.ru","password": "123456"}`,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "wrong json",
	},
}

func TestLoginSuccess(t *testing.T) {
	fillMockDB()
	for _, test := range testTableLoginSuccess {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/login"+test.inQuery, bodyReader)
		Login(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func TestLoginFailure(t *testing.T) {
	fillMockDB()
	for _, test := range testTableLoginFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/login"+test.inQuery, bodyReader)
		Login(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}
