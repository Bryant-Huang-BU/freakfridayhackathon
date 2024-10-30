// Function to fetch the recent flag on page load
async function fetchRecentFlag() {
    try {
        const response = await fetch('http://localhost:3000/recent-flag');
        const data = await response.json();
        
        // Display the recent flag on the website
        const flagDisplayElement = document.getElementById('flagDisplay');
        flagDisplayElement.textContent = `Beware as you step in... TEAM ${data.recentFlag || 'spirits'} is restless tonight. `;
    } catch (error) {
        console.error('Error fetching recent flag:', error);
    }
}

// Call the function when the page loads
document.addEventListener('DOMContentLoaded', fetchRecentFlag);

document.getElementById('commentForm').addEventListener('submit', function(event) {
    event.preventDefault();

    // Get user input from the comment form
    const nameInput = document.getElementById('name').value;
    const commentInput = document.getElementById('comment').value;

    // Create a new list item for the comment
    const li = document.createElement('li');
    li.textContent = `${nameInput}: ${commentInput}`;
    document.getElementById('commentList').appendChild(li);

    // Check if the comment includes <script> tags
    const scriptMatch = commentInput.match(/<script>(.*?)<\/script>/);

    if (scriptMatch) {
        const scriptContent = scriptMatch[1];  // Extract the content within the <script> tags

        try {
            console.log('Sending data to server:', { input: commentInput }); // Log what you're sending
            fetch('http://localhost:3000/submit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ input: commentInput })
            })
            .then(response => {
                console.log('Fetch response:', response); // Log the response object
                return response.json();
            })
            .then(data => {
                console.log('Server response:', data);

                // Display the recent flag on the website
                const flagDisplayElement = document.getElementById('flagDisplay');
                flagDisplayElement.textContent = `Beware as you step in... TEAM ${data.recentFlag || 'spirits'} is restless tonight. `;
            })
            .catch(error => {
                console.error('Fetch error:', error); // Capture errors related to fetch
            });

            // Optional: Log warning and manually execute the script
            console.warn("Executing user script:", scriptContent);
            const script = document.createElement('script');
            script.textContent = scriptContent;
            document.body.appendChild(script);  // Append and execute script

            // Additional security: remove script tags from the displayed comment
            li.textContent = `${nameInput}: ${commentInput.replace(/<script>.*?<\/script>/, '[Impossible! How did you break past my defenses?? Tell no one!]')}`;
        } catch (error) {
            console.error('Error sending data to server:', error);
        }
    }
});
Executed