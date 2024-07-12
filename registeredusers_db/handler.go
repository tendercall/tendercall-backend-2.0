package registeredusers_db

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

func CreateRegisteredUser(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateRegisteredUser handler called")

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

		fmt.Println("Querying the database for RegisteredUser")

		var nextID int
		err := DB.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM registered_users").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var user models.User
		err = decoder.Decode(&user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		if user.Name == "" || user.Phone_number == "" {
			http.Error(w, "Missing required RegisteredUser fields", http.StatusBadRequest)
			fmt.Println("Missing required RegisteredUser fields")
			return
		}

		if len(user.Phone_number) != 10 {
			http.Error(w, "Phone number must be 10 digits", http.StatusBadRequest)
			fmt.Println("Phone number must be 10 digits")
			return
		}

		var exists bool
		var isVendor bool
		err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM registered_users WHERE phone_number = $1), COALESCE((SELECT is_vendor FROM registered_users WHERE phone_number = $1), false) FROM registered_users WHERE phone_number = $1", user.Phone_number).Scan(&exists, &isVendor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error checking for existing user:", err)
			return
		}

		if exists {
			response := map[string]interface{}{
				"message":   "Phone number already exists",
				"is_vendor": isVendor,
			}
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			fmt.Println("Phone number already exists")
			return
		}

		currentTime := time.Now()
		_, err = DB.Exec("INSERT INTO registered_users(id, name, phone_number, created_at, is_vendor) VALUES ($1, $2, $3, $4, $5)", nextID, user.Name, user.Phone_number, currentTime, user.Is_vendor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error inserting new user:", err)
			return
		}

		fmt.Println("RegisteredUser inserted successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "RegisteredUser inserted successfully"})
	}
}

func GetRegisteredUser(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetRegisteredUser handler called")

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

		fmt.Println("Querying the database for RegisteredUser")
		rows, err := DB.Query("SELECT id,name,phone_number,created_at,is_vendor FROM registered_users")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var Users []models.User
		for rows.Next() {
			var User models.User
			err := rows.Scan(&User.ID, &User.Name, &User.Phone_number, &User.CreatedAt, &User.Is_vendor)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			Users = append(Users, User)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(Users); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetRegisteredUserById(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting RegisteredUser by userid\n")

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
			fmt.Println("Missing userid in request")
			return
		}

		fmt.Println("Querying the database by RegisteredUser userid:", ID)
		var User models.User
		err := DB.QueryRow("SELECT * FROM registered_users WHERE id = $1", ID).Scan(&User.ID, &User.Name, &User.Phone_number, &User.CreatedAt, &User.Is_vendor)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "RegisteredUser not found", http.StatusNotFound)
				fmt.Println("RegisteredUser not found")
			} else {
				http.Error(w, "Failed to query RegisteredUser", http.StatusInternalServerError)
				fmt.Printf("Failed to query RegisteredUser: %v\n", err)
			}
			return
		}

		fmt.Println("RegisteredUser found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(User)
	}
}

func UpdateRegisteredUser(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating RegisteredUser by id")

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
		var User models.User
		err = json.NewDecoder(r.Body).Decode(&User)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the ID from URL is used
		User.ID = uint(ID)

		// Update the vendor in the database
		query := `UPDATE registered_users SET name=$1, phone_number=$2, is_vendor=$3 WHERE id=$4`
		_, err = DB.Exec(query, User.Name, User.Phone_number, User.Is_vendor, User.ID)

		if err != nil {
			http.Error(w, "Failed to update RegisteredUser", http.StatusInternalServerError)
			fmt.Printf("Failed to update RegisteredUser: %v\n", err)
			return
		}

		fmt.Println("RegisteredUser updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(User)
	}
}

func DeleteRegisteredUser(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting RegisteredUser by id")

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
		_, err := DB.Exec("DELETE FROM registered_users WHERE id=$1", ID)
		if err != nil {
			http.Error(w, "Failed to delete RegisteredUser", http.StatusInternalServerError)
			fmt.Printf("Failed to delete RegisteredUser: %v\n", err)
			return
		}

		fmt.Println("RegisteredUser deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
