package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

func dbConnect() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "3232002"
	dbName := "goblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("template/*"))

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	rows, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for rows.Next() {
		var id int
		var name, city string
		err = rows.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", emp)
	defer db.Close()
}

// Show single item
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	nID := r.URL.Query().Get("id")
	rows, err := db.Query("SELECT * FROM Employee WHERE id = ?", nID)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for rows.Next() {
		var id int
		var name, city string
		err = rows.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
}

func Edit(w http.ResponseWriter, r *http.Request) {
}

func Insert(w http.ResponseWriter, r *http.Request) {
}

func Update(w http.ResponseWriter, r *http.Request) {
}

func Delete(w http.ResponseWriter, r *http.Request) {
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/Update", Update)
	http.HandleFunc("/Delete", Delete)
	http.ListenAndServe(":8080", nil)
}
