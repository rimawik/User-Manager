package models

// User represents a user in the system
// @Description User represents a user in the system
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int `json:"age"`
}
