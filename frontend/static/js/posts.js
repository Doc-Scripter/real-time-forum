let currentPage = 1;
let isLoading = false;
let hasMorePosts = true;
let currentCategory = '';
let currentFilter = '';
let fetchPostURL = '/api/posts';

async function loadFilterCategories() {
    try {
        const response = await fetch('/api/categories');
        const categories = await response.json();
        const container = document.getElementById('categoryFilter');
        
        // category icons
        const categoryIcons = {
            "Technology": "fas fa-laptop-code",
            "Sports": "fas fa-football-ball",
            "Entertainment": "fas fa-film",
            "Science": "fas fa-flask",
            "Politics": "fas fa-landmark",
            "Health": "fas fa-heartbeat",
            "Travel": "fas fa-plane",
            "Food": "fas fa-utensils",
            "Other": "fas fa-question-circle"
        };

        // All Categories button
        container.innerHTML = `
            <button 
                class="category-filter-btn active"
                data-category="">
                <i class="fas fa-list"></i> All Posts
            </button>
        `;

        // Tother categories buttons
        container.innerHTML += categories.map(category => `
            <button 
                class="category-filter-btn"
                data-category="${category.id}">
                <i class="${categoryIcons[category.name] || 'fas fa-tag'}"></i> ${category.name}
            </button>
        `).join('');

        // active state on click management
        container.querySelectorAll(".category-filter-btn").forEach(button => {
            button.addEventListener("click", function() {

                container.querySelectorAll(".category-filter-btn").forEach(btn => btn.classList.remove("active"));
                
                this.classList.add("active");
                
                if (this.dataset.category === 'all') {
                    resetPosts();
                } else {
                    applyCategoryFilter(this.dataset.category);
                }
            });
        });

        fetchPosts();

    } catch (error) {
        console.error('Error loading categories:', error);
    }
}

async function loadPostCategories() {
    try {
        const response = await fetch('/api/categories');
        const categories = await response.json();
        const container = document.getElementById('postCategories');
        
        container.innerHTML = categories.map(category => `
            <label class="category-checkbox">
                <input type="checkbox" value="${category.id}">
                <span>${category.name}</span>
            </label>
        `).join('');
    } catch (error) {
        console.error('Error loading categories:', error);
    }
}

async function openCreatePostModal() {
    try {
        const response = await fetch('/api/protected/api/auth/status');
        const data = await response.json();

        if (!data.authenticated) {
            handleError('Please login to create post');
            return;
        }

        // If authenticated, show modal and load categories
        document.getElementById('createPostModal').classList.add('active');
        loadPostCategories();
    } catch (error) {
        console.error('Error checking auth status:', error);
        handleError('Please login to create a post');
    }
}

function createPostElement(post) {
    const postDiv = document.createElement('div');
    postDiv.className = 'post';
    postDiv.dataset.postId = post.id;
    
    // Truncate content if it's longer than 800 characters
    const isLongPost = post.content.length > 800;
    const truncatedContent = isLongPost ? post.content.slice(0, 800) + '...' : post.content;
    
    postDiv.innerHTML = `
        <div class="post-header">
            <h3>${post.title}</h3>
            <small class="post-meta-info">Posted by ${post.nickname} on ${new Date(post.created_at).toLocaleString()}</small>
        </div>
        <div class="post-content">
            <p>${truncatedContent}</p>
            ${isLongPost ? `
                <div class="read-more-section">
                    <button onclick="toggleFullPost(this, \`${encodeURIComponent(post.content)}\`)" class="read-more-btn">
                        Read More
                    </button>
                </div>
            ` : ''}
        </div>
        <div class="post-footer">
            <div class="categories">
                ${post.categories ? post.categories.map(cat => 
                    `<span class="category-tag">${cat.name}</span>`
                ).join('') : ''}
            </div>
            <div class="post-actions">
                <button onclick="handleLike(${post.id}, true)" class="like-btn">
                    👍 <span>${post.likes || 0}</span>
                </button>
                <button onclick="handleLike(${post.id}, false)" class="dislike-btn">
                    👎 <span>${post.dislikes || 0}</span>
                </button>
                <button onclick="toggleComments(${post.id})" class="comment-btn">
                    💬 Comments (${(post.comments || []).length})
                </button>
            </div>
        </div>
        <div class="comments-container" id="comments-container-${post.id}" style="display: none;">
            <div class="comments-section" id="comments-${post.id}">
                ${renderComments(post.comments || [])}
            </div>
            <button onclick="showCommentForm(${post.id})" class="add-comment-btn">Add Comment</button>
            <div class="comment-form" id="comment-form-${post.id}" style="display: none;">
                <textarea placeholder="Write a comment..."></textarea>
                <button onclick="submitComment(${post.id})">Submit</button>
            </div>
        </div>
    `;
    return postDiv;
}

