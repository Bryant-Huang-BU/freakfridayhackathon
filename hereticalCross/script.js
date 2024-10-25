document.getElementById('commentForm').addEventListener('submit', function(e) {
    e.preventDefault(); // Prevent the form from submitting traditionally

    // Get input values
    const name = document.getElementById('name').value;
    const comment = document.getElementById('comment').value;

    // Create a new list item for the comment
    const newComment = document.createElement('li');
    newComment.innerHTML = `<strong>${name}:</strong> ${comment}`;

    // Add the new comment to the list
    document.getElementById('commentList').appendChild(newComment);

    // Clear form fields
    document.getElementById('name').value = '';
    document.getElementById('comment').value = '';
});