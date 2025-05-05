// Example selectors; update as needed for your HTML structure
const inboxBtn = document.getElementById('inbox-btn');
const mainContent = document.getElementById('main-content');

// Mock message data; replace with real API call as needed
const messages = [];
let ws;
// let currentUser = null; // Set this after fetching user info

// Fetch messages from API
async function fetchMessages() {
    const res = await fetch('/api/messages');
    if (res.ok) {
        messages = await res.json();
    } else {
        messages = [];
    }
}

// Helper to group messages by chat partner
function getConversations() {
    const convMap = {};
    messages.forEach(msg => {
        // Assuming msg.sender and msg.receiver fields
        const partner = msg.Sender ;
        // if (partner === currentUser) return;
        if (!convMap[partner]) convMap[partner] = [];
        convMap[partner].push(msg);
    });
    return Object.entries(convMap)
        .map(([partner, msgs]) => ({
            partner,
            lastMsg: msgs.slice().sort((a, b) => b.id - a.id)[0]
        }))
        .sort((a, b) => b.lastMsg.id - a.lastMsg.id);
}

// Render inbox: conversation list or chat view
function renderInbox(username = null, receiverId = null) {
    if (!username) {
        // Show conversation list
        const conversations = getConversations();
        let inboxHTML = `
            <div class="inbox-section">
                <h2>Inbox</h2>
                <div class="conversation-list">
                    ${conversations.length === 0 ? '<div>No conversations yet.</div>' : conversations.map(conv => `
                        <div class="conversation-item" 
                             data-username="${conv.partner}" 
                             data-receiver-id="${conv.lastMsg.receiver}">
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

        // Delegated click handler for conversation items
        const convList = mainContent.querySelector('.conversation-list');
        if (convList) {
            convList.onclick = (e) => {
                const item = e.target.closest('.conversation-item');
                if (item) {
                    const username = item.dataset.username;
                    const receiverId = item.dataset.receiverId;
                    openInboxWithUser(username, receiverId);
                }
            };
        }
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
            <form id="send-message-form" style="display:flex;gap:8px;margin-top:12px;">
                <input id="message-input" type="text" placeholder="Type a message..." style="flex:1;" autocomplete="off"/>
                <button type="submit">Send</button>
            </form>
        </div>
    `;
    mainContent.innerHTML = inboxHTML;

    // Attach send handler
    document.getElementById('send-message-form').onsubmit = (e) => {
        e.preventDefault();
        const input = document.getElementById('message-input');
        const text = input.value.trim();
        if (text) {
            sendMessage(receiverId || username, text); // Prefer receiverId if available
            input.value = '';
        }
    };
}

// Initialize WebSocket connection
function initWebSocket() {
    ws = new WebSocket('ws://localhost:8080/ws'); // Adjust URL as needed

    ws.onopen = () => {
        console.log('WebSocket connected');
        // Optionally authenticate or identify user
        ws.send(JSON.stringify({ type: 'auth', user: currentUser }));
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'message') {
            messages.push(data.message);
            // If chat with sender is open, re-render
            const openChat = document.querySelector('.inbox-section h2')?.textContent?.includes(data.message.sender);
            if (openChat) renderInbox(data.message.sender);
        }
    };

    ws.onclose = () => {
        console.log('WebSocket disconnected, retrying...');
        setTimeout(initWebSocket, 2000); // Reconnect
    };
}

// Send message via WebSocket
function sendMessage(receiver, text) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
            type: 'message',
            sender: currentUser,
            receiver,
            text,
            time: new Date().toLocaleTimeString()
        }));
    }
}

// Attach event listener for inbox button
if (inboxBtn && mainContent) {
    inboxBtn.addEventListener('click', () => renderInbox());
}

// Expose function globally for status.js and conversation items
window.openInboxWithUser = function(username, receiverId) {
    renderInbox(username, receiverId);
};

// On page load, fetch user info and initialize WebSocket
window.addEventListener('DOMContentLoaded', async () => {
    // Fetch current user (implement as needed)
    const res = await fetch('/api/messages');
    if (res.ok) {
        currentUser = (await res.json()).sender;
        initWebSocket();
    }
});

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