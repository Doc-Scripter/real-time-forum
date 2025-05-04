// Example selectors; update as needed for your HTML structure
const inboxBtn = document.getElementById('inbox-btn');
const mainContent = document.getElementById('main-content');

// Mock message data; replace with real API call as needed
const messages = [
    { id: 3, sender: 'Alice', text: 'Hey, how are you?', time: '10:30', self: false },
    { id: 4, sender: 'You', text: 'I am good, thanks!', time: '10:31', self: true },
    { id: 5, sender: 'Alice', text: 'Want to catch up later?', time: '10:32', self: false },
    { id: 6, sender: 'Bob', text: 'Hello!', time: '11:00', self: false },
    { id: 7, sender: 'You', text: 'Hi Bob!', time: '11:01', self: true },
];

// Helper to get chat partner name
function getChatPartner(msg) {
    return msg.self ? msg.sender : msg.sender;
}

// Render inbox, optionally filtered by username
function renderInbox(username = null) {
    let filteredMessages = messages;
    let chatTitle = 'Inbox';
    if (username) {
        filteredMessages = messages.filter(
            msg => msg.sender === username || (msg.self && username !== 'You')
        );
        chatTitle = `Chat with ${username}`;
    }
    const sortedMessages = filteredMessages.slice().sort((a, b) => b.id - a.id);
    let inboxHTML = `
        <div class="inbox-section">
            <h2>${chatTitle}</h2>
            <div class="messages-list">
                ${sortedMessages.map(msg => `
                    <div class="message-bubble ${msg.self ? 'self' : 'other'}">
                        <div class="message-sender">${msg.self ? 'You' : msg.sender}</div>
                        <div class="message-text">${msg.text}</div>
                        <div class="message-time">${msg.time}</div>
                    </div>
                `).join('')}
            </div>
        </div>
    `;
    mainContent.innerHTML = inboxHTML;
}

// Attach event listener for inbox button
if (inboxBtn && mainContent) {
    inboxBtn.addEventListener('click', () => renderInbox());
}

// Expose function globally for status.js
window.openInboxWithUser = function(username) {
    renderInbox(username);
};

// Add minimal CSS for WhatsApp-like style (should be moved to a CSS file)
if (!document.getElementById('inbox-style')) {
    const inboxStyle = document.createElement('style');
    inboxStyle.id = 'inbox-style';
    inboxStyle.innerHTML = `
    .inbox-section { width: 100%; max-width: 500px; margin: 0 auto; background: #f0f0f0; border-radius: 8px; padding: 16px; }
    .messages-list { display: flex; flex-direction: column; gap: 12px; }
    .message-bubble { padding: 10px 14px; border-radius: 18px; max-width: 80%; }
    .message-bubble.self { background: #dcf8c6; align-self: flex-end; }
    .message-bubble.other { background: #fff; align-self: flex-start; }
    .message-sender { font-size: 0.85em; color: #555; }
    .message-text { margin: 4px 0; }
    .message-time { font-size: 0.75em; color: #999; text-align: right; }
    `;
    document.head.appendChild(inboxStyle);
}