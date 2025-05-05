document.addEventListener('DOMContentLoaded', function() {
    const productId = parseInt(document.getElementById('product-data').getAttribute('data-product-id'));
    const commentsList = document.getElementById('comments-list');
    const commentText = document.getElementById('comment-text');
    const likesProduct = document.getElementById('likes-product');
    const submitComment = document.getElementById('submit-comment');
    
    // Initialize comments
    loadComments();

    // Set up auto audio playback if it's an audio product
    const audioPlayer = document.getElementById('audio-player');
    if (audioPlayer) {
        const audio = document.getElementById('audio');
        if (audio) {
            audio.src = `/audio/${productId}`;
            audio.oncanplay = () => audio.play();
        }
    }
    
    // Comment submission handler
    if (submitComment) {
        submitComment.addEventListener('click', function() {
            if (!commentText.value.trim()) {
                alert('Будь ласка, введіть текст коментаря');
                return;
            }
            
            const comment = {
                product_id: productId,
                comment: commentText.valuez,
                likes_product: likesProduct.checked
            };
            
            fetch('/api/comments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(comment)
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    commentText.value = '';
                    likesProduct.checked = false;
                    loadComments();
                } 
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Сталася помилка при додаванні коментаря');
            });
        });
    }
    
    // Load comments from the API
    function loadComments() {
        commentsList.innerHTML = '<div class="loading-comments">Завантаження коментарів...</div>';
        
        fetch(`/api/comments/${productId}`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                displayComments(data.comments);
            } else {
                commentsList.innerHTML = '<div class="no-comments">Помилка завантаження коментарів</div>';
            }
        })
        .catch(error => {
            console.error('Error:', error);
            commentsList.innerHTML = '<div class="no-comments">Помилка завантаження коментарів</div>';
        });
    }
    
    // Display comments in the UI
    function displayComments(comments) {
        if (!comments || comments.length === 0) {
            commentsList.innerHTML = '<div class="no-comments">Немає коментарів. Будьте першим!</div>';
            return;
        }
        
        commentsList.innerHTML = '';
        
        comments.forEach(comment => {
            const commentDate = new Date(comment.created_at).toLocaleString('uk-UA');
            const defaultAvatar = '/static/images/default-avatar.png';
            const avatarSrc = comment.profile_photo ? comment.profile_photo : defaultAvatar;
            
            // Check if current user is the author of the comment
            const isAuthor = document.getElementById('user-data') && 
                document.getElementById('user-data').getAttribute('data-user-id') === comment.user_id.toString();
            const deleteButton = isAuthor ? 
                `<span class="comment-delete" data-id="${comment.id}">Видалити</span>` : '';

            const commentHtml = `
                <div class="comment-item" data-id="${comment.id}">
                    <img src="${avatarSrc}" alt="${comment.username}" class="comment-avatar">
                    <div class="comment-content">
                        <div class="comment-header">
                            <span class="comment-username">${comment.username}</span>
                            <span class="comment-date">${commentDate} ${deleteButton}</span>
                        </div>
                        <div class="comment-text">${comment.comment}</div>
                        ${comment.likes_product ? '<div class="comment-likes">❤️ Подобається продукт</div>' : ''}
                    </div>
                </div>
            `;
            
            commentsList.innerHTML += commentHtml;
        });

        // Add event listeners to delete buttons
        document.querySelectorAll('.comment-delete').forEach(button => {
            button.addEventListener('click', function() {
                const commentId = this.getAttribute('data-id');
                
                if (confirm('Ви впевнені, що хочете видалити цей коментар?')) {
                    deleteComment(commentId);
                }
            });
        });
    }
    
    // Delete a comment
    function deleteComment(commentId) {
        fetch(`/api/comments/${commentId}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                loadComments();
            } else {
                alert('Помилка: ' + data.message);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Сталася помилка при видаленні коментаря');
        });
    }

    // Make playAudio globally accessible
    window.playAudio = playAudio;
});