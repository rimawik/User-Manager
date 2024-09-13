package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/user-manager/app"
	"github.com/user-manager/data"
	"github.com/user-manager/models"
)

// GetUsers godoc
// @Summary Get a user
// @Description Get users by ID
// @Tags users
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func GetUsers(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {

			var response []byte

			defer wg.Done()

			vars := mux.Vars(r)
			id := vars["id"]
			if id == "15" {
				time.Sleep(time.Second * 20)
			}
			user, err := data.GetUsers(app, id)
			if err != nil {
				status := http.StatusInternalServerError
				errorStr := err.Error()
				log.Errorf("couldn't get users from database: %s",
					err.Error())
				if err == sql.ErrNoRows {
					status = http.StatusNotFound
					errorStr = "User not found"
				}
				http.Error(w,
					errorStr,
					status)

				return
			}

			response, err = json.Marshal(user)
			if err != nil {
				log.Errorf("couldn't marshal response: %s",
					err.Error())
				http.Error(w,
					"couldn't marshal response",
					http.StatusInternalServerError)

				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}()

		wg.Wait()
		app.LogChannel() <- models.RequestLog{
			Method:     r.Method,
			Path:       r.URL.Path,
			StartTime:  startTime,
			FinishTime: time.Now(),
			Duration:   time.Since(startTime),
		}
	}

}

// AddUsers godoc
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Router /users [post]

func AddUsers(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Errorf("couldn't read request body: %v", err)
				http.Error(w,
					"couldn't read request body",
					http.StatusInternalServerError)
				return
			}
			err = r.Body.Close()
			if err != nil {
				log.Errorf("couldn't close body: %v", err)
				http.Error(w,
					"couldn't close body",
					http.StatusInternalServerError)
				return
			}
			var user models.User
			err = json.Unmarshal(body, &user)
			if err != nil {
				log.Errorf("couldn't unmarshal payload: %v", err)
				http.Error(w,
					"couldn't unmarshal payload",
					http.StatusInternalServerError)
				return
			}

			AddedUser, err := data.AddUsers(app, user)
			if err != nil {
				log.Errorf("couldn't add user to database: %s",
					err.Error())
				http.Error(w,
					"couldn't add user",
					http.StatusInternalServerError)
				return
			}

			log.Info("user was added successfully")

			response, err := json.Marshal(AddedUser)
			if err != nil {
				log.Errorf("couldn't marshal response: %s",
					err.Error())
				http.Error(w,
					"couldn't marshal response",
					http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}()
		wg.Wait()
		app.LogChannel() <- models.RequestLog{
			Method:     r.Method,
			Path:       r.URL.Path,
			StartTime:  startTime,
			FinishTime: time.Now(),
			Duration:   time.Since(startTime),
		}
	}
}

// EditUsers godoc
// @Summary Update a user
// @Description Update user details
// @Tags users
// @Accept json
// @Param user body models.User true "User"
// @Param id path int true "User ID"
// @Success 200
// @Router /users/{id} [patch]
func EditUsers(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			vars := mux.Vars(r)
			id := vars["id"]

			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Errorf("couldn't read request body: %v", err)
				http.Error(w,
					"couldn't read body",
					http.StatusInternalServerError)
				return
			}
			err = r.Body.Close()
			if err != nil {
				log.Errorf("couldn't close body: %v", err)
				http.Error(w,
					"couldn't close body",
					http.StatusInternalServerError)
				return
			}

			var user models.User
			err = json.Unmarshal(body, &user)
			if err != nil {
				log.Errorf("couldn't unmarshal payload: %v", err)
				http.Error(w,
					"couldn't unmarshal payload",
					http.StatusInternalServerError)
				return
			}

			err = data.EditUsers(app, user, id)
			if err != nil {
				log.Errorf("couldn't edit users in database: %s",
					err.Error())
				http.Error(w,
					"couldn't edit user details",
					http.StatusInternalServerError)
				return
			}

			log.Info("user was edited successfully")
			w.WriteHeader(200)
		}()
		wg.Wait()
		app.LogChannel() <- models.RequestLog{
			Method:     r.Method,
			Path:       r.URL.Path,
			StartTime:  startTime,
			FinishTime: time.Now(),
			Duration:   time.Since(startTime),
		}
	}
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200
// @Router /users/{id} [delete]
func DeleteUsers(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			vars := mux.Vars(r)
			id := vars["id"]
			err := data.DeleteUsers(app, id)

			if err != nil {

				log.Errorf("couldn't delete user %s",
					err.Error())
				errorStr := "couldn't delete user: " + err.Error()
				status := http.StatusInternalServerError
				if err == sql.ErrNoRows {
					status = http.StatusNotFound
					errorStr = "User not found"
				}
				http.Error(w,
					errorStr,
					status)
				return
			}

			log.Info("user was deleted successfully")
			w.WriteHeader(200)
		}()
		wg.Wait()
		app.LogChannel() <- models.RequestLog{
			Method:     r.Method,
			Path:       r.URL.Path,
			StartTime:  startTime,
			FinishTime: time.Now(),
			Duration:   time.Since(startTime),
		}
	}
}
