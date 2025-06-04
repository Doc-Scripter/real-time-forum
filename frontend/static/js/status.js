async function fetchAndDisplayOnlineUsers() {
    try {
        const response = await fetch('/api/protected/api/status');
        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }
        let users = await response.json();
        if (!Array.isArray(users)) users = [];
        
        if (!users.length) {
            const trendingSection = document.querySelector('.trending-section');
            trendingSection.innerHTML = `<h3>Users</h3><div class="online-users-list">No users found</div>`;
            return;
        }

        // Separate users into recent conversations and others
        const recentConversations = users.filter(user => user.last_conversation);
        const otherUsers = users.filter(user => !user.last_conversation);

        const trendingSection = document.querySelector('.trending-section');
        trendingSection.innerHTML = `
            <h3>Users</h3>
            <div class="online-users-list">
                ${recentConversations.length > 0 ? `
                    <div class="user-section">
                        <div class="section-header">Recent Conversations</div>
                        ${recentConversations.map(user =>
                            `<div class="online-user" 
                                 style="cursor:pointer;" 
                                 data-receiver-id="${user.receiver}" 
                                 onclick="openInboxWithUser('${user.nickname}',${user.receiver})">
                                <span class="status-dot ${user.online ? 'online' : 'offline'}"></span>
                                <span class="receiver">${user.nickname}</span>
                                <span class="conversation-indicator">💬</span>
                            </div>
                        `).join('')}
                    </div>
                ` : ''}
                ${otherUsers.length > 0 ? `
                    <div class="user-section">
                        <div class="section-header">All Users</div>
                        ${otherUsers.map(user =>
                            `<div class="online-user" 
                                 style="cursor:pointer;" 
                                 data-receiver-id="${user.receiver}" 
                                 onclick="openInboxWithUser('${user.nickname}',${user.receiver})">
                                <span class="status-dot ${user.online ? 'online' : 'offline'}"></span>
                                <span class="receiver">${user.nickname}</span>
                            </div>
                        `).join('')}
                    </div>
                ` : ''}
            </div>
        `;
    } catch (error) {
    }
}

// Enhanced CSS for better visual separation
const styling = document.createElement('style');
styling.innerHTML = `
.online-users-list { display: flex; flex-direction: column; gap: 8px; }
.online-user { display: flex; align-items: center; gap: 8px; padding: 4px 8px; border-radius: 4px; }
.online-user:hover { background-color: #f0f0f0; }
.status-dot { width: 10px; height: 10px; border-radius: 50%; display: inline-block; }
.status-dot.online { background: #4caf50; }
.status-dot.offline { background: #ccc; }
.user-section { margin-bottom: 16px; }
.section-header { font-weight: bold; font-size: 12px; color: #666; margin-bottom: 8px; text-transform: uppercase; }
.conversation-indicator { font-size: 12px; margin-left: auto; }
`;
document.head.appendChild(styling);

