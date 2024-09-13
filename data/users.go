package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
	"github.com/user-manager/app"
	"github.com/user-manager/models"
)

func GetUsers(app *app.App,
	id string) (user models.User,
	err error) {

	err = app.PostgresDB().Conn.QueryRow(`SELECT id,
	 email,
	 name,
	 age from users where id = $1`, id).Scan(&user.ID,
		&user.Email,
		&user.Name,
		&user.Age)

	if err != nil {
		log.Errorf("Couldn't query user details: %v", err)
		return user, err
	}
	if user.ID == 0 {
		log.Errorf("user doesn't exist")
		return user, sql.ErrNoRows
	}
	return user, nil
}

func AddUsers(app *app.App,
	userTobeAdded models.User) (models.User,
	error) {

	err := app.PostgresDB().Conn.QueryRow(`INSERT INTO users ("name","age","email") values($1,$2,$3) returning id`,
		userTobeAdded.Name,
		userTobeAdded.Age,
		userTobeAdded.Email).Scan(&userTobeAdded.ID)

	if err != nil {
		log.Errorf("Couldn't insert user: %v", err)
		return models.User{}, err
	}

	return userTobeAdded, nil
}

func EditUsers(app *app.App,
	userTobeAdded models.User, id string) error {

	_, err := app.PostgresDB().Conn.Exec(`update users set name = coalesce($1, name),
	age = coalesce($2, age),
	email= coalesce($3, email),
	where id = $4`,
		userTobeAdded.Name,
		userTobeAdded.Age,
		userTobeAdded.Email,
		id)

	if err != nil {
		log.Errorf("Couldn't Edit user: %v", err)
		return err
	}

	return nil
}
func DeleteUsers(app *app.App,
	id string) error {

	result, err := app.PostgresDB().Conn.Exec(`delete from users where id = $1`, id)

	if err != nil {
		log.Errorf("Couldn't delete user: %v", err)
		return err
	}
	if x, _ := result.RowsAffected(); x == 0 {
		log.Errorf("not found")
		return sql.ErrNoRows
	}

	return nil
}
