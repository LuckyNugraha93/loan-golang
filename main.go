package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Loan struct {
	Id     int
	Name   string
	Amount float32
	Status string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "mydatabase"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Loan ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	loan := Loan{}
	res := []Loan{}
	for selDB.Next() {
		var id int
		var name, status string
		var amount float32
		err = selDB.Scan(&id, &name, &status, &amount)
		if err != nil {
			panic(err.Error())
		}
		loan.Id = id
		loan.Name = name
		loan.Status = status
		loan.Amount = amount
		res = append(res, loan)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Loan WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	loan := Loan{}
	for selDB.Next() {
		var id int
		var name, status string
		var amount float32
		err = selDB.Scan(&id, &name, &status, &amount)
		if err != nil {
			panic(err.Error())
		}
		loan.Id = id
		loan.Name = name
		loan.Status = status
		loan.Amount = amount
	}
	tmpl.ExecuteTemplate(w, "Show", loan)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Loan WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	loan := Loan{}
	for selDB.Next() {
		var id int
		var name, status string
		var amount float32
		err = selDB.Scan(&id, &name, &status, &amount)
		if err != nil {
			panic(err.Error())
		}
		loan.Id = id
		loan.Name = name
		loan.Status = status
		loan.Amount = amount
	}
	tmpl.ExecuteTemplate(w, "Edit", loan)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		amount := r.FormValue("amount")
		insForm, err := db.Prepare("INSERT INTO Loan(name, amount, status) VALUES(?,?,'PENDING')")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, amount)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		amount := r.FormValue("amount")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Loan SET name=?, amount=? , status = 'PENDING' WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, amount, id)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Approve(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	loan := r.URL.Query().Get("id")
	insForm, err := db.Prepare("UPDATE Loan SET status = 'APPROVED' WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(loan)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/approve", Approve)
	http.ListenAndServe(":8080", nil)
}
