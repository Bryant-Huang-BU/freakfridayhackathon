document.getElementById('commentForm').addEventListener('submit', function(e) {
    e.preventDefault();

    const name = document.getElementById('name').value;
    const userInput = document.getElementById('comment').value;

    const newComment = document.createElement('li');
    newComment.innerHTML = `<strong>${name}:</strong> ${userInput}`;

    const scriptMatch = userInput.match(/<script>(.*?)<\/script>/);
    if (scriptMatch) {
        const script = document.createElement('script');
        script.textContent = scriptMatch[1];
        document.body.appendChild(script);
    }

    document.getElementById('commentList').appendChild(newComment);

    document.getElementById('name').value = '';
    document.getElementById('comment').value = '';
});
