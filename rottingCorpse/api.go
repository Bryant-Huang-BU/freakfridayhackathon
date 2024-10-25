package main

import (
    "os"
    "strings"
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors" // Import the cors package
)

func main() {
    r := chi.NewRouter()

    // CORS middleware to allow requests from all origins
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"}, // Allow all origins
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allowed methods
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Allowed headers
        ExposedHeaders:   []string{"Link"}, // Exposed headers
        AllowCredentials: false, // Do not allow credentials
        MaxAge:           300, // Cache duration
    }))

    r.Post("/rottingCorpse", func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Error parsing form", http.StatusBadRequest)
            return
        }
        flag := r.FormValue("input")
        // rot-10 the flag
        rotFlag := ""
        for _, c := range flag {
            if c >= 'A' && c <= 'Z' {
                rotFlag += string((c-'A'+10)%26 + 'A')
            } else if c >= 'a' && c <= 'z' {
                rotFlag += string((c-'a'+10)%26 + 'a')
            } else {
                rotFlag += string(c)
            }
        }
        teams, err := os.ReadFile("teams.txt")
        if err != nil {
            http.Error(w, "Error reading teams file", http.StatusInternalServerError)
            return
        }
        // Check if the rotated flag is in the teams list
        if strings.Contains(string(teams), rotFlag) {
            // Write the new team's name to currentflag.txt
            err = os.WriteFile("currentflag.txt", []byte(flag), 0644)
            if err != nil {
                http.Error(w, "Error writing to currentflag file", http.StatusInternalServerError)
                return
            }
            fmt.Fprintf(w, "Flag found and written to currentflag.txt")
        } else {
            http.Error(w, "Flag not found in teams list", http.StatusNotFound)
        }
    })

    r.Get("/winningteam", func(w http.ResponseWriter, r *http.Request) {
        content, err := os.ReadFile("currentflag.txt")
        if err != nil {
            http.Error(w, "Error reading currentflag file", http.StatusInternalServerError)
            return
        }
        w.Write(content)
    })

    log.Println("Starting server on :4000")
    log.Fatal(http.ListenAndServe(":4000", r))
}
