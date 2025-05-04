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

// Helper to group messages by chat partner
function getConversations() {
    const convMap = {};
    messages.forEach(msg => {
        const partner = msg.self ? msg.sender : msg.sender;
        if (partner === 'You') return; // skip self
        if (!convMap[partner]) convMap[partner] = [];
        convMap[partner].push(msg);
    });
    // Convert to array and sort by last message id (descending)
    return Object.entries(convMap)
        .map(([partner, msgs]) => ({
            partner,
            lastMsg: msgs.slice().sort((a, b) => b.id - a.id)[0]
        }))
        .sort((a, b) => b.lastMsg.id - a.lastMsg.id);
}

// Render inbox: conversation list or chat view
function renderInbox(username = null) {
    if (!username) {
        // Show conversation list
        const conversations = getConversations();
        let inboxHTML = `
            <div class="inbox-section">
                <h2>Inbox</h2>
                <div class="conversation-list">
                    ${conversations.length === 0 ? '<div>No conversations yet.</div>' : conversations.map(conv => `
                        <div class="conversation-item" onclick="openInboxWithUser('${conv.partner}')">
                            <div class="avatar">${conv.partner[0].toUpperCase()}</div>
                            <div class="conv-details">
                                <div class="conv-name">${conv.partner}</div>
                                <div class="conv-preview">${conv.lastMsg.text.length > 32 ? conv.lastMsg.text.slice(0, 32) + '…' : conv.lastMsg.text}</div>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
        mainContent.innerHTML = inboxHTML;
        return;
    }

    // Show chat with specific user
    let filteredMessages = messages.filter(
        msg => msg.sender === username || (msg.self && username !== 'You')
    );
    const sortedMessages = filteredMessages.slice().sort((a, b) => b.id - a.id);
    let inboxHTML = `
        <div class="inbox-section">
            <button class="back-btn" onclick="renderInbox()" style="margin-bottom:10px;">← Back</button>
            <h2>Chat with ${username}</h2>
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

// Expose function globally for status.js and conversation items
window.openInboxWithUser = function(username) {
    renderInbox(username);
};

// Add minimal CSS for WhatsApp-like style (should be moved to a CSS file)
if (!document.getElementById('inbox-style')) {
    const inboxStyle = document.createElement('style');
    inboxStyle.id = 'inbox-style';
    inboxStyle.innerHTML = `
    .inbox-section { width: 100%; max-width: 500px; margin: 0 auto; background: #f0f0f0; border-radius: 8px; padding: 16px; }
    .conversation-list { display: flex; flex-direction: column; gap: 12px; }
    .conversation-item { display: flex; align-items: center; gap: 12px; padding: 10px; border-radius: 10px; background: #fff; cursor: pointer; transition: background 0.2s; }
    .conversation-item:hover { background: #e6e6e6; }
    .avatar { width: 40px; height: 40px; border-radius: 50%; background: #bbb; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 1.3em; font-weight: bold; }
    .conv-details { flex: 1; }
    .conv-name { font-weight: bold; }
    .conv-preview { color: #666; font-size: 0.95em; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    .messages-list { display: flex; flex-direction: column; gap: 12px; }
    .message-bubble { padding: 10px 14px; border-radius: 18px; max-width: 80%; }
    .message-bubble.self { background: #dcf8c6; align-self: flex-end; }
    .message-bubble.other { background: #fff; align-self: flex-start; }
    .message-sender { font-size: 0.85em; color: #555; }
    .message-text { margin: 4px 0; }
    .message-time { font-size: 0.75em; color: #999; text-align: right; }
    .back-btn { background: none; border: none; color: #007bff; font-size: 1em; cursor: pointer; }
    `;
    document.head.appendChild(inboxStyle);
}