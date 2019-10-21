package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
	"time"
)
var err error
var st *sql.Stmt
// checkDatabase will check if the db already exists else it will create it
func checkDatabase() {
	var info os.FileInfo



	info, err = os.Stat(configuration.DBName)
	if !os.IsNotExist(err) {
		log.Print("Db already created")
	}
	// Directory doesn't exist so we will create it
	log.Println("Creating db file")
	os.Create(configuration.DBName)
	time.Sleep(1000000000) // delay to give system time to create file just in case

	if info.IsDir(){
		log.Panic("A directory with db name.... not gonna work")
	}
	db, err := sql.Open("sqlite3", "db.sqlite3")
	checkErr(err)

	st, err = db.Prepare("CREATE TABLE IF NOT EXISTS people " +
		"(dni VARCHAR PRIMARY KEY, " +
		"name TEXT)" )

	res, err := st.Exec()

	st, err = db.Prepare("CREATE TABLE IF NOT EXISTS ticket " +
		"(id INTEGER  NOT NULL, " +
		"hash binary(20), dni VARCHAR,  PRIMARY KEY (id))" )

	res, err = st.Exec()

	checkErr(err)
	log.Println(res)


	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertData( students []Student)  {
	log.Print("Inserting People data to db")
	db, err := sql.Open("sqlite3", "db.sqlite3")
	checkErr(err)
	for _, student := range students {
		st, err := db.Prepare("SELECT dni, name FROM people WHERE dni=$1 ")
		checkErr(err)
		rows, err := st.Query(student.Dni)
		checkErr(err)
		if !rows.Next() {
			st, err = db.Prepare("INSERT INTO people (dni, name) VALUES (?, ?)" )
			checkErr(err)
			_ , err = st.Exec(student.Dni, student.Name)
			checkErr(err)

		}
		_ = rows.Close()
		for j := 0; j < len(student.Keys); j++ {
			st, err = db.Prepare("INSERT INTO ticket (dni, hash) VALUES (?, ?)" )
			checkErr(err)
			_ , err = st.Exec(student.Dni, student.Keys[j])
			checkErr(err)
		}
	}
	db.Close()
	log.Print("People data saved")
}