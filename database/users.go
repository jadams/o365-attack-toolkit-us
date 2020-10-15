package database

import (
	_ "database/sql"
	"log"
	"github.com/jadams/o365-attack-toolkit-us/model"

	_ "github.com/mattn/go-sqlite3"
)

func GetUsers() []model.User {

	var users []model.User

	rows, err := db.Query(model.GetUsersQuery)
	user := model.User{}

	if err != nil {
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.DisplayName, &user.Mail, &user.JobTitle, &user.UserPrincipalName, &user.AccessToken, &user.AccessTokenActive, &user.RefreshToken)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}

func InsertUser(user model.User) {

	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare(model.InsertUserQuery)

	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	_, err := stmt.Exec(user.Id, user.DisplayName, user.Mail, user.JobTitle, user.UserPrincipalName, user.AccessToken, user.AccessTokenActive, user.RefreshToken)
	tx.Commit()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

}
