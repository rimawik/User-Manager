package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ditointernet/go-assert"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/user-manager/app"
	"github.com/user-manager/config"
	"github.com/user-manager/db"
	"github.com/user-manager/handlers"
	"github.com/user-manager/models"
)

func TestAddUsers(t *testing.T) {
	// Create a new router
	r := mux.NewRouter()

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}

	logChannel := make(chan models.RequestLog)
	testApp := app.BuildApp(cfg, postgresDB, logChannel)

	r.HandleFunc("/users", handlers.AddUsers(testApp)).Methods("POST")

	userToBeAdded := map[string]interface{}{
		"name":  "Jack",
		"email": "Jack@gmail.com",
		"age":   28,
	}
	userJSON, _ := json.Marshal(userToBeAdded)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var user models.User
	err = json.NewDecoder(rr.Body).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}
	// Check the response body

	assert.Equal(t, user.Name,
		userToBeAdded["name"])
}

func TestDeleteUsers(t *testing.T) {

	//prepare db and configs
	r := mux.NewRouter()

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}
	logChannel := make(chan models.RequestLog)

	testApp := app.BuildApp(cfg, postgresDB, logChannel)

	//test case 1: user not found
	r.HandleFunc("/users/{id}", handlers.DeleteUsers(testApp)).Methods("DELETE")

	req, err := http.NewRequest("DELETE", "/users/1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	//test case 2: user successfully deleted
	r.HandleFunc("/users/{id}", handlers.DeleteUsers(testApp)).Methods("DELETE")

	req, err = http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

}

func TestGetUsers(t *testing.T) {

	//prepare db and configs
	r := mux.NewRouter()

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}

	logChannel := make(chan models.RequestLog)

	testApp := app.BuildApp(cfg, postgresDB, logChannel)

	//test case 1: user not found
	r.HandleFunc("/users/{id}", handlers.GetUsers(testApp)).Methods("GET")

	req, err := http.NewRequest("GET", "/users/1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	//test case 2: get user successfully
	r.HandleFunc("/users/{id}", handlers.GetUsers(testApp)).Methods("GET")

	req, err = http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var user models.User
	err = json.NewDecoder(rr.Body).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user.ID, 2)

}
