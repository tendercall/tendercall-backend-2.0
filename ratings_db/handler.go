package ratings_db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"tendercall.com/main/models"
)

func CreateRating(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateRating handler called")

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

		fmt.Println("Querying the database for Rating")
		var nextID int
		err := DB.QueryRow("SELECT MAX(id) + 1 FROM ratings").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var Rating models.Ratings
		err = decoder.Decode(&Rating)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		currentTime := time.Now()

		// Insert the ad_request data into the database
		_, err = DB.Exec("INSERT INTO ratings (id,created_at,rating,vendor_id,remark) VALUES ($1, $2, $3, $4, $5)", nextID, currentTime, Rating.Rating, Rating.Vendor_id, Rating.Remark)
		if err != nil {
			log.Printf("Failed to insert rating: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Println("Rating inserted successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Rating inserted successfully")
	}
}

func GetRating(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetRating handler called")

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

		fmt.Println("Querying the database for ratings")
		rows, err := DB.Query("SELECT id,created_at,ratings,vendor_id,remark FROM ratings")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var ratings []models.Ratings
		for rows.Next() {
			var Rating models.Ratings
			err := rows.Scan(&Rating.ID, &Rating.Created_at, &Rating.Rating, &Rating.Vendor_id, &Rating.Remark)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			ratings = append(ratings, Rating)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ratings); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetRatingById(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting Rating by id")

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
		ID, ok := vars["id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing id in request")
			return
		}

		fmt.Println("Querying the database by Rating id:", ID)
		var Rating models.Ratings
		err := DB.QueryRow("SELECT * FROM ratings WHERE id = $1", ID).Scan(&Rating.ID, &Rating.Created_at, &Rating.Rating, &Rating.Vendor_id, &Rating.Remark)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Rating not found", http.StatusNotFound)
				fmt.Println("Rating not found")
			} else {
				http.Error(w, "Failed to query ratings", http.StatusInternalServerError)
				fmt.Printf("Failed to query ratings: %v\n", err)
			}
			return
		}

		fmt.Println("Rating found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Rating)
	}
}

func UpdateRating(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating Rating by id")

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
		IDStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing unique_id in request")
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
		var Rating models.Ratings
		err = json.NewDecoder(r.Body).Decode(&Rating)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the UniqueID from URL is used
		Rating.ID = uint(ID)

		// Update the vendor in the database
		query := `UPDATE ratings SET created_at = $1, rating = $2, vendor_id = $3, remark = $4 WHERE id = $5`
		_, err = DB.Exec(query, Rating.Created_at, Rating.Rating, Rating.Vendor_id, Rating.Remark, Rating.ID)

		if err != nil {
			http.Error(w, "Failed to update Rating", http.StatusInternalServerError)
			fmt.Printf("Failed to update Rating: %v\n", err)
			return
		}

		fmt.Println("Rating updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Rating)
	}
}

func DeleteRating(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting Rating by id")

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
		_, err := DB.Exec("SELECT id FROM ratings WHERE id = $1", ID)
		if err != nil {
			http.Error(w, "Failed to delete Rating", http.StatusInternalServerError)
			fmt.Printf("Failed to delete Rating: %v\n", err)
			return
		}

		fmt.Println("Rating deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
