const express = require('express');
const cors = require('cors');  // Import the CORS middleware
const app = express();
const port = 3000;

// Use CORS to allow cross-origin requests from any origin
app.use(cors());

// Middleware to parse JSON bodies
app.use(express.json());

// Store the most recent flag
let recentFlag = null;

// API to handle incoming data and check for the correct flag pattern
app.post('/submit', (req, res) => {
    const userInput = req.body.input;

    // Regular expression to match the script with the alert flag
    const flagPattern = /<script>alert\((.*?)\)<\/script>/;

    // Check if the input matches the pattern
    const match = userInput.match(flagPattern);
    if (match) {
        // Store the captured flag (inside the alert function)
        recentFlag = match[1];
        console.log(`Flag captured: ${recentFlag}`);
    }

    // Respond with the most recent flag (whether updated or not)
    res.json({ recentFlag: recentFlag });
});

// Add a simple route to test the server and CORS headers
app.get('/test', (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');  // Set CORS header manually
    res.send('recentflag');
});

// Start the server
app.listen(port, () => {
    console.log(`Server running on http://localhost:${port}`);
});
