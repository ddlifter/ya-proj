package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	database "1/internal/database"

	"github.com/gorilla/mux"
)

func Home(w http.ResponseWriter, r *http.Request) {
	database.Get()
}

// Вывести все задачи
func GetExpressions(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Expressions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	expressions := map[string]database.Expression{} // map of users
	for rows.Next() {
		var u database.Expression
		if err := rows.Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result); err != nil {
			log.Fatal(err)
		}
		expressions[u.ID] = u
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(expressions)
}

// Вывести задачу по id
func GetExpression(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u database.Expression
	err := db.QueryRow("SELECT * FROM Expressions WHERE id = $1", id).Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// Добавить задачу
func AddExpression(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	var u database.Expression
	json.NewDecoder(r.Body).Decode(&u)

	res, err := db.Exec("INSERT INTO Expressions (MathExpr, Status, Result) VALUES ($1, $2, $3)", u.MathExpr, "waiting", u.Result)
	if err != nil {
		log.Fatal(err)
	}

	id, _ := res.LastInsertId()
	u.ID = string(id)

	database.Rabbit(u.MathExpr)

	json.NewEncoder(w).Encode(u)
}

// Удалить задачу
func DeleteExpression(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u database.Expression
	err := db.QueryRow("SELECT * FROM Expressions WHERE id = $1", id).Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		_, err := db.Exec("DELETE FROM Expressions WHERE id = $1", id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("User deleted")
	}
}

// Удалить все задачи
func DeleteExpressions(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	rows, err := db.Query("SELECT ID FROM Expressions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var u database.Expression
		if err := rows.Scan(&u.ID); err != nil {
			log.Fatal(err)
		}
		_, err := db.Exec("DELETE FROM Expressions WHERE id = $1", u.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
