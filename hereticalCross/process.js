const express = require('express');
const cors = require('cors');
const app = express();
const port = 3000;
const fs = require('fs'); // Import the fs module


// Use CORS to allow cross-origin requests from any origin
app.use(cors());

// Middleware to parse JSON bodies
app.use(express.json());

// Store the most recent flag
let recentFlag = null;
const writeFlagToFile = (flag) => {
    const jsonData = { recentFlag: flag };
    fs.writeFileSync('./lol.json', JSON.stringify(jsonData, null, 2), 'utf-8');
};
const readFlagFromFile = () => {
    const data = fs.readFileSync('./lol.json', 'utf-8'); // Read the file
    return JSON.parse(data); // Parse and return the JSON object
};
// API to handle incoming data and check for the correct flag pattern
app.post('/submit', (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');
    const userInput = req.body.input;

    // Regular expression to match the script with the alert flag
    const flagPattern = /<script>alert\((.*?)\)<\/script>/;

    // Check if the input matches the pattern
    const match = userInput.match(flagPattern);
    if (match) {
        // Store the captured flag (inside the alert function)
        recentFlag = match[1];
        console.log(`Flag captured: ${recentFlag}`);
        writeFlagToFile(recentFlag)
    }

    // Respond with the most recent flag (whether updated or not)
    res.json({ recentFlag: recentFlag });
});

// New API to get the most recent flag
app.get('/recent-flag', (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');
    console.log(readFlagFromFile());
    res.json(readFlagFromFile());
});

// Add a simple route to test the server and CORS headers
app.get('/test', (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.send('test');
});

// Start the server
app.listen(port, () => {
    console.log(`Server running on http://localhost:${port}`);
});
