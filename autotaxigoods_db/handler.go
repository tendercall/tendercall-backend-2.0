package autotaxigoods_db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"tendercall.com/main/models"
)

func CreateAutoTaxiGoods(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateAutoTaxiGoods handler called")

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

		fmt.Println("Querying the database for AutoTaxiGoods")

		var nextID int
		err := DB.QueryRow("SELECT MAX(id) + 1 FROM auto_taxi_goods").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var AutoTaxiGoods models.AutoTaxiGoods
		err = decoder.Decode(&AutoTaxiGoods)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		if AutoTaxiGoods.Name == "" || AutoTaxiGoods.Phone_number == "" {
			http.Error(w, "Missing required AutoTaxiGoods fields", http.StatusBadRequest)
			fmt.Println("Missing required AutoTaxiGoods fields")
			return
		}

		if len(AutoTaxiGoods.Phone_number) != 10 {
			http.Error(w, "Phone number must be 10 digits", http.StatusBadRequest)
			fmt.Println("Phone number must be 10 digits")
			return
		}

		_, err = DB.Exec("INSERT INTO auto_taxi_goods (id,name,phone_number,unique_id,service_type,experience,district,panchayat,profile_image,image,image_views,rating,created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", nextID, &AutoTaxiGoods.Name, &AutoTaxiGoods.Phone_number, &AutoTaxiGoods.Unique_id, &AutoTaxiGoods.Service_type, &AutoTaxiGoods.Experience, &AutoTaxiGoods.District, &AutoTaxiGoods.Panchayat, &AutoTaxiGoods.Profile_image, &AutoTaxiGoods.Image, &AutoTaxiGoods.Image_views, &AutoTaxiGoods.Rating, &AutoTaxiGoods.Created_at)
		if err != nil {
			log.Printf("Failed to insert AutoTaxiGoods: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Println("AutoTaxiGoods inserted successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "AutoTaxiGoods inserted successfully")
	}
}

func GetAutoTaxiGoods(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetAutoTaxiGoods handler called")

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

		fmt.Println("Querying the database for AutoTaxiGoods")
		rows, err := DB.Query("SELECT id,auto_taxi_goods_id,name,phone_number,userid,token,district,panchayat,profile_image,join_date,is_block FROM customers")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var customers []models.Customer
		for rows.Next() {
			var customer models.Customer
			err := rows.Scan(&customer.ID, &customer.Auto_taxi_goods_id, &customer.Name, &customer.Phone_number, &customer.Userid, &customer.Token, &customer.District, &customer.Panchayat, &customer.Profile_image, &customer.Join_date, &customer.Is_block)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			customers = append(customers, customer)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(customers); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetAutoTaxiGoodsByUserId(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting AutoTaxiGoods by userid\n")

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
		Userid, ok := vars["userid"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		fmt.Println("Querying the database by AutoTaxiGoods userid:", Userid)
		var customer models.Customer
		err := DB.QueryRow("SELECT * FROM customers WHERE userid = $1", Userid).Scan(&customer.ID, &customer.Auto_taxi_goods_id, &customer.Name, &customer.Phone_number, &customer.Userid, &customer.Token, &customer.District, &customer.Panchayat, &customer.Profile_image, &customer.Join_date, &customer.Is_block)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "AutoTaxiGoods not found", http.StatusNotFound)
				fmt.Println("AutoTaxiGoods not found")
			} else {
				http.Error(w, "Failed to query AutoTaxiGoods", http.StatusInternalServerError)
				fmt.Printf("Failed to query AutoTaxiGoods: %v\n", err)
			}
			return
		}

		fmt.Println("AutoTaxiGoods found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	}
}

func UpdateAutoTaxiGoods(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating AutoTaxiGoods by id")

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
		Userid, ok := vars["userid"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		// Decode JSON request body into Vendor struct
		var customer models.Customer
		err := json.NewDecoder(r.Body).Decode(&customer)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the UniqueID from URL is used
		customer.Userid = Userid

		// Update the vendor in the database
		query := `UPDATE customers SET id = $1, auto_taxi_goods_id = $2, name = $3, phone_number = $4, token = $5, district = $6, panchayat =$7, profile_image = $8, join_date = $9, is_block = $10 WHERE userid = $11`
		_, err = DB.Exec(query, customer.ID, customer.Auto_taxi_goods_id, customer.Name, customer.Phone_number, customer.Token, customer.District, customer.Panchayat, customer.Profile_image, customer.Join_date, customer.Is_block, customer.Userid)

		if err != nil {
			http.Error(w, "Failed to update AutoTaxiGoods", http.StatusInternalServerError)
			fmt.Printf("Failed to update AutoTaxiGoods: %v\n", err)
			return
		}

		fmt.Println("AutoTaxiGoods updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	}
}

func DeleteAutoTaxiGoods(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting AutoTaxiGoods by id")

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

		// Extract userid from URL parameters
		vars := mux.Vars(r)
		Userid, ok := vars["userid"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		// Delete the vendor from the database
		_, err := DB.Exec("DELETE FROM customers WHERE userid=$1", Userid)
		if err != nil {
			http.Error(w, "Failed to delete AutoTaxiGoods", http.StatusInternalServerError)
			fmt.Printf("Failed to delete AutoTaxiGoods: %v\n", err)
			return
		}

		fmt.Println("AutoTaxiGoods deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
