package customer_db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"

	_ "image/gif"
	_ "image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/nfnt/resize"
	"tendercall.com/main/models"
)

func CreateCustomer(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateCustomer handler called")

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

		fmt.Println("Querying the database for vendors")

		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error parsing multipart form:", err)
			return
		}

		var nextID int
		err = DB.QueryRow("SELECT MAX(id) + 1 FROM customers").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		// Parse customer data from form fields
		customer := models.Customer{
			Name:         r.FormValue("name"),
			Phone_number: r.FormValue("phone_number"),
			Userid:       r.FormValue("userid"),
			Token:        r.FormValue("token"),
			District:     r.FormValue("district"),
			Panchayat:    r.FormValue("panchayat"),
			Join_date:    r.FormValue("join_date"),
			Is_block:     false, // Set default value
		}

		if customer.Name == "" || customer.Phone_number == "" {
			http.Error(w, "Missing required customer fields", http.StatusBadRequest)
			fmt.Println("Missing required customer fields")
			return
		}

		if len(customer.Phone_number) != 10 {
			http.Error(w, "Phone number must be 10 digits", http.StatusBadRequest)
			fmt.Println("Phone number must be 10 digits")
			return
		}

		// r.FormFile("profile_image") retrieves the uploaded file named "profile_image" from the HTTP request.

		file, _, err := r.FormFile("profile_image")
		if err != nil {
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			fmt.Println("Error uploading file:", err)
			return
		}
		defer file.Close()

		// ioutil.ReadAll(file) reads all the content of file into fileBytes, which is a []byte.

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file content", http.StatusInternalServerError)
			fmt.Println("Error reading file content:", err)
			return
		}

		/*
			This checks if the size of fileBytes is greater than 3 MB
			If it is, the image resizing process begins
		*/

		if len(fileBytes) > 3*1024*1024 {

			// Decodes the image from fileBytes using the image package.

			img, _, err := image.Decode(bytes.NewReader(fileBytes))

			if err != nil {
				http.Error(w, "Error decoding image", http.StatusInternalServerError)
				fmt.Println("Error decoding image:", err)
				return
			}

			/*
				resize.Resize resizes img to have a maximum width of 800 pixels (800), while maintaining the aspect ratio (0 for height).
				resize.Lanczos3 is the interpolation method used for resizing, which provides high-quality results.
			*/

			newImage := resize.Resize(800, 0, img, resize.Lanczos3)

			// jpeg.Encode encodes newImage (resized image) into JPEG format and writes it into buf

			var buf bytes.Buffer
			err = jpeg.Encode(&buf, newImage, nil)

			if err != nil {
				http.Error(w, "Error encoding compressed image", http.StatusInternalServerError)
				fmt.Println("Error encoding compressed image:", err)
				return
			}

			/*
			  After resizing and encoding,
			  fileBytes is updated to contain the resized image data (buf.Bytes()).
			*/

			fileBytes = buf.Bytes()
		}

		/*
			Initialize AWS session with credentials and configuration
			aws.Config is a struct that holds the configuration for the AWS SDK for Go.
			session.NewSession initializes a new AWS session with the specified configuration..
			Region specifies the AWS region where AWS service requests are sent.
			Credentials allows specifying AWS credentials used to authenticate requests.
		*/

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				"AKIAYS2NVN4MBSHP33FF",                     // replace with your access key ID
				"aILySGhiQAB7SaFnqozcRZe1MhZ0zNODLof2Alr4", // replace with your secret access key
				""), // optional token, leave blank if not using
		})
		if err != nil {
			log.Printf("Failed to create AWS session: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// s3.New(sess) creates a new S3 service client using the AWS session (sess) created in the previous step.

		svc := s3.New(sess)

		/*
			Upload file to s3
			fmt.Sprintf constructs a unique object key (imageKey) for the image file in the S3 bucket.
			svc.PutObject uploads the image (fileBytes) to the specified S3 bucket ("your-bucket-name") with the specified object key (imageKey).
			bytes.NewReader(fileBytes) sets the content of the object to be uploaded.
		*/

		imageKey := fmt.Sprintf("profile_images/%d.jpg", nextID) // Adjust key as needed
		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String("tendercall-db"),
			Key:    aws.String(imageKey),
			Body:   bytes.NewReader(fileBytes),
		})
		if err != nil {
			log.Printf("Failed to upload image to S3: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Construct and return the URL of the uploaded image

		imageURL := fmt.Sprintf("https://tendercall-db.s3.amazonaws.com/%s", imageKey)

		// Insert customer details into database with image URL
		_, err = DB.Exec("INSERT INTO customers (id, auto_taxi_goods_id, name, phone_number, userid, token, district, panchayat, profile_image, join_date, is_block) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			nextID, customer.Auto_taxi_goods_id, customer.Name, customer.Phone_number, customer.Userid, customer.Token, pq.Array([]string{customer.District}), pq.Array([]string{customer.Panchayat}), imageURL, customer.Join_date, customer.Is_block)
		if err != nil {
			log.Printf("Failed to insert customer: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Println("Customer inserted successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Customer inserted successfully")
	}
}

func GetCustomer(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetCustomer handler called")

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

		fmt.Println("Querying the database for customers")
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

func GetCustomerByUserId(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting Customer by userid\n")

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

		fmt.Println("Querying the database by customers userid:", Userid)
		var customer models.Customer
		err := DB.QueryRow("SELECT * FROM customers WHERE userid = $1", Userid).Scan(&customer.ID, &customer.Auto_taxi_goods_id, &customer.Name, &customer.Phone_number, &customer.Userid, &customer.Token, &customer.District, &customer.Panchayat, &customer.Profile_image, &customer.Join_date, &customer.Is_block)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Customer not found", http.StatusNotFound)
				fmt.Println("Customer not found")
			} else {
				http.Error(w, "Failed to query Customer", http.StatusInternalServerError)
				fmt.Printf("Failed to query Customer: %v\n", err)
			}
			return
		}

		fmt.Println("Customers found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	}
}

func UpdateCustomer(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating Customer by id")

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
			http.Error(w, "Failed to update Customer", http.StatusInternalServerError)
			fmt.Printf("Failed to update Customer: %v\n", err)
			return
		}

		fmt.Println("Customers updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	}
}

func DeleteCustomer(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting Customer by id")

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
			http.Error(w, "Failed to delete customers", http.StatusInternalServerError)
			fmt.Printf("Failed to delete customers: %v\n", err)
			return
		}

		fmt.Println("Customers deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
