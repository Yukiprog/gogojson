package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	_ "os"
	"strconv"
)

type data struct {
	id   int
	date string
	name string
}

var db *sql.DB
var err error

func main() {
	db_connect()
	routing()
	db.Close()
}
func routing() {
	r := mux.NewRouter()

	//Create
	r.HandleFunc("/upload", uploadFile).Methods("POST")

	//Read
	r.HandleFunc("/read", readRecord).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var record data
	var line []string
	for {
		line, err = reader.Read()
		if err != nil {
			panic(err)
		}
		record.id, _ = strconv.Atoi(line[0])
		record.date = line[1]
		record.name = line[2]
		statement := "insert into w_movie(id,w_date,m_name) values(?,?,?)"
		stmt, err := db.Prepare(statement)
		if err != nil {
			panic(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(record.id, record.date, record.name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d %s %s\n", record.id, record.date, record.name)
	}
}
func readRecord(w http.ResponseWriter, r *http.Request) {
	var record data
	rows, err := db.Query("select * from w_movie")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&record.id, &record.date, &record.name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d %s %s\n", record.id, record.date, record.name)
	}
}
func db_connect() {
	fmt.Printf(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "test", "test", "api-mysql-1", "3306", "movie"))
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "test", "test", "api-mysql-1", "3306", "movie"))
	if err != nil {
		panic(err)
	}
}
