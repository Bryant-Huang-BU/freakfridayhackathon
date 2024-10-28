package main

import (
    "database/sql"
    "fmt"
    "net/http"
    "strings"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    _ "github.com/go-sql-driver/mysql"
)

func SplitSQLStatements(query string) []string {
    statements := strings.Split(query, ";")
    var cleanedStatements []string
    for _, stmt := range statements {
        trimmedStmt := strings.TrimSpace(stmt)
        if trimmedStmt != "" {
            cleanedStatements = append(cleanedStatements, trimmedStmt)
        }
    }
    return cleanedStatements
}

func main() {
    r := chi.NewRouter()

    // CORS middleware
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: false,
        MaxAge:           300,
    }))

    db, err := sql.Open("mysql", "root:cybears@tcp(localhost:3306)/SPOOKY")
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return
    }
    defer db.Close()

    r.Post("/query", func(w http.ResponseWriter, r *http.Request) {
        flag := r.URL.Query().Get("flag")
        if flag == "" {
            http.Error(w, "Flag parameter is missing", http.StatusBadRequest)
            return
        }

        statements := SplitSQLStatements(flag)
        var results []string

        for _, stmt := range statements {
            fmt.Printf("Executing statement: %s\n", stmt)

            // Check if the statement is a SELECT
            if strings.HasPrefix(stmt, "SELECT") || strings.HasPrefix(stmt, "show") {
                rows, err := db.Query(stmt)
                if err != nil {
                    fmt.Printf("Error executing query: %v\n", err)
                    results = append(results, fmt.Sprintf("Error: %v", err))
                    continue
                }
                defer rows.Close()

                var result string
                for rows.Next() {
                    var id int
                    var name string
                    if err := rows.Scan(&id, &name); err != nil {
                        fmt.Printf("Error scanning result: %v\n", err)
                        results = append(results, fmt.Sprintf("Error scanning: %v", err))
                        continue
                    }
                    result += fmt.Sprintf("ID: %d, Name: %s\n", id, name)
                }
                results = append(results, result)
            } else {
                // Handle non-SELECT statements
                result, err := db.Exec(stmt)
                if err != nil {
                    fmt.Printf("Error executing statement: %v\n", err)
                    results = append(results, fmt.Sprintf("Error: %v", err))
                    continue
                }

                // Get affected rows
                rowsAffected, err := result.RowsAffected()
                if err != nil {
                    fmt.Printf("Error getting rows affected: %v\n", err)
                    results = append(results, fmt.Sprintf("Error: %v", err))
                    continue
                }
                results = append(results, fmt.Sprintf("Rows affected: %d", rowsAffected))
            }
        }

        // Combine results and send them as a response
        finalResult := strings.Join(results, "\n")
        w.Write([]byte(finalResult))
    })

    fmt.Println("Starting server on :4000")
    http.ListenAndServe(":4000", r)
}
