package onesignal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationRequest struct {
	AppID            string            `json:"app_id"`
	IncludedSegments []string          `json:"included_segments"`
	Headings         map[string]string `json:"headings"`
	Contents         map[string]string `json:"contents"`
}

type NotificationInput struct {
	Headings map[string]string `json:"headings" binding:"required"`
	Contents map[string]string `json:"contents" binding:"required"`
}

// SendNotification function declaration
// Send a notification to a specific device

func SendNotification(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		appID := "4beb7160-79b8-49e9-9da2-b7520531055f"
		apiKey := "Yzk4MWI3NDItNjRhNy00YWY2LThjNzktNzY5NjY2NTE2NjM3"

		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Println("Invalid request method:", r.Method)
			return
		}

		// Decoding JSON Request Body
		// json.NewDecoder decode the JSON request body (r.Body) into a NotificationInput struct (input)

		decoder := json.NewDecoder(r.Body)
		var input NotificationInput
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Building the NotificationRequest struct

		notificationRequest := NotificationRequest{
			AppID:            appID,
			IncludedSegments: []string{"All"},
			Headings:         input.Headings,
			Contents:         input.Contents,
		}

		// Encoding the NotificationRequest struct into JSON

		requestBody, err := json.Marshal(notificationRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// creates an HTTP POST request (req) to the OneSignal API endpoint

		req, err := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(requestBody))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Sets request headers for JSON content type and API key authentication.

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", apiKey))

		// Sends the request using an HTTP client

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		fmt.Println("Notification send successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Notification send successfully")
	}
}
