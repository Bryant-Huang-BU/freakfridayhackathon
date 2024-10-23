package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

func main() 
	// Create a new HTTP server
	http.HandleFunc"/", func(w http.ResponseWriter, r *http.Request) 
		// Check if the request is a GET request
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}

		// Get the query parameter from the URL
		flagID := r.URL.Query().Get("flag_id")

		// Connect to the database
		db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Query the flags table with the user-provided flag_id
        
		// This is vulnerable to SQL injection attacks
		query := fmt.Sprintf("SELECT * FROM employees WHERE id = '%s'", flagID)
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Fetch the results
		var results []map[string]string
		for rows.Next() {
			var id int
			var flag string
			err := rows.Scan(&id, &flag)
			if err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			results = append(results, map[string]string{"id": fmt.Sprintf("%d", id), "flag": flag})
		}

		// Return the results as JSON
		json.NewEncoder(w).Encode(results)
	

	log.Fatal(http.ListenAndServe(":8080", nil))
