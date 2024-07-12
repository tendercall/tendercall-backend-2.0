package chat_db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"tendercall.com/main/models"
)

func CreateChat(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateChat handler called")

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

		fmt.Println("Querying the database for Chat")

		var nextID int
		err := DB.QueryRow("SELECT COALESCE(MAX(id) + 1, 1) FROM chat_assistant").Scan(&nextID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var chatassistant models.ChatAssistant
		err = decoder.Decode(&chatassistant)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		if chatassistant.Message == "" {
			http.Error(w, "Message field is required", http.StatusBadRequest)
			return
		}

		currentTime := time.Now()

		// Insert the chat data into the database
		_, err = DB.Exec("INSERT INTO chat_assistant (id, message, userid, createdat) VALUES ($1, $2, $3, $4)", nextID, chatassistant.Message, chatassistant.UserID, currentTime)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error inserting new chat:", err)
			return
		}

		fmt.Println("Chat inserted successfully")

		responseMessage := ""
		if chatassistant.Message == "Hello" || chatassistant.Message == "Hi" {
			responseMessage = "How can I help you?"
		} else {
			responseMessage = "Sorry I don't understand"
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"Response": responseMessage})
	}
}

func GetChat(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetChat handler called")

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

		fmt.Println("Querying the database for Chat")
		rows, err := DB.Query("SELECT id,message,userid,createdat FROM chat_assistant")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		var chatassistants []models.ChatAssistant
		for rows.Next() {
			var chatassistant models.ChatAssistant
			err := rows.Scan(&chatassistant.ID, &chatassistant.Message, &chatassistant.UserID, &chatassistant.CreatedAt)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			chatassistants = append(chatassistants, chatassistant)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error with rows:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(chatassistants); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Response sent successfully")
	}
}

func GetChatById(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting Chat by userid\n")

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
			fmt.Println("Missing id in request")
			return
		}

		fmt.Println("Querying the database by Chat id:", ID)
		var chatAssistant models.ChatAssistant
		err := DB.QueryRow("SELECT * FROM chat_assistant WHERE id = $1", ID).Scan(&chatAssistant.ID, &chatAssistant.Message, &chatAssistant.UserID, &chatAssistant.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Chat not found", http.StatusNotFound)
				fmt.Println("Chat not found")
			} else {
				http.Error(w, "Failed to query Chat", http.StatusInternalServerError)
				fmt.Printf("Failed to query Chat: %v\n", err)
			}
			return
		}

		fmt.Println("Chat found successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chatAssistant)
	}
}
