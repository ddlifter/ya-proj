package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Expression struct {
	ID       string  `json:"id"`
	MathExpr string  `json:"mathExpr"`
	Status   string  `json:"status"`
	Result   float64 `json:"result"`
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})

}

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func GetExpressions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// get user by id
func GetExpression(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// create user
func AddExpression(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Expression
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO Expressions (MathExpr, Status, Result) VALUES ($1, $2, $3) RETURNING ID", u.MathExpr, u.Status, u.Result).Scan(&u.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

func DeleteExpression(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
