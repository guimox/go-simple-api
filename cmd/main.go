package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	mu  sync.Mutex
	err error

	users   = make(map[int]User)
	lockers = make(map[int]Locker)

	nextUID int
	nextLID int
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "lockers_db"
)

type User struct {
	ID        int      `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Password  string   `json:"-"`
	Lockers   []Locker `json:"lockers"`
}

type Locker struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	Status string `json:"status"`
	UserID int    `json:"user_id"`
}

func main() {
	// Build the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to ensure the connection is established
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to the database", err)
	}

	log.Println("Connected to the database!")

	createTables()

	// Router setup
	r := mux.NewRouter()
	r.HandleFunc("/users", getAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	r.HandleFunc("/lockers", getAllLockers).Methods("GET")
	r.HandleFunc("/lockers/{id}", getLocker).Methods("GET")
	r.HandleFunc("/users/{user_id}/lockers", createLocker).Methods("POST")
	r.HandleFunc("/lockers/{id}", updateLocker).Methods("PUT")
	r.HandleFunc("/lockers/{id}", deleteLocker).Methods("DELETE")

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}

func createTables() {
	// Create users table
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	);`

	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	log.Println("Users table created or already exists.")

	// Create lockers table
	lockerTable := `
	CREATE TABLE IF NOT EXISTS lockers (
		id SERIAL PRIMARY KEY,
		number VARCHAR(50) NOT NULL,
		status VARCHAR(50) NOT NULL,
		user_id INT,
		CONSTRAINT fk_user
			FOREIGN KEY(user_id)
			REFERENCES users(id)
			ON DELETE SET NULL
	);`

	_, err = db.Exec(lockerTable)
	if err != nil {
		log.Fatalf("Error creating lockers table: %v", err)
	}
	log.Println("Lockers table created or already exists.")
}

// User handlers

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var usersList []User
	for _, user := range users {
		usersList = append(usersList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersList)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if user, ok := users[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}

	http.NotFound(w, r)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = nextUID
	nextUID++

	users[user.ID] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user, ok := users[id]; ok {
		updatedUser.ID = user.ID
		users[id] = updatedUser

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedUser)
		return
	}

	http.NotFound(w, r)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if _, ok := users[id]; ok {
		delete(users, id)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.NotFound(w, r)
}

// Locker handlers

func getAllLockers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var lockersList []Locker
	for _, locker := range lockers {
		lockersList = append(lockersList, locker)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lockersList)
}

func getLocker(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if locker, ok := lockers[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(locker)
		return
	}

	http.NotFound(w, r)
}

func createLocker(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["user_id"])

	// Check if user exists
	user, ok := users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode the new locker
	var locker Locker
	if err := json.NewDecoder(r.Body).Decode(&locker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assign a new ID and associate with the user
	locker.ID = nextLID
	locker.UserID = userID
	nextLID++

	// Add the locker to the lockers map
	lockers[locker.ID] = locker

	// Update the user's lockers
	user.Lockers = append(user.Lockers, locker)
	users[userID] = user

	// Respond with the created locker
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(locker)
}

func updateLocker(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedLocker Locker
	if err := json.NewDecoder(r.Body).Decode(&updatedLocker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if locker, ok := lockers[id]; ok {
		updatedLocker.ID = locker.ID
		updatedLocker.UserID = locker.UserID
		lockers[id] = updatedLocker

		// Update locker in user's list
		user := users[locker.UserID]
		for i, l := range user.Lockers {
			if l.ID == locker.ID {
				user.Lockers[i] = updatedLocker
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedLocker)
		return
	}

	http.NotFound(w, r)
}

func deleteLocker(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if locker, ok := lockers[id]; ok {
		// Remove locker from user's list
		user := users[locker.UserID]
		for i, l := range user.Lockers {
			if l.ID == locker.ID {
				user.Lockers = append(user.Lockers[:i], user.Lockers[i+1:]...)
				break
			}
		}

		delete(lockers, id)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.NotFound(w, r)
}
