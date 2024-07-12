package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"tendercall.com/main/adrequest_db"
	"tendercall.com/main/applogs_db"
	"tendercall.com/main/autotaxigoods_db"
	"tendercall.com/main/banner_db"
	"tendercall.com/main/chat_db"
	"tendercall.com/main/customer_db"
	"tendercall.com/main/logs_db"
	"tendercall.com/main/onesignal"
	"tendercall.com/main/order_db"
	"tendercall.com/main/ratings_db"
	"tendercall.com/main/registeredusers_db"
	"tendercall.com/main/upload"
	"tendercall.com/main/vendor_db"
)

var DB *sql.DB

func main() {

	// PostgreSQL connection parameters
	const (
		host     = "ep-empty-resonance-a1usmq29.ap-southeast-1.pg.koyeb.app"
		port     = 5432
		user     = "koyeb-adm"
		password = "M3OKWT2pcHYa"
		dbname   = "koyebdb"
	)

	// Construct the connection string
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	// Attempt to connect to the database
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Database connection established")

	router := mux.NewRouter()

	//Upload router
	router.HandleFunc("/upload", upload.HandleFileUpload).Methods("POST")
	router.HandleFunc("/upload", upload.HandleGetURL).Methods("GET")

	//Notification router
	router.HandleFunc("/notification", onesignal.SendNotification(DB)).Methods("POST")

	//Chat router
	router.HandleFunc("/chat", chat_db.CreateChat(DB)).Methods("POST")
	router.HandleFunc("/chat", chat_db.GetChat(DB)).Methods("GET")
	router.HandleFunc("/chat/{id}", chat_db.GetChatById(DB)).Methods("GET")

	//Vendor router
	router.HandleFunc("/vendors", vendor_db.CreateVendor(DB)).Methods("POST")
	router.HandleFunc("/vendors", vendor_db.GetVendor(DB)).Methods("GET")
	router.HandleFunc("/vendors", vendor_db.GetVendorByDistrictAndPanchayat(DB)).Methods("GET")
	router.HandleFunc("/vendors/{unique_id}", vendor_db.GetVendorByUniqueId(DB)).Methods("GET")
	router.HandleFunc("/vendors/{unique_id}", vendor_db.UpdateVendor(DB)).Methods("PUT")
	router.HandleFunc("/vendors/{unique_id}", vendor_db.DeleteVendor(DB)).Methods("DELETE")

	//Customer router
	router.HandleFunc("/customers", customer_db.CreateCustomer(DB)).Methods("POST")
	router.HandleFunc("/customers", customer_db.GetCustomer(DB)).Methods("GET")
	router.HandleFunc("/customers/{userid}", customer_db.GetCustomerByUserId(DB)).Methods("GET")
	router.HandleFunc("/customers/{userid}", customer_db.UpdateCustomer(DB)).Methods("PUT")
	router.HandleFunc("/customers/{userid}", customer_db.DeleteCustomer(DB)).Methods("DELETE")

	//Banner router
	router.HandleFunc("/banners", banner_db.CreateBanner(DB)).Methods("POST")
	router.HandleFunc("/banners", banner_db.GetBanner(DB)).Methods("GET")
	router.HandleFunc("/banners/{id}", banner_db.GetBannerById(DB)).Methods("GET")
	router.HandleFunc("/banners/{id}", banner_db.UpdateBanner(DB)).Methods("PUT")
	router.HandleFunc("/banners/{id}", banner_db.DeleteBanner(DB)).Methods("DELETE")

	//Ratings router
	router.HandleFunc("/ratings", ratings_db.CreateRating(DB)).Methods("POST")
	router.HandleFunc("/ratings", ratings_db.GetRating(DB)).Methods("GET")
	router.HandleFunc("/ratings/{id}", ratings_db.GetRatingById(DB)).Methods("GET")
	router.HandleFunc("/ratings/{id}", ratings_db.UpdateRating(DB)).Methods("PUT")
	router.HandleFunc("/ratings/{id}", ratings_db.DeleteRating(DB)).Methods("DELETE")

	//AdRequest router
	router.HandleFunc("/ad_request", adrequest_db.CreateAdRequest(DB)).Methods("POST")
	router.HandleFunc("/ad_request", adrequest_db.GetAdRequest(DB)).Methods("GET")
	router.HandleFunc("/ad_request/{unique_id}", adrequest_db.GetAdRequestById(DB)).Methods("GET")
	router.HandleFunc("/ad_request/{unique_id}", adrequest_db.UpdateAdRequest(DB)).Methods("PUT")
	router.HandleFunc("/ad_request/{unique_id}", adrequest_db.DeleteAdRequest(DB)).Methods("DELETE")

	//Order router
	router.HandleFunc("/orders", order_db.CreateOrder(DB)).Methods("POST")
	router.HandleFunc("/orders", order_db.GetOrder(DB)).Methods("GET")
	router.HandleFunc("/orders/{user_id}", order_db.GetOrderByUserId(DB)).Methods("GET")
	router.HandleFunc("/orders/{user_id}", order_db.UpdateOrder(DB)).Methods("PUT")
	router.HandleFunc("/orders/{user_id}", order_db.DeleteOrder(DB)).Methods("DELETE")

	//AutoTaxiGoods router
	router.HandleFunc("/autotaxigoods", autotaxigoods_db.CreateAutoTaxiGoods(DB)).Methods("POST")
	router.HandleFunc("/autotaxigoods", autotaxigoods_db.GetAutoTaxiGoods(DB)).Methods("GET")
	router.HandleFunc("/autotaxigoods/{userid}", autotaxigoods_db.GetAutoTaxiGoodsByUserId(DB)).Methods("GET")
	router.HandleFunc("/autotaxigoods/{userid}", autotaxigoods_db.UpdateAutoTaxiGoods(DB)).Methods("PUT")
	router.HandleFunc("/autotaxigoods/{userid}", autotaxigoods_db.DeleteAutoTaxiGoods(DB)).Methods("DELETE")

	//RegisteredUser router
	router.HandleFunc("/users", registeredusers_db.CreateRegisteredUser(DB)).Methods("POST")
	router.HandleFunc("/users", registeredusers_db.GetRegisteredUser(DB)).Methods("GET")
	router.HandleFunc("/users/{id}", registeredusers_db.GetRegisteredUserById(DB)).Methods("GET")
	router.HandleFunc("/users/{id}", registeredusers_db.UpdateRegisteredUser(DB)).Methods("PUT")
	router.HandleFunc("/users/{id}", registeredusers_db.DeleteRegisteredUser(DB)).Methods("DELETE")

	//Logs router
	router.HandleFunc("/logs", logs_db.CreateLogs(DB)).Methods("POST")
	router.HandleFunc("/logs", logs_db.GetLogs(DB)).Methods("GET")
	router.HandleFunc("/logs/{id}", logs_db.GetLogsById(DB)).Methods("GET")
	router.HandleFunc("/logs/{id}", logs_db.UpdateLogs(DB)).Methods("PUT")
	router.HandleFunc("/logs/{id}", logs_db.DeleteLogs(DB)).Methods("DELETE")

	//Applogs router
	router.HandleFunc("/applogs", applogs_db.CreateAppLogs(DB)).Methods("POST")
	router.HandleFunc("/applogs", applogs_db.GetAppLogs(DB)).Methods("GET")
	router.HandleFunc("/applogs/{userid}", applogs_db.GetAppLogsById(DB)).Methods("GET")
	router.HandleFunc("/applogs/{userid}", applogs_db.UpdateAppLogs(DB)).Methods("PUT")
	router.HandleFunc("/applogs/{userid}", applogs_db.DeleteAppLogs(DB)).Methods("DELETE")

	log.Println("Server is starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Define the upload endpoint
// router.POST("/upload", func(c *gin.Context) {
// 	// Get the file from the request
// 	file, header, err := c.Request.FormFile("image")
// 	if err != nil {
// 		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
// 		return
// 	}
// 	defer file.Close()

// 	// 	// Create a unique file name
// 	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))

// 	// 	// Upload file to Firebase Storage
// 	writer := bucket.Object(fileName).NewWriter(ctx)
// 	if _, err := io.Copy(writer, file); err != nil {
// 		c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
// 		return
// 	}

// 	if err := writer.Close(); err != nil {
// 		c.String(http.StatusInternalServerError, fmt.Sprintf("close writer err: %s", err.Error()))
// 		return
// 	}

// 	// 	// Construct the URL of the uploaded image
// 	imageURL := "https://firebasestorage.googleapis.com/v0/b/" + bucketName + "//" + fileName

// 	// 	// File uploaded successfully, return the URL
// 	c.JSON(http.StatusAccepted, gin.H{"image": imageURL})

// })
