package banner_db

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

func CreateBanner(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateBanner handler called")

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

		fmt.Println("Querying the database for next banner ID")

		var nextID int
		err := DB.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM banners").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var banner models.Banner
		err = decoder.Decode(&banner)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		currentTime := time.Now()

		_, err = DB.Exec("INSERT INTO banners (id, image, district, panchayat, services, service_category, createdat) VALUES ($1, $2, $3, $4, $5, $6, $7)", nextID, banner.Image, banner.District, banner.Panchayat, banner.Services, banner.Service_category, currentTime)
		if err != nil {
			log.Printf("Failed to insert banner: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Println("Banner inserted successfully")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Banner inserted successfully")
	}
}

func GetBanner(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetBanner handler called")

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

		fmt.Println("Querying the database for Banner")
		rows, err := DB.Query("SELECT id,image,district,panchayat,services,service_category,createdat FROM banners")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		defer rows.Close()

		var banners []models.Banner
		for rows.Next() {
			var Banner models.Banner
			err := rows.Scan(&Banner.ID, &Banner.Image, &Banner.District, &Banner.Panchayat, &Banner.Services, &Banner.Service_category, &Banner.Created_at)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			banners = append(banners, Banner)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(banners); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetBannerById(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting Banner by userid\n")

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

		fmt.Println("Querying the database by Banner id:", ID)
		var Banner models.Banner
		err := DB.QueryRow("SELECT * FROM banners WHERE id = $1", ID).Scan(&Banner.ID, &Banner.Image, &Banner.District, &Banner.Panchayat, &Banner.Services, &Banner.Service_category, &Banner.Created_at)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Banner not found", http.StatusNotFound)
				fmt.Println("Banner not found")
			} else {
				http.Error(w, "Failed to query Banner", http.StatusInternalServerError)
				fmt.Printf("Failed to query Banner: %v\n", err)
			}
			return
		}

		fmt.Println("Banner found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Banner)
	}
}

func UpdateBanner(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating Banner by id")

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
			fmt.Println("Missing id in request")
			return
		}

		// Convert ID from string to int
		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Invalid id format:", err)
			return
		}

		// Decode JSON request body into Banner struct
		var Banner models.Banner
		err = json.NewDecoder(r.Body).Decode(&Banner)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the ID from URL is used
		Banner.ID = uint(ID)

		// Update the vendor in the database
		query := `UPDATE banners SET image = $1, district = $2, panchayat = $3, services = $4, service_category = $5, createdat = $6 WHERE id = $7`
		_, err = DB.Exec(query, Banner.Image, Banner.District, Banner.Panchayat, Banner.Services, Banner.Service_category, Banner.Created_at, Banner.ID)

		if err != nil {
			http.Error(w, "Failed to update Banner", http.StatusInternalServerError)
			fmt.Printf("Failed to update Banner: %v\n", err)
			return
		}

		fmt.Println("Banner updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Banner)
	}
}

func DeleteBanner(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting Banner by id")

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
		ID, ok := vars["id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		// Delete the vendor from the database
		_, err := DB.Exec("DELETE FROM banners WHERE id=$1", ID)
		if err != nil {
			http.Error(w, "Failed to delete Banner", http.StatusInternalServerError)
			fmt.Printf("Failed to delete Banner: %v\n", err)
			return
		}

		fmt.Println("Banner deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
