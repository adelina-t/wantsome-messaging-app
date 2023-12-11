package server

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func checkIfUserIsInDb() {

}

func addUserInDb() {

}

func insertMessageInDb() {

}

func createDb() {
	queryString := "CREATE TABLE 'test' ('username' VARCHAR(64) null)"
	db, err := sql.Open("sqlite3", "foo.db")
	if err != nil {
		log.Printf("error connecting to the db %s\n", err)
	}
	_, err = db.Exec(queryString)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}

func createUserTable() {

}

func createMessageTable() {

}
