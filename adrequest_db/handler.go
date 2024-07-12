package adrequest_db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"tendercall.com/main/models"
)

func CreateAdRequest(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateAdRequest handler called")

		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")

		if tokenString != "Bearer eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcxNTU4Njc4MywiaWF0IjoxNzE1NTg2NzgzfQ.f3OxHxEJ-IX2D3f98VliSurFKWKh3GI5Mh3yGwsS16E" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println("Unauthorized request")
			return
		}

		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method:", r.Method)
			return
		}

		fmt.Println("Querying the database for AdRequest")
		var nextID int
		err := DB.QueryRow("SELECT MAX(id) + 1 FROM ad_request").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var AdRequest models.AdRequest
		err = decoder.Decode(&AdRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Insert the ad_request data into the database
		_, err = DB.Exec("INSERT INTO ad_request (id,ad_category,phone_number,unique_id,category,district,start_date,request_date) VALUES ($1, $2, $3, $4, $5, $6, $7, 48)", nextID, AdRequest.Ad_category, AdRequest.Phone_number, AdRequest.Unique_id, AdRequest.Category, AdRequest.District, AdRequest.Start_date, AdRequest.Request_date)
		if err != nil {
			log.Printf("Failed to insert AdRequest: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Println("AdRequest inserted successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "AdRequest inserted successfully")
	}
}

func GetAdRequest(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetAdRequest handler called")

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

		fmt.Println("Querying the database for AdRequest")
		rows, err := DB.Query("SELECT id,ad_category,phone_number,unique_id,category,district,start_date,request_date FROM ad_request")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var AdRequests []models.AdRequest
		for rows.Next() {
			var AdRequest models.AdRequest
			err := rows.Scan(&AdRequest.ID, &AdRequest.Ad_category, &AdRequest.Phone_number, &AdRequest.Unique_id, &AdRequest.Category, &AdRequest.District, &AdRequest.Start_date, &AdRequest.Request_date)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			AdRequests = append(AdRequests, AdRequest)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(AdRequests); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetAdRequestById(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting AdRequest by id:")

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

		// Extracting UniqueID from the URL parameters
		vars := mux.Vars(r)
		Unique_id, ok := vars["unique_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing Unique_id in request")
			return
		}

		fmt.Println("Querying the database by AdRequest Unique_id:", Unique_id)
		var AdRequest models.AdRequest
		err := DB.QueryRow("SELECT * FROM ad_request WHERE unique_id = $1", Unique_id).Scan(&AdRequest.ID, &AdRequest.Phone_number, &AdRequest.Ad_category, &AdRequest.Unique_id, &AdRequest.Category, &AdRequest.District, &AdRequest.Start_date, &AdRequest.Request_date)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "AdRequest not found", http.StatusNotFound)
				fmt.Println("AdRequest not found")
			} else {
				http.Error(w, "Failed to query AdRequest", http.StatusInternalServerError)
				fmt.Printf("Failed to query AdRequest: %v\n", err)
			}
			return
		}

		fmt.Println("AdRequest found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AdRequest)
	}
}

func UpdateAdRequest(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating AdRequest by id")

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

		// Extract UniqueID from URL parameters
		vars := mux.Vars(r)
		Unique_id, ok := vars["Unique_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing Unique_id in request")
			return
		}

		// Decode JSON request body into Vendor struct
		var AdRequests models.AdRequest
		err := json.NewDecoder(r.Body).Decode(&AdRequests)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the UniqueID from URL is used
		AdRequests.Unique_id = Unique_id

		// Update the vendor in the database
		query := `UPDATE ad_request SET ad_category = $1, phone-number= $2, unique_id = $3, category = $4, district = $5, start_date = $6, request_date = $7 WHERE id = $8`
		_, err = DB.Exec(query, AdRequests.ID, AdRequests.Ad_category, AdRequests.Phone_number, AdRequests.Category, AdRequests.District, AdRequests.Start_date, AdRequests.Request_date, AdRequests.Unique_id)

		if err != nil {
			http.Error(w, "Failed to update AdRequest", http.StatusInternalServerError)
			fmt.Printf("Failed to update AdRequest: %v\n", err)
			return
		}

		fmt.Println("AdRequest updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AdRequests)
	}
}

func DeleteAdRequest(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting AdRequest by id")

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

		// Extract UniqueID from URL parameters
		vars := mux.Vars(r)
		Unique_id, ok := vars["Unique_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing Unique_id in request")
			return
		}

		// Delete the vendor from the database
		_, err := DB.Exec("SELECT id FROM ad_request WHERE id = $1", Unique_id)
		if err != nil {
			http.Error(w, "Failed to delete AdRequest", http.StatusInternalServerError)
			fmt.Printf("Failed to delete AdRequest: %v\n", err)
			return
		}

		fmt.Println("AdRequest deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
