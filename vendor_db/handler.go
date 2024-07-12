package vendor_db

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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/nfnt/resize"
	"tendercall.com/main/models"
)

func CreateVendor(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateVendor handler called")

		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")

		if tokenString != "Bearer eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcxNTU4Njc4MywiaWF0IjoxNzE1NTg2NzgzfQ.f3OxHxEJ-IX2D3f98VliSurFKWKh3GI5Mh3yGwsS16E" { // Update the token validation
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println("Unauthorized request")
			return
		}

		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method:", r.Method)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error parsing multipart form:", err)
			return
		}

		fmt.Println("Querying the database for vendors")

		var nextID int
		err = DB.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM vendors").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		// Parse vendor data from form fields
		vendor := models.Vendor{
			Name:            r.FormValue("name"),
			BusinessName:    r.FormValue("business_name"),
			BusinessAddress: r.FormValue("business_address"),
			PhoneNumber:     r.FormValue("phone_number"),
			UniqueID:        r.FormValue("unique_id"),
			Token:           r.FormValue("token"),
			District:        pq.StringArray(r.Form["district"]),
			Panchayat:       pq.StringArray(r.Form["panchayat"]),
			Services:        r.FormValue("services"),
			ServiceCategory: r.FormValue("service_category"),
			Experience:      r.FormValue("experience"),
			Reference:       r.FormValue("reference"),
			Premium:         r.FormValue("premium"),
			PlanPurchase:    r.FormValue("plan_purchase"),
			ImageView:       r.FormValue("image_view"),
			Favorites:       r.FormValue("favorites"),
			CreatedAt:       r.FormValue("created_at"),
			Latitude:        r.FormValue("latitude"),
			Longitude:       r.FormValue("longitude"),
		}

		if vendor.Name == "" || vendor.BusinessName == "" || vendor.BusinessAddress == "" || vendor.PhoneNumber == "" {
			http.Error(w, "Missing required vendor fields", http.StatusBadRequest)
			fmt.Println("Missing required vendor fields")
			return
		}

		if len(vendor.PhoneNumber) != 10 {
			http.Error(w, "Phone number must be 10 digits", http.StatusBadRequest)
			fmt.Println("Phone number must be 10 digits")
			return
		}

		/*
			an anonymous function uploadImage that takes fileKey as a parameter and returns a string and an error.
			fileKey is used to retrieve the uploaded file from the request form data.
			r.FormFile(fileKey) retrieves the uploaded file identified by fileKey from the HTTP request.
		*/

		uploadImage := func(fileKey string) (string, error) {
			file, _, err := r.FormFile(fileKey)
			if err != nil {
				return "", fmt.Errorf("error uploading file: %v", err)
			}
			defer file.Close()

			// ioutil.ReadAll(file) reads all the content of file into fileBytes, which is a []byte.

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				return "", fmt.Errorf("error reading file content: %v", err)
			}

			/*
				This checks if the size of fileBytes is greater than 3 MB
				If it is, the image resizing process begins
			*/

			if len(fileBytes) > 3*1024*1024 {

				// Decodes the image from fileBytes using the image package.

				img, _, err := image.Decode(bytes.NewReader(fileBytes))
				if err != nil {
					return "", fmt.Errorf("error decoding image: %v", err)
				}

				/*
					resize.Resize resizes image to have a maximum width of 800 pixels (800), while maintaining the aspect ratio (0 for height).
					resize.Lanczos3 is the interpolation method used for resizing, which provides high-quality results.
				*/

				newImage := resize.Resize(800, 0, img, resize.Lanczos3)

				// jpeg.Encode encodes newImage (resized image) into JPEG format and writes it into buf

				var buf bytes.Buffer
				err = jpeg.Encode(&buf, newImage, nil)
				if err != nil {
					return "", fmt.Errorf("error encoding compressed image: %v", err)
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
				return "", fmt.Errorf("failed to create AWS session: %v", err)
			}

			// s3.New(sess) creates a new S3 service client using the AWS session (sess) created in the previous step.

			svc := s3.New(sess)

			/*
				fmt.Sprintf constructs a unique object key (imageKey) for the image file in the S3 bucket.
				svc.PutObject uploads the image (fileBytes) to the specified S3 bucket ("your-bucket-name") with the specified object key (imageKey).
				bytes.NewReader(fileBytes) sets the content of the object to be uploaded.
			*/

			imageKey := fmt.Sprintf("images/%d_%s.jpg", nextID, fileKey)
			_, err = svc.PutObject(&s3.PutObjectInput{
				Bucket: aws.String("tendercall-db"),
				Key:    aws.String(imageKey),
				Body:   bytes.NewReader(fileBytes),
			})
			if err != nil {
				return "", fmt.Errorf("failed to upload image to S3: %v", err)
			}

			// Construct and return the URL of the uploaded image

			return fmt.Sprintf("https://tendercall-db.s3.amazonaws.com/%s", imageKey), nil
		}

		// Upload ProfileImage
		profileImageURL, err := uploadImage("profile_image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		vendor.ProfileImage = profileImageURL

		// Upload Image
		imageURL, err := uploadImage("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		vendor.Image = imageURL

		// Insert vendor details into database with image URLs
		_, err = DB.Exec("INSERT INTO vendors (id, name, business_name, business_address, phone_number, unique_id, token, district, panchayat, services, service_category, experience, reference, premium, profile_image, plan_purchase, image, image_view, favorites, created_at, latitude, longitude) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)", nextID, vendor.Name, vendor.BusinessName, vendor.BusinessAddress, vendor.PhoneNumber, vendor.UniqueID, vendor.Token, vendor.District, vendor.Panchayat, pq.Array([]string{vendor.Services}), pq.Array([]string{vendor.ServiceCategory}), vendor.Experience, vendor.Reference, vendor.Premium, vendor.ProfileImage, vendor.PlanPurchase, vendor.Image, pq.Array([]string{vendor.ImageView}), pq.Array([]string{vendor.Favorites}), vendor.CreatedAt, vendor.Latitude, vendor.Longitude)
		if err != nil {
			log.Printf("Failed to insert vendor: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Println("Vendor inserted successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Vendor inserted successfully")
	}
}

func GetVendor(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("GetVendor handler called\n")

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

		fmt.Println("Querying the database for vendors")
		rows, err := DB.Query("SELECT id, name, business_name, business_address, phone_number, unique_id, token, district, panchayat, services, service_category, experience, reference, premium, profile_image, plan_purchase, image, image_view, favorites, created_at, latitude, longitude FROM vendors")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var vendors []models.Vendor
		for rows.Next() {
			var vendor models.Vendor
			err := rows.Scan(&vendor.ID, &vendor.Name, &vendor.BusinessName, &vendor.BusinessAddress, &vendor.PhoneNumber, &vendor.UniqueID, &vendor.Token, (*pq.StringArray)(&vendor.District), (*pq.StringArray)(&vendor.Panchayat), &vendor.Services, &vendor.ServiceCategory, &vendor.Experience, &vendor.Reference, &vendor.Premium, &vendor.ProfileImage, &vendor.PlanPurchase, &vendor.Image, &vendor.ImageView, &vendor.Favorites, &vendor.CreatedAt, &vendor.Latitude, &vendor.Longitude)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			vendors = append(vendors, vendor)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(vendors); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetVendorByDistrictAndPanchayat(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Getting vendor by district and panchayat\n")

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

		districtParam := r.URL.Query().Get("district")
		if districtParam == "" {
			http.Error(w, `{"error": "district parameter is required"}`, http.StatusBadRequest)
			return
		}

		districts := strings.Split(districtParam, ",")

		panchayatParam := r.URL.Query().Get("panchayat")
		if panchayatParam == "" {
			http.Error(w, `{"error": "panchayat parameter is required"}`, http.StatusBadRequest)
			return
		}

		panchayats := strings.Split(panchayatParam, ",")

		query := `SELECT id, name, business_name, business_address, phone_number, unique_id, token, district, panchayat, services, service_category, experience, reference, premium, profile_image, plan_purchase, image, image_view, favorites, created_at, latitude, longitude 
				  FROM vendors WHERE district && $1 AND panchayat && $2`

		rows, err := DB.Query(query, pq.StringArray(districts), pq.StringArray(panchayats))
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var vendors []models.Vendor
		for rows.Next() {
			var vendor models.Vendor
			err := rows.Scan(&vendor.ID, &vendor.Name, &vendor.BusinessName, &vendor.BusinessAddress, &vendor.PhoneNumber, &vendor.UniqueID, &vendor.Token, (*pq.StringArray)(&vendor.District), (*pq.StringArray)(&vendor.Panchayat), &vendor.Services, &vendor.ServiceCategory, &vendor.Experience, &vendor.Reference, &vendor.Premium, &vendor.ProfileImage, &vendor.PlanPurchase, &vendor.Image, &vendor.ImageView, &vendor.Favorites, &vendor.CreatedAt, &vendor.Latitude, &vendor.Longitude)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
				return
			}
			vendors = append(vendors, vendor)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(vendors); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func GetVendorByUniqueId(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Getting vendor by unique id\n")

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
		UniqueID, ok := vars["unique_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing unique_id in request")
			return
		}

		fmt.Println("Querying the database by vendor uniqueid:", UniqueID)
		var vendor models.Vendor
		err := DB.QueryRow("SELECT * FROM vendors WHERE unique_id = $1", UniqueID).Scan(&vendor.ID, &vendor.Name, &vendor.BusinessName, &vendor.BusinessAddress, &vendor.PhoneNumber, &vendor.UniqueID, &vendor.Token, (*pq.StringArray)(&vendor.District), (*pq.StringArray)(&vendor.Panchayat), &vendor.Services, &vendor.ServiceCategory, &vendor.Experience, &vendor.Reference, &vendor.Premium, &vendor.ProfileImage, &vendor.PlanPurchase, &vendor.Image, &vendor.ImageView, &vendor.Favorites, &vendor.CreatedAt, &vendor.Latitude, &vendor.Longitude)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Vendor not found", http.StatusNotFound)
				fmt.Println("Vendor not found")
			} else {
				http.Error(w, "Failed to query vendor", http.StatusInternalServerError)
				fmt.Printf("Failed to query vendor: %v\n", err)
			}
			return
		}

		fmt.Println("Vendor found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(vendor)
	}
}

func UpdateVendor(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Updating vendor by id\n")

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
		UniqueID, ok := vars["unique_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing unique_id in request")
			return
		}

		// Decode JSON request body into Vendor struct
		var vendor models.Vendor
		err := json.NewDecoder(r.Body).Decode(&vendor)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Failed to decode request body:", err)
			return
		}

		// Ensure the UniqueID from URL is used
		vendor.UniqueID = UniqueID

		// Update the vendor in the database
		query := `UPDATE vendors SET latitude = $1, longitude = $2 WHERE unique_id = $3`
		_, err = DB.Exec(query, vendor.Latitude, vendor.Longitude, vendor.UniqueID)

		if err != nil {
			http.Error(w, "Failed to update vendor", http.StatusInternalServerError)
			fmt.Printf("Failed to update vendor: %v\n", err)
			return
		}

		fmt.Println("Vendor updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(vendor)
	}
}

func DeleteVendor(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Deleting vendor by id\n")

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
		UniqueID, ok := vars["unique_id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			fmt.Println("Missing unique_id in request")
			return
		}

		// Delete the vendor from the database
		_, err := DB.Exec("DELETE FROM vendors WHERE unique_id=$1", UniqueID)
		if err != nil {
			http.Error(w, "Failed to delete vendor", http.StatusInternalServerError)
			fmt.Printf("Failed to delete vendor: %v\n", err)
			return
		}

		fmt.Println("Vendor deleted successfully")
		w.WriteHeader(http.StatusOK)
	}
}
