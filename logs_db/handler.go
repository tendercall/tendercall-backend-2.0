package logs_db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"tendercall.com/main/models"
)

func CreateLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateLogs handler called")

		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")

		if tokenString != "Bearer eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcxNTU4Njc4MywiaWF0IjoxNzE1NTg2NzgzfQ.f3OxHxEJ-IX2D3f98VliSurFKWKh3GI5Mh3yGwsS16E" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println("Unauthorized request")
			return
		}

		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method:", r.Method)
			return
		}

		fmt.Println("Querying the database for Logs")

		var nextID int
		err := DB.QueryRow("SELECT MAX(id) + 1 FROM logs").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var Log models.Log
		err = decoder.Decode(&Log)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		currentTime := time.Now()

		// Insert the logs data into the database
		_, err = DB.Exec("INSERT INTO logs (id,created_date,log_number,log_description,function,estimate_time) VALUES ($1, $2, $3, $4, $5, $6)", nextID, currentTime, Log.Log_Number, Log.Log_Description, Log.Function, Log.Estimate_Time)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error inserting new user:", err)
			return
		}

		fmt.Println("Logs inserted successfully")
		w.WriteHeader(http.StatusOK)
	}
}

func GetLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetLogs handler called")

		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")

		if tokenString != "Bearer eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcxNTU4Njc4MywiaWF0IjoxNzE1NTg2NzgzfQ.f3OxHxEJ-IX2D3f98VliSurFKWKh3GI5Mh3yGwsS16E" {
			http.Error(w, "StatusUnauthorized", http.StatusUnauthorized)
			fmt.Println("Unauthorized request")
			return
		}

		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method:", r.Method)
			return
		}

		fmt.Println("Querying the database for Logs")
		rows, err := DB.Query("SELECT id,created_date,log_number,log_description,function,estimate_time FROM logs")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var Logs []models.Log
		for rows.Next() {
			var Log models.Log
			err := rows.Scan(&Log.ID, &Log.CreatedDate, &Log.Log_Number, &Log.Log_Description, &Log.Function, &Log.Estimate_Time)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			Logs = append(Logs, Log)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(Logs); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetLogsById(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting Logs by userid\n")

		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")

		if tokenString != "Bearer eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcxNTU4Njc4MywiaWF0IjoxNzE1NTg2NzgzfQ.f3OxHxEJ-IX2D3f98VliSurFKWKh3GI5Mh3yGwsS16E" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println("Unauthorized request")
			return
		}

		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method")
			return
		}

		// Extracting id from the URL parameters
		vars := mux.Vars(r)
		ID, ok := vars["id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing id in request")
			return
		}

		fmt.Println("Querying the database by Logs userid:", ID)
		var Log models.Log
		err := DB.QueryRow("SELECT * FROM logs WHERE id = $1", ID).Scan(&Log.ID, &Log.CreatedDate, &Log.Log_Number, &Log.Log_Description, &Log.Function, &Log.Estimate_Time)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Logs not found", http.StatusNotFound)
				fmt.Println("Logs not found")
			} else {
				http.Error(w, "Failed to query Logs", http.StatusInternalServerError)
				fmt.Printf("Failed to query Logs: %v\n", err)
			}
			return
		}

		fmt.Println("Logs found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Log)
	}
}

func UpdateLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating Logs by id")

		// Check authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString != "Bearer eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcxNTU4Njc4MywiaWF0IjoxNzE1NTg2NzgzfQ.f3OxHxEJ-IX2D3f98VliSurFKWKh3GI5Mh3yGwsS16E" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println("Unauthorized request")
			return
		}

		// Check request method
		if r.Method != "PUT" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method:", r.Method)
			return
		}

		// Extract ID from URL parameters
		vars := mux.Vars(r)
		IDStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		// Convert ID from string to int
		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Invalid id format:", err)
			return
		}

		// Decode JSON request body into Vendor struct
		var Log models.Log
		err = json.NewDecoder(r.Body).Decode(&Log)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the ID from URL is used
		Log.ID = uint(ID)

		// Update the vendor in the database
		query := `UPDATE logs SET created_date=$1, log_number=$2, log_description=$3, function=$4, estimate_time=$5 WHERE id=$6`
		_, err = DB.Exec(query, Log.CreatedDate, Log.Log_Number, Log.Log_Description, Log.Function, Log.Estimate_Time, Log.ID)

		if err != nil {
			http.Error(w, "Failed to update RegisteredUser", http.StatusInternalServerError)
			fmt.Printf("Failed to update RegisteredUser: %v\n", err)
			return
		}

		fmt.Println("Logs updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Log)
	}
}

func DeleteLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting Logs by id")

		// Check authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString != "Bearer eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcxNTU4Njc4MywiaWF0IjoxNzE1NTg2NzgzfQ.f3OxHxEJ-IX2D3f98VliSurFKWKh3GI5Mh3yGwsS16E" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println("Unauthorized request")
			return
		}

		// Check request method
		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method:", r.Method)
			return
		}
		// Extract ID from URL parameters
		vars := mux.Vars(r)
		ID, ok := vars["id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing unique_id in request")
			return
		}

		// Delete the vendor from the database
		_, err := DB.Exec("DELETE FROM logs WHERE id=$1", ID)
		if err != nil {
			http.Error(w, "Failed to delete Logs", http.StatusInternalServerError)
			fmt.Printf("Failed to delete Logs: %v\n", err)
			return
		}

		fmt.Println("Logs deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
