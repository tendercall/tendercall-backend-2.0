package applogs_db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"tendercall.com/main/models"
)

func CreateAppLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateAppLogs handler called")

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

		fmt.Println("Querying the database for AppLogs")

		var nextID int
		err := DB.QueryRow("SELECT MAX(id) + 1 FROM app_logs").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var appLog models.App
		err = decoder.Decode(&appLog)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		currentTime := time.Now()

		// Insert the logs data into the database
		_, err = DB.Exec("INSERT INTO app_logs (id,log_number,log_description,function,created_date,userid,device,platform,estimate_time,exception_case) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", nextID, appLog.Log_Number, appLog.Log_Description, appLog.Function, currentTime, appLog.Userid, appLog.Device, appLog.Platform, appLog.Estimate_Time, appLog.Exceptional_Case)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error inserting new user:", err)
			return
		}

		fmt.Println("AppLogs inserted successfully")
		w.WriteHeader(http.StatusOK)
	}
}

func GetAppLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetAppLogs handler called")

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

		fmt.Println("Querying the database for AppLogs")
		rows, err := DB.Query("SELECT id,log_number,log_description,function,created_date,userid,device,platform,estimate_time,exception_case FROM app_logs")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var appLogs []models.App
		for rows.Next() {
			var appLog models.App
			err := rows.Scan(&appLog.ID, &appLog.Log_Number, &appLog.Log_Description, &appLog.Function, &appLog.CreatedDate, &appLog.Userid, &appLog.Device, &appLog.Platform, &appLog.Estimate_Time, &appLog.Exceptional_Case)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			appLogs = append(appLogs, appLog)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(appLogs); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetAppLogsById(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting AppLogs by userid\n")

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

		// Extracting userid from the URL parameters
		vars := mux.Vars(r)
		Userid, ok := vars["userid"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		fmt.Println("Querying the database by AppLogs userid:", Userid)
		var appLog models.App
		err := DB.QueryRow("SELECT * FROM app_logs WHERE userid = $1", Userid).Scan(&appLog.ID, &appLog.Log_Number, &appLog.Log_Description, &appLog.Function, &appLog.CreatedDate, &appLog.Userid, &appLog.Device, &appLog.Platform, &appLog.Estimate_Time, &appLog.Exceptional_Case)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Logs not found", http.StatusNotFound)
				fmt.Println("Logs not found")
			} else {
				http.Error(w, "Failed to query AppLogs", http.StatusInternalServerError)
				fmt.Printf("Failed to query AppLogs: %v\n", err)
			}
			return
		}

		fmt.Println("AppLogs found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appLog)
	}
}

func UpdateAppLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating AppLogs by Userid")

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

		// Extract Userid from URL parameters
		vars := mux.Vars(r)
		Userid, ok := vars["userid"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		// Decode JSON request body into AppLogs struct
		var appLog models.App
		err := json.NewDecoder(r.Body).Decode(&appLog)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the Userid from URL is used
		appLog.Userid = Userid

		// Update the AppLogs in the database
		query := `UPDATE app_logs SET created_date=$1, log_number=$2, log_description=$3, function=$4, userid=$5, device=$6, platform=$7, estimate_time=$8, exception_case=$9 WHERE id=$10`
		_, err = DB.Exec(query, appLog.CreatedDate, appLog.Log_Number, appLog.Log_Description, appLog.Function, appLog.Userid, appLog.Device, appLog.Platform, appLog.Estimate_Time, appLog.Exceptional_Case, appLog.ID)

		if err != nil {
			http.Error(w, "Failed to update AppLogs", http.StatusInternalServerError)
			fmt.Printf("Failed to update AppLogs: %v\n", err)
			return
		}

		fmt.Println("AppLogs updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appLog)
	}
}

func DeleteAppLogs(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting AppLogs by Userid")

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
		// Extract Userid from URL parameters
		vars := mux.Vars(r)
		Userid, ok := vars["userid"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing Userid in request")
			return
		}

		// Delete the AppLogs from the database
		_, err := DB.Exec("DELETE FROM app_logs WHERE id=$1", Userid)
		if err != nil {
			http.Error(w, "Failed to delete Logs", http.StatusInternalServerError)
			fmt.Printf("Failed to delete Logs: %v\n", err)
			return
		}

		fmt.Println("AppLogs deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
