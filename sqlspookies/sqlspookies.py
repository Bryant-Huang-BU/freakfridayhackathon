from flask import Flask, request, jsonify
from flask_cors import CORS
import pymysql

app = Flask(__name__)
CORS(app)

# Database connection details
db_config = {
    'host': 'localhost',
    'user': 'root',
    'password': 'yourpassword',
    'database': 'users',
    'client_flag': pymysql.constants.CLIENT.MULTI_STATEMENTS
}

def run_query(subject):
    connection = None
    try:
        # Connect to the database
        connection = pymysql.connect(**db_config, autocommit=True)
        print(f"Running query with input: {subject}")
        
        with connection.cursor() as cursor:
            # Modify the SQL injection attempt with space after `--`
            # Example input: "'; SELECT * FROM users; -- "
            #subject = "'; SELECT * FROM users; -- "
            print("Executing query with injection attempt:", subject)  # Log the injection input
            
            query = f"SELECT * FROM users WHERE name = '{subject}'"
            print("Final query executed:", query)
            
            cursor.execute(query)
            results = []
            
            while True:
                fetch_result = cursor.fetchall()
                if fetch_result:
                    print("Fetch result:", fetch_result)
                    results.extend(fetch_result)
                
                # Check for additional result sets
                if not cursor.nextset():
                    break
            
            # Log final results
            print("Query Results:", results if results else "No results found.")
            return results
    
    except Exception as e:
        print("Error executing query:", e)  # Detailed error logging
        return {"error": str(e)}
    
    finally:
        if connection:
            connection.close()
            print("Database connection closed.")


# Example input for testing
test_input = "'; SELECT * FROM users; --"
results = run_query(test_input)
print("Results:", results)



@app.route('/query', methods=['POST'])
def execute_query():
    # Get the query from the URL-encoded parameter
    user_query = request.form.get('flag')
    
    # Check if user_query is None or empty
    if not user_query:
        return jsonify({"error": "No query provided"}), 400
    
    # Run the query and fetch results
    results = run_query(user_query)
    
    # Check if results is a list, otherwise wrap in dict for error message
    if not isinstance(results, list):
        return jsonify({"error": results})  # Return the error message
    
    return jsonify(results)

if __name__ == '__main__':
    app.run(debug=True)
