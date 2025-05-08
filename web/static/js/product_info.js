document.addEventListener('DOMContentLoaded', function() {
    console.clear(); // Clear any previous console logs
    console.log('=== PRODUCT INFO SCRIPT STARTED ===');
    
    // Basic elements
    const productData = document.getElementById('product-data');
    const commentsList = document.getElementById('comments-list');
    const commentText = document.getElementById('comment-text');
    const submitComment = document.getElementById('submit-comment');
    const likeButton = document.getElementById('like-button');
    
    // Check if these elements exist on the page
    console.log('Product data element exists:', !!productData);
    console.log('Comments list element exists:', !!commentsList);
    console.log('Comment text element exists:', !!commentText);
    console.log('Submit comment button exists:', !!submitComment);
    console.log('Like button exists:', !!likeButton);
    
    // Make sure we have the necessary elements before proceeding
    if (!productData || !commentsList) {
        console.error('Essential elements are missing from the page');
        return;
    }
    
    // Get the product ID from the data attribute
    const productId = productData.getAttribute('data-product-id');
    console.log('Product ID:', productId);
    
    // Handle audio player if present
    setupAudioPlayer();
    
    // Load comments on page load
    loadComments();
    
    // Initialize like functionality if the button exists
    if (likeButton) {
        initializeLikeButton();
    } else {
        // If the like button doesn't exist (user not logged in), we still load the like count
        updateLikeCount();
    }
    
    // Set up comment submission if form elements exist
    if (submitComment && commentText) {
        submitComment.addEventListener('click', function(e) {
            e.preventDefault();
            console.log('Submit comment button clicked');
            submitNewComment();
        });
    }
    
    /**
     * Initialize the like button functionality
     */
    function initializeLikeButton() {
        // First check if the user has already liked this product
        fetch(`/api/likes/${productId}`)
            .then(response => response.json())
            .then(data => {
                console.log('Like status:', data);
                
                if (data.success) {
                    updateLikeButtonState(data.is_liked, data.count);
                }
                
                // Add click handler to toggle like
                likeButton.addEventListener('click', toggleLike);
            })
            .catch(error => {
                console.error('Error checking like status:', error);
            });
    }
    
    /**
     * Toggle product like status
     */
    function toggleLike() {
        const isLiked = likeButton.classList.contains('liked');
        const method = isLiked ? 'DELETE' : 'POST';
        
        fetch(`/api/likes/${productId}`, {
            method: method
        })
        .then(response => response.json())
        .then(data => {
            console.log('Like toggle response:', data);
            
            if (data.success) {
                updateLikeButtonState(data.is_liked, data.count);
            }
        })
        .catch(error => {
            console.error('Error toggling like:', error);
        });
    }
    
    /**
     * Update the like button state based on user's like status
     */
    function updateLikeButtonState(isLiked, count) {
        if (isLiked) {
            likeButton.classList.add('liked');
        } else {
            likeButton.classList.remove('liked');
        }
        
        // Update the count
        const countElement = likeButton.querySelector('.like-count');
        if (countElement) {
            countElement.textContent = count;
        }
    }
    
    /**
     * Update like count for non-logged in users
     */
    function updateLikeCount() {
        fetch(`/api/likes/${productId}/count`)
            .then(response => response.json())
            .then(data => {
                console.log('Like count:', data);
                
                if (data.success) {
                    // Find the count element in any like button element (logged in or not)
                    const countElements = document.querySelectorAll('.like-count');
                    countElements.forEach(element => {
                        element.textContent = data.count;
                    });
                }
            })
            .catch(error => {
                console.error('Error fetching like count:', error);
            });
    }
    
    /**
     * Basic function to load comments via GET request
     */
    function loadComments() {
        console.log('Loading comments for product ID:', productId);
        
        // Show loading state
        commentsList.innerHTML = '<p>Loading comments...</p>';
        
        // Fetch comments from API
        fetch('/api/comments/' + productId)
            .then(function(response) {
                console.log('Comments API response status:', response.status);
                return response.json();
            })
            .then(function(data) {
                console.log('Comments data:', data);
                
                // Handle empty or error responses
                if (!data.success) {
                    commentsList.innerHTML = '<p>Error loading comments</p>';
                    return;
                }
                
                renderComments(data.comments || []);
            })
            .catch(function(error) {
                console.error('Error fetching comments:', error);
                commentsList.innerHTML = '<p>Failed to load comments</p>';
            });
    }
    
    /**
     * Simple function to render comments
     */
    function renderComments(comments) {
        console.log('Rendering', comments.length, 'comments');
        
        if (comments.length === 0) {
            commentsList.innerHTML = '<p>No comments yet. Be the first to comment!</p>';
            return;
        }
        
        // Build HTML for all comments
        let html = '';
        
        comments.forEach(function(comment) {
            // Format date
            const date = new Date(comment.created_at);
            const formattedDate = date.toLocaleString();
            
            // Default avatar
            const avatar = comment.profile_photo || '/static/images/default-avatar.png';
            
            // Check if current user is the author
            const userData = document.getElementById('user-data');
            const isAuthor = userData && userData.getAttribute('data-user-id') === String(comment.user_id);
            
            // Build individual comment HTML
            html += `
                <div class="comment-item">
                    <img src="${avatar}" alt="${comment.username}" class="comment-avatar">
                    <div class="comment-content">
                        <div class="comment-header">
                            <span class="comment-username">${comment.username}</span>
                            <span class="comment-date">${formattedDate}</span>
                            ${isAuthor ? '<button class="delete-comment-btn" data-id="' + comment.id + '">Delete</button>' : ''}
                        </div>
                        <p class="comment-text">${comment.comment}</p>
                    </div>
                </div>
            `;
        });
        
        // Update the comments list
        commentsList.innerHTML = html;
        
        // Add event listeners to delete buttons
        const deleteButtons = document.querySelectorAll('.delete-comment-btn');
        deleteButtons.forEach(function(button) {
            button.addEventListener('click', function() {
                const commentId = this.getAttribute('data-id');
                deleteComment(commentId);
            });
        });
        
        console.log('Comments rendering complete');
    }
    
    /**
     * Submit a new comment
     */
    function submitNewComment() {
        if (!commentText.value.trim()) {
            alert('Please enter a comment');
            return;
        }
        
        const comment = {
            product_id: parseInt(productId),
            comment: commentText.value.trim()
        };
        
        console.log('Submitting comment:', comment);
        
        fetch('/api/comments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(comment)
        })
        .then(function(response) {
            console.log('Comment submission status:', response.status);
            return response.json();
        })
        .then(function(data) {
            console.log('Comment submission response:', data);
            
            if (data.success) {
                // Clear the form
                commentText.value = '';
                
                // Reload comments
                loadComments();
            } else {
                alert('Failed to add comment: ' + (data.message || 'Unknown error'));
            }
        })
        .catch(function(error) {
            console.error('Error submitting comment:', error);
            alert('Error submitting comment. Please try again.');
        });
    }
    
    /**
     * Delete a comment
     */
    function deleteComment(commentId) {
        if (!confirm('Are you sure you want to delete this comment?')) {
            return;
        }
        
        console.log('Deleting comment ID:', commentId);
        
        fetch('/api/comments/' + commentId, {
            method: 'DELETE'
        })
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            console.log('Delete comment response:', data);
            
            if (data.success) {
                loadComments();
            } else {
                alert('Failed to delete comment: ' + (data.message || 'Unknown error'));
            }
        })
        .catch(function(error) {
            console.error('Error deleting comment:', error);
            alert('Error deleting comment. Please try again.');
        });
    }
    
    /**
     * Set up audio player if exists
     */
    function setupAudioPlayer() {
        const playButton = document.getElementById('play-button');
        if (!playButton) return;
        
        playButton.addEventListener('click', function() {
            const productId = this.getAttribute('data-product-id');
            const productName = this.getAttribute('data-product-name');
            const productOwner = this.getAttribute('data-product-owner');
            
            console.log('Play button clicked for:', productName);
            
            if (typeof playTrack === 'function') {
                const productImage = document.querySelector('.product-image').src;
                playTrack(parseInt(productId), productName, productOwner, productImage);
            } else {
                console.error('playTrack function not found');
            }
        });
    }
    
    console.log('=== PRODUCT INFO SCRIPT INITIALIZATION COMPLETE ===');
});