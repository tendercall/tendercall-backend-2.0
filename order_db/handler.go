package order_db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"tendercall.com/main/models"
)

func CreateOrder(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateOrder handler called")

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

		fmt.Println("Querying the database for Orders")

		var nextID int
		err := DB.QueryRow("SELECT MAX(id) + 1 FROM orders").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var Order models.Orders
		err = decoder.Decode(&Order)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		if Order.Name == "" || Order.Phone_number == "" {
			http.Error(w, "Missing required Order fields", http.StatusBadRequest)
			fmt.Println("Missing required Order fields")
			return
		}

		if len(Order.Phone_number) != 10 {
			http.Error(w, "Phone number must be 10 digits", http.StatusBadRequest)
			fmt.Println("Phone number must be 10 digits")
			return
		}

		_, err = DB.Exec("INSERT INTO orders (id,order_id,name,address,category,created_at,description,assigned_vendors,end_date,image,service,sqf,start_date,user_id,phone_number,population,function_type,ac_available,food_available,event_date,program_type,travel_experience,origin,destination,vehicle_type,dining,accommodation,property_type,budget,property_location,Area,quantity,tool_type,seat_capacity,rent_period,service_type,type,product_type,course_type,called_vendors,building_type,isRating_enable,property_address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43)", nextID, &Order.Order_id, &Order.Name, &Order.Address, &Order.Category, &Order.Created_at, &Order.Description, &Order.Assigned_vendors, &Order.End_date, &Order.Image, &Order.Service, &Order.Sqf, &Order.Start_date, &Order.User_id, &Order.Phone_number, &Order.Population, &Order.Function_type, &Order.Ac_available, &Order.Food_available, &Order.Event_date, &Order.Program_type, &Order.Travel_experience, &Order.Origin, &Order.Destination, &Order.Vehicle_type, &Order.Dining, &Order.Accommodation, &Order.Property_type, &Order.Budget, &Order.Property_location, &Order.Area, &Order.Quantity, &Order.Tool_type, &Order.Seat_capacity, &Order.Rent_period, &Order.Service_type, &Order.Type, &Order.Product_type, &Order.Course_type, &Order.Called_vendors, &Order.Building_type, &Order.IsRating_Enable, &Order.Property_address)
		if err != nil {
			log.Printf("Failed to insert Order: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Println("Order inserted successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Order inserted successfully")
	}
}

func GetOrder(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetOrder handler called")

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

		fmt.Println("Querying the database for Orders")
		rows, err := DB.Query("SELECT id,order_id,name,address,category,created_at,description,assigned_vendors,end_date,image,service,sqf,start_date,user_id,phone_number,population,function_type,ac_available,food_available,event_date,program_type,travel_experience,origin,destination,vehicle_type,dining,accommodation,property_type,budget,property_location,Area,quantity,tool_type,seat_capacity,rent_period,service_type,type,product_type,course_type,called_vendors,building_type,isRating_enable,property_address FROM orders")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var Orders []models.Orders
		for rows.Next() {
			var Order models.Orders
			err := rows.Scan(&Order.ID, &Order.Order_id, &Order.Name, &Order.Address, &Order.Category, &Order.Created_at, &Order.Description, &Order.Assigned_vendors, &Order.End_date, &Order.Image, &Order.Service, &Order.Sqf, &Order.Start_date, &Order.User_id, &Order.Phone_number, &Order.Population, &Order.Function_type, &Order.Ac_available, &Order.Food_available, &Order.Event_date, &Order.Program_type, &Order.Travel_experience, &Order.Origin, &Order.Destination, &Order.Vehicle_type, &Order.Dining, &Order.Accommodation, &Order.Property_type, &Order.Budget, &Order.Property_location, &Order.Area, &Order.Quantity, &Order.Tool_type, &Order.Seat_capacity, &Order.Rent_period, &Order.Service_type, &Order.Type, &Order.Product_type, &Order.Course_type, &Order.Called_vendors, &Order.Building_type, &Order.IsRating_Enable, &Order.Property_address)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			Orders = append(Orders, Order)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(Orders); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetOrderByUserId(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting Order by userid\n")

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

		// Extracting User_id from the URL parameters
		vars := mux.Vars(r)
		User_id, ok := vars["user_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		fmt.Println("Querying the database by Orders user_id:", User_id)
		var Order models.Orders
		err := DB.QueryRow("SELECT * FROM orders WHERE user_id = $1", User_id).Scan(&Order.ID, &Order.Order_id, &Order.Name, &Order.Address, &Order.Category, &Order.Created_at, &Order.Description, &Order.Assigned_vendors, &Order.End_date, &Order.Image, &Order.Service, &Order.Sqf, &Order.Start_date, &Order.User_id, &Order.Phone_number, &Order.Population, &Order.Function_type, &Order.Ac_available, &Order.Food_available, &Order.Event_date, &Order.Program_type, &Order.Travel_experience, &Order.Origin, &Order.Destination, &Order.Vehicle_type, &Order.Dining, &Order.Accommodation, &Order.Property_type, &Order.Budget, &Order.Property_location, &Order.Area, &Order.Quantity, &Order.Tool_type, &Order.Seat_capacity, &Order.Rent_period, &Order.Service_type, &Order.Type, &Order.Product_type, &Order.Course_type, &Order.Called_vendors, &Order.Building_type, &Order.IsRating_Enable, &Order.Property_address)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Order not found", http.StatusNotFound)
				fmt.Println("Order not found")
			} else {
				http.Error(w, "Failed to query Order", http.StatusInternalServerError)
				fmt.Printf("Failed to query Order: %v\n", err)
			}
			return
		}

		fmt.Println("Order found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Order)
	}
}

func UpdateOrder(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating Order by user_id")

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
		User_id, ok := vars["user_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		// Decode JSON request body into Vendor struct
		var Order models.Orders
		err := json.NewDecoder(r.Body).Decode(&Order)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the Userid from URL is used
		Order.User_id = User_id

		// Update the vendor in the database
		query := `UPDATE orders SET order_id = $1, name = $2, address = $3, category = $4, created_at = $5, description = $6, assigned_vendors = $7, end_date = $8, image = $9, service = $10, sqf = $11, start_date = $12, user_id = $13, phone_number = $14, population = $15, function_type = $16, ac_available = $17, food_available = $18, event_date = $19, program_type = $20, travel_experience = $21, origin = $22, destination = $23, vehicle_type = $24, dining = $25, accommodation = $26, property_type = $27, budget = $28, property_location = $29, Area = $30, quantity = $31, tool_type = $32, seat_capacity = $33, rent_period = $34, service_type = $35, type = $36, product_type = $37, course_type = $38, called_vendors = $39, building_type = $40, isRating_enable = $41, property_address = $42 WHERE id = $43`
		_, err = DB.Exec(query, Order.ID, Order.Order_id, Order.Name, Order.Address, Order.Category, Order.Created_at, Order.Description, Order.Assigned_vendors, Order.End_date, Order.Image, Order.Service, Order.Sqf, Order.Start_date, Order.Phone_number, Order.Population, Order.Function_type, Order.Ac_available, Order.Food_available, Order.Event_date, Order.Program_type, Order.Travel_experience, Order.Origin, Order.Destination, Order.Vehicle_type, Order.Dining, Order.Accommodation, Order.Property_type, Order.Budget, Order.Property_location, Order.Area, Order.Quantity, Order.Tool_type, Order.Seat_capacity, Order.Rent_period, Order.Service_type, Order.Type, Order.Product_type, Order.Course_type, Order.Called_vendors, Order.Building_type, Order.IsRating_Enable, Order.Property_address, Order.User_id)

		if err != nil {
			http.Error(w, "Failed to update Order", http.StatusInternalServerError)
			fmt.Printf("Failed to update Order: %v\n", err)
			return
		}

		fmt.Println("Order updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Order)
	}
}

func DeleteOrder(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting Order by id")

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
		User_id, ok := vars["user_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing userid in request")
			return
		}

		// Delete the vendor from the database
		_, err := DB.Exec("DELETE FROM orders WHERE user_id=$1", User_id)
		if err != nil {
			http.Error(w, "Failed to delete Order", http.StatusInternalServerError)
			fmt.Printf("Failed to delete Order: %v\n", err)
			return
		}

		fmt.Println("Order deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
