package server

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Expression struct {
	ID       string  `json:"id"`
	MathExpr string  `json:"mathExpr"`
	Status   string  `json:"status"`
	Result   float64 `json:"result"`
}

func Database() *sql.DB {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA journal_mode=WAL")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Expressions (id INTEGER PRIMARY KEY AUTOINCREMENT, MathExpr TEXT, Status TEXT, Result REAL)")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Home(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl.Execute(w, nil)
}

func GetExpressions(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl.Execute(w, nil)

	db := Database()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Expressions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := map[string]Expression{} // map of users
	for rows.Next() {
		var u Expression
		if err := rows.Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result); err != nil {
			log.Fatal(err)
		}
		users[u.ID] = u
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(users)
}

// get user by id
func GetExpression(w http.ResponseWriter, r *http.Request) {
	db := Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u Expression
	err := db.QueryRow("SELECT * FROM Expressions WHERE id = ?", id).Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// create user
func AddExpression(w http.ResponseWriter, r *http.Request) {
	db := Database()
	defer db.Close()
	var u Expression
	json.NewDecoder(r.Body).Decode(&u)

	res, err := db.Exec("INSERT INTO Expressions (MathExpr, Status, Result) VALUES (?, ?, ?)", u.MathExpr, "pending", u.Result)
	if err != nil {
		log.Fatal(err)
	}

	id, _ := res.LastInsertId()
	u.ID = string(id)

	json.NewEncoder(w).Encode(u)
}

func DeleteExpression(w http.ResponseWriter, r *http.Request) {
	db := Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u Expression
	err := db.QueryRow("SELECT * FROM Expressions WHERE id = ?", id).Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		_, err := db.Exec("DELETE FROM Expressions WHERE id = ?", id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("User deleted")
	}
}

func DeleteExpressions(w http.ResponseWriter, r *http.Request) {
	db := Database()
	defer db.Close()
	rows, err := db.Query("SELECT ID FROM Expressions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var u Expression
		if err := rows.Scan(&u.ID); err != nil {
			log.Fatal(err)
		}
		_, err := db.Exec("DELETE FROM Expressions WHERE id = ?", u.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
