package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/aslon1213/comnet_task/internal/app/models"
	"github.com/aslon1213/comnet_task/internal/pkg/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// var router *app.App

func TestMain(M *testing.M) {
	// create db
	os.Chdir("/Users/aslonkhamidov/Desktop/code/tasks/comnet_task/")
	file, err := os.Create("db/db.sqlite3")
	log.Println("db created")
	if err != nil {
		panic(err)

	}
	// time.Sleep(time.Second * 15)
	file.Close()
	gin.SetMode(gin.TestMode)
	exit := M.Run()

	// delete db
	os.Exit(exit)
}
func router() *app.App {
	router := app.New()
	return router
}

func makeRequest(method, path string, body []byte, form_data url.Values, authrequired bool) *httptest.ResponseRecorder {

	buf := bytes.Buffer{}
	buf.Write(body)
	// fmt.Println(buf)
	request, _ := http.NewRequest(method, path, &buf)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if authrequired {

		token := GetToken()
		request.AddCookie(token)
	}

	writer := httptest.NewRecorder()

	router().Gin.ServeHTTP(writer, request)
	return writer
}

func TestRegister(t *testing.T) {
	user := models.User{
		Name:     "test1",
		Login:    "test1",
		Password: "test1",
		Age:      25,
	}
	user2 := models.User{
		Name:     "test2",
		Login:    "test2",
		Password: "test2",
		Age:      25,
	}
	users := []models.User{user, user2}
	for _, user := range users {
		data := url.Values{}
		data.Add("name", user.Name)
		data.Add("login", user.Login)
		data.Add("password", user.Password)
		data.Add("age", fmt.Sprintf("%d", user.Age))

		writer := makeRequest("POST", "/user/register", []byte(data.Encode()), data, false)
		// fmt.Println(string(writer.Body.Bytes()))
		assert.Equal(t, http.StatusCreated, writer.Code)
	}
}

func TestRegisterFail(t *testing.T) {
	user := models.User{
		Name:     "test1",
		Login:    "test1",
		Password: "test1",
		Age:      25,
	}
	data := url.Values{}
	data.Add("name", user.Name)
	data.Add("login", user.Login)
	data.Add("password", user.Password)
	data.Add("age", fmt.Sprintf("%d", user.Age))

	writer := makeRequest("POST", "/user/register", []byte(data.Encode()), data, false)
	// fmt.Println(string(writer.Body.Bytes()))
	assert.Equal(t, http.StatusFound, writer.Code)
	assert.Contains(t, writer.Body.String(), "{\"error\":\"user already exists\"}")

}

func GetToken() *http.Cookie {
	var user_details struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	user_details.Login = "test1"
	user_details.Password = "test1"
	body, _ := json.Marshal(user_details)
	writer := makeRequest("GET", "/user/auth", body, nil, false)
	if writer.Code != http.StatusOK {
		panic("error while getting token")
	}
	cookies := writer.Result().Cookies()
	return cookies[0]

}

func CreateTestUsers(t *testing.T) {

}

func TestGetUserByName(t *testing.T) {
	writer := makeRequest("GET", "/user/test1", nil, nil, true)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "{\"age\":25,\"id\":1,\"name\":\"test1\"}")
	// fmt.Println(writer.Body.String())
}

func TestGetUserByNameFail(t *testing.T) {
	writer := makeRequest("GET", "/user/test3", nil, nil, true)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.NotContains(t, writer.Body.String(), "test3")
	// fmt.Println(writer.Body.String())
}

func TestCreateUserPhone(t *testing.T) {
	var phone struct {
		Phone       string `json:"phone"`
		Description string `json:"description"`
		IsFax       bool   `json:"is_fax"`
	}
	phone.Phone = "+9998998998"
	phone.Description = "test phone"
	phone.IsFax = false
	body, _ := json.Marshal(phone)
	writer := makeRequest("POST", "/user/phone", body, nil, true)
	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Contains(t, writer.Body.String(), "{\"inserted_id\":1,\"message\":\"Phone created succesfully\"}")
	// fmt.Println(writer.Body.String())
}

func TestCreateUserPhoneFailByLength(t *testing.T) {
	var phone struct {
		Phone       string `json:"phone"`
		Description string `json:"description"`
		IsFax       bool   `json:"is_fax"`
	}
	phone.Phone = "+9998998998000000000"
	phone.Description = "test phone"
	phone.IsFax = false
	body, _ := json.Marshal(phone)
	writer := makeRequest("POST", "/user/phone", body, nil, true)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Contains(t, writer.Body.String(), "length")

}

func TestCreateUserPhoneFailByUniqueConstraint(t *testing.T) {
	var phone struct {
		Phone       string `json:"phone"`
		Description string `json:"description"`
		IsFax       bool   `json:"is_fax"`
	}
	phone.Phone = "+9998998998"
	phone.Description = "test phone"
	phone.IsFax = false
	body, _ := json.Marshal(phone)
	writer := makeRequest("POST", "/user/phone", body, nil, true)
	assert.Equal(t, 500, writer.Code)
	assert.Contains(t, writer.Body.String(), "{\"error\":\"UNIQUE constraint failed: phones.phone\"}")

}

func TestGetPhonesByQuery(t *testing.T) {
	writer := makeRequest("GET", "/user/phone?q=8998", nil, nil, true)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "8998")
	// fmt.Println(writer.Body.String())
}

func TestUpdatePhone(t *testing.T) {
	inp := models.PhoneUpdateInput{}
	inp.Phone = "+9998998998"
	inp.Description = "test phone"
	inp.IsFax = false
	inp.PhoneId = 1
	body, _ := json.Marshal(inp)
	writer := makeRequest("PUT", "/user/phone", body, nil, true)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "{\"message\":\"phone updated successfully\"}")
	// fmt.Println(writer.Body.String())
}