function toggleFullPost(button, content) {
    const decodedContent = decodeURIComponent(content);
    const postContent = button.closest('.post-content');
    const paragraph = postContent.querySelector('p');
    
    if (button.textContent === 'Read More') {
        paragraph.textContent = decodedContent;
        button.textContent = 'Show Less';
    } else {
        paragraph.textContent = decodedContent.slice(0, 800) + '...';
        button.textContent = 'Read More';
        // Scroll back to the top of the post
        button.closest('.post').scrollIntoView({ behavior: 'smooth' });
    }
}

async function fetchPosts(append = false) {
    if (isLoading || (!append && !hasMorePosts)) return;
    
    try {
        isLoading = true;
        const response = await fetch(`${fetchPostURL}?page=${currentPage}&category=${currentCategory}&filter=${currentFilter}`);
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to fetch posts');
        }
        const posts = await response.json();
        const mainContent = document.getElementById('main-content');
         // Ensure posts-list exists inside main-content
         let postsList = document.getElementById('posts-list');
         if (!postsList) {
             postsList = document.createElement('div');
             postsList.id = 'posts-list';
             mainContent.innerHTML= '';
             mainContent.appendChild(postsList);
            }
        // const postsList = document.getElementById('posts-list');
        if (!posts || posts.length === 0) {
            if (!append) {
                postsList.innerHTML = '<p>No posts yet. Be the first to create one!</p>';
            }
            hasMorePosts = false;
            return;
        }
        // Check if we've reached the end
        if (posts.length < 8) {
            hasMorePosts = false;
        }
        
        if (!append) {
            postsList.innerHTML = '';
        }

        posts.forEach(post => {
            const postElement = createPostElement(post);
            postsList.appendChild(postElement);
        });

        // Re-setup infinite scroll after adding new posts
        setupInfiniteScroll();

    } catch (error) {
        console.error('Error fetching posts:', error);    
        handleError(error.message);
    } finally {
        isLoading = false;
    }
}

function closeCreatePostModal() {
    document.getElementById('createPostModal').classList.remove('active');
    document.getElementById('createPostForm').reset();
}

// post and title validation
function validatePostForm(title, content) {
    const errorDiv = document.getElementById('post-error-message');
    
    if (!title || title.trim() === '') {
        showPostError('Please enter a post title');
        return false;
    }
    
    if (!content || content.trim() === '') {
        showPostError('Please enter post content');
        return false;
    }
    
    errorDiv.style.display = 'none';
    return true;
}

// Add this helper function for consistent error styling
function showPostError(message) {
    const errorDiv = document.getElementById('post-error-message');
    errorDiv.innerHTML = `
        <div class="error-content">
            <i class="fas fa-exclamation-circle"></i>
            <span>${message}</span>
        </div>
    `;
    errorDiv.style.display = 'block';
}

// Modify the createPost function to include validation
async function handleCreatePost(event) {
    event.preventDefault();
    
    const title = document.getElementById('postTitle').value;
    const content = document.getElementById('postContent').value;
    
    const selectedCategories = Array.from(document.querySelectorAll('#postCategories input:checked')).map(input => parseInt(input.value));

    try {
        const response = await fetch('/api/protected/api/posts/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ title, content, raw_categories: selectedCategories })
        });
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to create post');
        }
        handleSuccess('Post created successfully');
        closeCreatePostModal();
        resetPosts();
    } catch (e) {
        if (e.message === 'Title and content are required') {
            showPostError('Title and content are required')
            return
        }
        handleError(e.message)
    }
}

