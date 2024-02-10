package server

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Expression struct {
	ID       string  `json:"id"`
	MathExpr string  `json:"mathExpr"`
	Status   string  `json:"status"`
	Result   float64 `json:"result"`
}

func database() *sql.DB {
	connStr := "user=postgres password=12345 dbname=postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Expressions (id SERIAL PRIMARY KEY, MathExpr TEXT, Status TEXT, Result INTEGER)")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("internal/pages/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func GetExpressions(w http.ResponseWriter, r *http.Request) {
	db := database()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Expressions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := []Expression{} // array of users
	for rows.Next() {
		var u Expression
		if err := rows.Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(users)
}

// get user by id
func GetExpression(w http.ResponseWriter, r *http.Request) {
	db := database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u Expression
	err := db.QueryRow("SELECT * FROM Expressions WHERE id = $1", id).Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// create user
func AddExpression(w http.ResponseWriter, r *http.Request) {
	db := database()
	defer db.Close()
	var u Expression
	json.NewDecoder(r.Body).Decode(&u)

	err := db.QueryRow("INSERT INTO Expressions (MathExpr, Status, Result) VALUES ($1, $2, $3) RETURNING ID", u.MathExpr, u.Status, u.Result).Scan(&u.ID)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(u)
}

func DeleteExpression(w http.ResponseWriter, r *http.Request) {
	db := database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u Expression
	err := db.QueryRow("SELECT * FROM Expressions WHERE id = $1", id).Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		_, err := db.Exec("DELETE FROM Expressions WHERE id = $1", id)
		if err != nil {
			// todo : fix error handling
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("User deleted")
	}
}
