package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Employee Struct
type Employee struct {
	ID   int
	Name string
	City string
}

// Open Connection with MySQL Driver
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

// Read All Templates on folder template
var tmpl = template.Must(template.ParseGlob("template/*"))

// Index Page
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
		emp.ID = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

// Show Single Item
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
		emp.ID = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

// Show New Page
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Edit Item
func Edit(w http.ResponseWriter, r *http.Request) {
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
		emp.ID = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

// Insert Item
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		insert, err := db.Prepare("INSERT INTO Employee (name, city) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Update Item
func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		insert, err := db.Prepare("UPDATE Employee SET name = ?, city = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Delete Item
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	emp := r.URL.Query().Get("id")
	delete, err := db.Prepare("DELETE FROM Employee WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	delete.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	// Routes
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	// Start server on port 8080
	http.ListenAndServe(":8080", nil)
}
