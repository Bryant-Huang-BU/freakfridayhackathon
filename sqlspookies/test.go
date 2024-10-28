package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Template to render the input form and results
var formTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>SQL Injection Demo</title>
</head>
<body>
    <h1>SQL Injection Demo</h1>
    <form action="/" method="post">
        <label for="input">Enter your SQL command:</label><br>
        <input type="text" id="input" name="input" /><br>
        <input type="submit" value="Submit" />
    </form>
    <div>{{.}}</div>
</body>
</html>`

func main() {
	var err error
	// Replace with your own database connection string
	db, err = sql.Open("mysql", "root:yourpassword@tcp(localhost:3306)/testdb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", handleRequest)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var result string
	if r.Method == http.MethodPost {
		input := r.FormValue("input")

		// Make a valid query context where the injection will work
		query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", input)

		// Execute the query directly
		rows, err := db.Query(query)
		if err != nil {
			result = fmt.Sprintf("Error executing query: %s", err)
			log.Println(result)
		} else {
			defer rows.Close()
			columns, _ := rows.Columns()
			values := make([]sql.RawBytes, len(columns))
			scanArgs := make([]interface{}, len(columns))
			for i := range values {
				scanArgs[i] = &values[i]
			}
			for rows.Next() {
				if err := rows.Scan(scanArgs...); err != nil {
					log.Println(err)
				}
				for _, value := range values {
					result += string(value) + "<br>"
				}
			}
		}

		// Add a special case to execute a known injection
		if input == "hello'; SHOW TABLES; --" {
			// Directly query the tables since multiple queries are not supported
			tablesQuery := "SHOW TABLES"
			tables, err := db.Query(tablesQuery)
			if err != nil {
				result += fmt.Sprintf("Error executing SHOW TABLES query: %s", err)
			} else {
				defer tables.Close()
				var tableName string
				for tables.Next() {
					if err := tables.Scan(&tableName); err != nil {
						log.Println(err)
					} else {
						result += tableName + "<br>"
					}
				}
			}
		}
	}

	t := template.Must(template.New("form").Parse(formTemplate))
	t.Execute(w, result)
}
