package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors" // Import the cors package
)

func main() {
	r := chi.NewRouter()
	teams, err := readTeamsFromFile("teams.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Apply ROT-23 transformation
	for i, team := range teams {
		teams[i] = rot23(team)
	}

	// Write the transformed teams back to teams.txt
	err = writeTeamsToFile("rotteams.txt", teams)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully transformed and written to teams.txt")
	// CORS middleware to allow requests from all origins
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},                                                       // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},                                  // Allowed methods
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Allowed headers
		ExposedHeaders:   []string{"Link"},                                                    // Exposed headers
		AllowCredentials: false,                                                               // Do not allow credentials
		MaxAge:           300,                                                                 // Cache duration
	}))

	r.Post("/rottingCorpse", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		flag := r.FormValue("input")
		// rot-23 the flag
		//read teams.txt but store each line in a string array
		rotteams, err := readTeamsFromFile("rotteams.txt")
		teams, err := readTeamsFromFile("teams.txt")
		if err != nil {
			http.Error(w, "Error reading teams file", http.StatusInternalServerError)
			return
		}
		cnt := 0
		for _, s := range rotteams {
			if s == flag {
				// Clear and  Write the new team's name to currentflag.txt
				os.Truncate("currentflag.txt", 0)
				err = os.WriteFile("currentflag.txt", []byte(teams[cnt]), 0644)
				if err != nil {
					http.Error(w, "Error writing to currentflag file", http.StatusInternalServerError)
					return
				}
				//fmt.Fprintf(w, "Flag found and written to currentflag.txt")
				successtring := "Braaaaiinnssss.... Successss... Teammmm.... " + teams[cnt]
				w.Write([]byte(successtring))
				return
			}
			cnt++
		}
		errstring := "Real Zombies are ROT-ten to the CO-23"
		w.Write([]byte(errstring))
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

func rot23(s string) string {
	rotated := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			rotated += string((c-'A'+23)%26 + 'A')
		} else if c >= 'a' && c <= 'z' {
			rotated += string((c-'a'+23)%26 + 'a')
		} else {
			rotated += string(c)
		}
	}
	return rotated
}

func writeTeamsToFile(filename string, teams []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, team := range teams {
		_, err := writer.WriteString(team + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func readTeamsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var teams []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		teams = append(teams, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