async function handleLike(postId, isLike) {
    try {
        const authResponse = await fetch('/api/protected/api/auth/status');
        const authData = await authResponse.json();
        
        if (!authData.authenticated) {
            handleError('Please login to like posts');
            return;
        }

        const response = await fetch('/api/protected/api/likes', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                post_id: postId,
                is_like: isLike
            })
        });

        if (!response.ok) {
            throw new Error('Login to leave a like');
        }

        await updatePost(postId);
        handleSuccess(isLike ? 'Post liked!' : 'Post disliked!');
    } catch (error) {
        console.error('Error handling like:', error);
        handleError('Login to like posts');
    }
}

function setupInfiniteScroll() {
    // Remove any existing trigger
    const existingTrigger = document.getElementById('load-more-trigger');
    if (existingTrigger) {
        existingTrigger.remove();
    }

    const loadMoreTrigger = document.createElement('div');
    loadMoreTrigger.id = 'load-more-trigger';
    loadMoreTrigger.innerHTML = `
        <div class="loading-spinner" style="display: none;"></div>
        <p class="no-more-posts" style="display: none;">No more posts to load</p>
    `;
    document.getElementById('posts-list').appendChild(loadMoreTrigger);

    updateNoMorePostsVisibility();

    // Set up intersection observer for infinite scroll
    const observer = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting && hasMorePosts && !isLoading) {
            loadMorePosts();
        } else if (entries[0].isIntersecting && !hasMorePosts) {
            updateNoMorePostsVisibility();
        }
    }, { threshold: 0.1 });
    observer.observe(loadMoreTrigger);
}

async function loadMorePosts() {
    if (isLoading || !hasMorePosts) return;
    
    const trigger = document.getElementById('load-more-trigger');
    const spinner = trigger.querySelector('.loading-spinner');
    
    spinner.style.display = 'inline-block';
    
    currentPage++;
    await fetchPosts(true);
    
    spinner.style.display = 'none';
    
    updateNoMorePostsVisibility();
}

// updates the visibility of the "no more posts" message
function updateNoMorePostsVisibility() {
    const trigger = document.getElementById('load-more-trigger');
    if (!trigger) return;

    const noMorePosts = trigger.querySelector('.no-more-posts');
    if (!hasMorePosts) {
        noMorePosts.style.display = 'block';
    } else {
        noMorePosts.style.display = 'none';
    }
}

// Update the existing functions
function resetPosts() {
    currentPage = 1;
    currentFilter = '';
    currentCategory = '';
    hasMorePosts = true;
    fetchPosts(false);
}

// Called when the page loads
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    resetPosts();
    setupInfiniteScroll();
    loadFilterCategories();
});

// Update other functions that fetch posts to use resetPosts()
function applyFilters(filter) {
    currentFilter = filter || '';
    currentCategory = '';
    currentPage = 1;
    hasMorePosts = true;
    fetchPostURL = 'api/protected/api/posts';
    fetchPosts(false);
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
                content: content
            })
        });

        if (!response.ok) {
            throw new Error('Failed to submit comment');
        }

        await response.json();
        textarea.value = '';
        commentForm.style.display = 'none';
        
        await updatePost(postId);
        handleSuccess('Comment posted successfully');
    } catch (error) {
        console.error('Error submitting comment:', error);
        handleError('Login to post comment');
    }
}

async function updatePost(postId) {
    try {
        const response = await fetch(`/api/posts/${postId}`);
        if (!response.ok) {
            throw new Error('Failed to fetch updated post');
        }
        const updatedPost = await response.json();
        
        // Find and update the specific post in the DOM
        const existingPost = document.querySelector(`.post[data-post-id="${postId}"]`);
        if (existingPost) {
            const newPost = createPostElement(updatedPost);
            // Preserve the comments display state
            const oldCommentsContainer = existingPost.querySelector('.comments-container');
            const newCommentsContainer = newPost.querySelector('.comments-container');
            if (oldCommentsContainer && oldCommentsContainer.style.display === 'block') {
                newCommentsContainer.style.display = 'block';
            }
            existingPost.replaceWith(newPost);
        }
    } catch (error) {
        console.error('Error updating post:', error);
    }
}

// Add this function to handle category filtering
function applyCategoryFilter(categoryId) {
    currentCategory = categoryId === '' ? '' : `category-${categoryId}`;
    currentFilter = '';
    currentPage = 1;
    hasMorePosts = true;
    fetchPostURL = 'api/posts';
    fetchPosts(false);
}
