// Comments related functions
function renderComments(comments) {
    if (!Array.isArray(comments) || comments.length === 0) {
        return '<p class="no-comments">No comments yet</p></br>';
    }
    
    return comments.map(comment => `
        <div class="comment">
            <div class="comment-content">
                <p>${comment.content}</p>
            </div>
            <div class="comment-footer">
                <small>By ${comment.username} on ${new Date(comment.created_at).toLocaleString()}</small>
                <div class="comment-actions">
                    <button onclick="handleCommentLike(${comment.id}, true)" class="action-btn">
                        👍 <span>${comment.likes || 0}</span>
                    </button>
                    <button onclick="handleCommentLike(${comment.id}, false)" class="action-btn">
                        👎 <span>${comment.dislikes || 0}</span>
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}
function toggleComments(postId) {
    const container = document.getElementById(`comments-container-${postId}`);
    if (container.style.display === 'none') {
        container.style.display = 'block';
        openCommentSections.add(postId);
    } else {
        container.style.display = 'none';
        openCommentSections.delete(postId);
    }
}
function showCommentForm(postId) {
    const commentsContainer = document.getElementById(`comments-container-${postId}`);
    const form = document.getElementById(`comment-form-${postId}`);
    commentsContainer.style.display = 'block';
    form.style.display = 'block';
}
async function submitComment(postId) {
    const commentForm = document.getElementById(`comment-form-${postId}`);
    const textarea = commentForm.querySelector('textarea');
    const content = textarea.value.trim();

    if (!content) {
        handleError('Comment cannot be empty');
        return;
    }

    try {
        // First check authentication status
        const authResponse = await fetch('/api/protected/api/auth/status');
        const authData = await authResponse.json();
        
        if (!authData.authenticated) {
            handleError('Please login to comment');
            return;
        }

        const response = await fetch('/api/protected/api/comments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                post_id: postId,
                content: sanitizeHTML(content)
            })
        });

        if (!response.ok) {
            throw new Error('Failed to submit comment');
        }

        // Wait for the response before clearing the form
        await response.json();
        
        textarea.value = '';
        commentForm.style.display = 'none';
        
        // Store the current open comments
        const openComments = new Set();
        document.querySelectorAll('.comments-container').forEach(container => {
            if (container.style.display !== 'none') {
                openComments.add(parseInt(container.id.split('-').pop()));
            }
        });
        
        // Fetch the updated post data
        await fetchPosts();
        
        // Restore open state of comments
        openComments.forEach(id => {
            const container = document.getElementById(`comments-container-${id}`);
            if (container) {
                container.style.display = 'block';
            }
        });
        
        handleSuccess('Comment posted successfully');
    } catch (error) {
        console.error('Error submitting comment:', error);
        handleError('Failed to post comment');
    }
}
async function handleCommentLike(commentId, isLike) {
    try {
        // Check authentication first
        const authResponse = await fetch('/api/protected/api/auth/status');
        const authData = await authResponse.json();
        
        if (!authData.authenticated) {
            handleError('Please login to like comments');
            return;
        }

        // Store the current open comments
        const openComments = new Set();
        document.querySelectorAll('.comments-container').forEach(container => {
            if (container.style.display !== 'none') {
                openComments.add(parseInt(container.id.split('-').pop()));
            }
        });

        const response = await fetch('/api/protected/api/comments/like', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                comment_id: commentId,
                is_like: isLike
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to update like');
        }
        
        await fetchPosts();
        
        // Restore open state of comments
        openComments.forEach(id => {
            const container = document.getElementById(`comments-container-${id}`);
            if (container) {
                container.style.display = 'block';
            }
        });
        
        handleSuccess(isLike ? 'Comment liked!' : 'Comment disliked!');
    } catch (error) {
        console.error('Error handling comment like:', error);
        handleError('Failed to update like');
    }
} 
