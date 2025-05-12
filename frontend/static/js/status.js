async function fetchAndDisplayOnlineUsers() {
    try {
        const response = await fetch('/api/protected/api/status');
        if (!response.ok) {
            throw new Error('Failed to fetch online users');
        }
        let users = await response.json();
        if (!Array.isArray(users)) users = [];
        if (!users.length ) {
            const trendingSection = document.querySelector('.trending-section');
            trendingSection.innerHTML = `<h3>Users</h3><div class="online-users-list">No users online</div>`;
            return;
        }


        const trendingSection = document.querySelector('.trending-section');
        trendingSection.innerHTML = `
            <h3>Users</h3>
            <div class="online-users-list">
                ${users.map(user =>
                    `<div class="online-user" 
                         style="cursor:pointer;" 
                         data-receiver-id="${user.receiver}" 
                         onclick="openInboxWithUser('${user.username}',${user.receiver})">
                        <span class="status-dot ${user.online ? 'online' : 'offline'}"></span>
                        <span class="receiver">${user.username}</span>
                    </div>
                `).join('')}
            </div>
        `;
    } catch (error) {
        console.error('Error fetching online users:', error);
    }
}

// CSS for status dots (add to your CSS file or inject as needed)
const styling = document.createElement('style');
styling.innerHTML = `
.online-users-list { display: flex; flex-direction: column; gap: 8px; }
.online-user { display: flex; align-items: center; gap: 8px; }
.status-dot { width: 10px; height: 10px; border-radius: 50%; display: inline-block; }
.status-dot.online { background: #4caf50; }
.status-dot.offline { background: #ccc; }
`;
document.head.appendChild(styling);

// Update online users every 30 seconds
setInterval(fetchAndDisplayOnlineUsers, 5000);

// Initial load
document.addEventListener('DOMContentLoaded', fetchAndDisplayOnlineUsers);