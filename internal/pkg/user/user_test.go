package user

import (
	"2021_2_MAMBa/internal/pkg/database"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan12","surname": "Ivanov","email": "ivan131@mail.ru","password": "123455","password_repeat": "123456"}`,
		out:        errorBadInput + "\n",
		status:     http.StatusBadRequest,
		name:       "unmatching passwords",
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
	apiPath := "/api/user/register"
	for _, test := range testTableRegisterFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		Register(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func fillMockDB() {
	db = database.UserMockDatabase{}
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
	apiPath := "/api/user/get/"
	for _, test := range testTableGetFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
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
	apiPath := "/api/user/login"
	for _, test := range testTableLoginSuccess {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		Login(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func TestLoginFailure(t *testing.T) {
	fillMockDB()
	apiPath := "/api/user/login"
	for _, test := range testTableLoginFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		Login(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

var testTableLogoutFailure = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        errorBadInput + "\n",
		status:     http.StatusForbidden,
		name:       "logout not logged in",
	},
}

func TestLogoutFailure(t *testing.T) {
	fillMockDB()
	apiPath := "/api/user/logout"
	for _, test := range testTableLogoutFailure {
		fmt.Fprintf(os.Stdout, "Test:"+test.name)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		Logout(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
		fmt.Fprintf(os.Stdout, " done\n")
	}
}

func TestLogoutSuccess(t *testing.T) {
	fillMockDB()
	bodyReader := strings.NewReader(testTableLoginSuccess[0].bodyString)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bodyReader)
	Login(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	r = httptest.NewRequest("GET", "/api/user/logout", bodyReader)
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	w = httptest.NewRecorder()
	Logout(w, r)
	assert.Equal(t, http.StatusOK, w.Code, "Test: Logout OK")

}

func TestCheckAuthSuccess(t *testing.T) {
	fillMockDB()
	bodyReader := strings.NewReader(testTableLoginSuccess[0].bodyString)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bodyReader)
	Login(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	r = httptest.NewRequest("GET", "/api/user/checkAuth", bodyReader)
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	w = httptest.NewRecorder()
	CheckAuth(w, r)
	assert.Equal(t, http.StatusOK, w.Code, "Test: Logout OK")
}

func TestCheckAuthFailure(t *testing.T) {
	fillMockDB()
	bodyReader := strings.NewReader("")
	r := httptest.NewRequest("GET", "/api/user/checkAuth", bodyReader)
	w := httptest.NewRecorder()
	CheckAuth(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Test: Logout OK")
}