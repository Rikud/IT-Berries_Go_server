package DA

import (
	"IT-Berries/auth/entities"
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "it_berries"

	connectError = "An error occurred while trying to connect to the database."
	pignError    = "An error occurred while trying to ping database"
	executeError = "An error occurred while trying to execute query."
	readRowError = "An error occurred while trying to read row"
)

var psqlInfo string

func init() {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db := connect()
	defer db.Close()
	err := db.Ping()
	errorCheck(err, pignError)
	fmt.Println("Successfully connected to database.")

}

func FindUserByEmail(email string) []*entities.User {
	db := connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM USERS_ WHERE email = $1", email)
	errorCheck(err, executeError)
	defer rows.Close()
	users := make([]*entities.User, 0)
	for rows.Next() {
		user := new(entities.User)
		err := rows.Scan(user.GetAvatarPoint(), user.GetEmailPoint(), user.GetPasswordPoint(), user.GetUsernamePoint())
		errorCheck(err, readRowError)
		users = append(users, user)
	}
	errorCheck(rows.Err(), readRowError)
	return users
}

func errorCheck(err error, message string) {
	if err != nil {
		log.Println(message, err)
		panic(err)
	}
}

func connect() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	errorCheck(err, connectError)
	return db
}
