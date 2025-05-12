let ws;

const inboxBtn = document.getElementById("inbox-btn");
const mainContent = document.getElementById("main-content");

let messages = [];
let lastMessages = [];

let currentUser = null; 
let currentUserId = null;
let currentReceiverId = null;
let currentPartner = null;

// Fetch current user from auth status endpoint
async function fetchCurrentUser() {
  try {
    const response = await fetch("/api/protected/api/auth/status");
    const data = await response.json();
    if (data.authenticated) {
      currentUser = data.username;
      currentUserId = data.user_id;
    } else {
      currentUser = null;
      currentUserId = null;
    }
  } catch (e) {
    console.error("Failed to fetch or parse current user:", e);
    currentUser = null;
    currentUserId = null;
  }
}

// Fetch messages from API
async function fetchMessages(receiverId) {
  try {
    const url = new URL("/api/protected/api/messages", window.location.origin);
    url.searchParams.append("receiver", receiverId);
    let res = await fetch(url.toString());
    if (res.ok) {
      messages = await res.json();
    } else {
      messages = [];
    }
  } catch (e) {
    console.error("Failed to fetch messages:", e);
    messages = [];
  }
}


async function fetchLastMessages() {
  try {
    let res = await fetch("/api/protected/api/messages");
    if (res.ok) {
      lastMessages = await res.json();
    } else {
      lastMessages = [];
    }
  } catch (e) {
    console.error("Failed to fetch messages:", e);
    lastMessages = [];
  }
}

let inboxRefreshInterval;

function startInboxRefresh() {
  // Clear any existing interval
  if (inboxRefreshInterval) {
    clearInterval(inboxRefreshInterval);
  }
  
  // Set new interval to fetch messages and re-render inbox every 5 seconds
  inboxRefreshInterval = setInterval(async () => {
    // Only refresh if inbox view is currently open
    if (mainContent && mainContent.querySelector('.inbox-section')) {
      await fetchLastMessages();
      renderInbox();
    }
  }, 5000);
}

function stopInboxRefresh() {
  if (inboxRefreshInterval) {
    clearInterval(inboxRefreshInterval);
    inboxRefreshInterval = null;
  }
} 

// Example: Initialize currentUser before using inbox functionality
async function initInbox() {
  await fetchCurrentUser();
  await fetchLastMessages();
  // Now currentUser is set and can be used throughout inbox logic
}

// Helper to group messages by chat partner
function getConversations() {
  if (!Array.isArray(lastMessages) || lastMessages.length === 0) {
    return [];
  }

  console.log("DEBUG - Last Messages:", lastMessages);
  console.log("DEBUG - Current User:", currentUser);

  return lastMessages.map((msg) => {
    console.log("DEBUG - Processing message:", msg);
    
    // Log values before deciding partner
    console.log("DEBUG - Sender:", msg.sender);
    console.log("DEBUG - Receiver:", msg.receiver);
    let receiverName = msg.receiver;
        const onlineUsers = document.querySelectorAll('.online-user');
        onlineUsers.forEach(user => {
            if (user.dataset.receiverId == msg.receiver) {
                receiverName = user.querySelector('.receiver').textContent;
            }
        });

    
    const isCurrentUserSender = msg.sender === currentUser;
    const partner = isCurrentUserSender ? receiverName : msg.sender;
    
    console.log("DEBUG - Determined partner:", partner);

    return {
      partner: partner,
      lastMsg: {
        data: msg.data || "",
        time: msg.time || "",
        receiver:receiverName, // Changed from msg.sender to msg.receiver
        receiverId: msg.receiver,
      },
    };
  });
}
// Render inbox: conversation list or chat view
async function renderInbox() {
  if (!currentUser) {
    mainContent.innerHTML = "<div>Please log in to view your inbox.</div>";
    return;
  }

  
  await fetchLastMessages()
  console.log("DEBUG - After fetchLastMessages:", lastMessages);

  startInboxRefresh()

  // Show conversation list
  const conversations = getConversations();
  console.log("DEBUG - Final conversations:", conversations);
  let inboxHTML = `
    <div class="inbox-section">
        <h2>Inbox</h2>
        <div class="conversation-list">
            ${
              conversations.length === 0
                ? "<div>No conversations yet.</div>"
                : conversations
                    .map((conv) => {
                      const partner = conv.partner || "";
                      const receiverId =
                        conv.lastMsg && conv.lastMsg.receiverId
                          ? parseInt(conv.lastMsg.receiverId)
                          : 0;
                      const lastMsgText = conv.lastMsg
                        ? conv.lastMsg.data || conv.lastMsg.data || ""
                        : "";
                      const lastMsgTime =
                        conv.lastMsg && conv.lastMsg.time
                          ? conv.lastMsg.time
                          : "";
                      return `
                    <div class="conversation-item" 
                         data-username="${partner}" 
                         data-receiver-id="${receiverId}">
                        <div class="avatar">${
                          partner[0] ? partner[0].toUpperCase() : "?"
                        }</div>
                        <div class="conv-details">
                            <div class="conv-name">${partner}</div>
                            <div class="conv-preview">
                                ${
                                  lastMsgText.length > 32
                                    ? lastMsgText.slice(0, 32) + "…"
                                    : lastMsgText
                                }
                            </div>
                            <div class="conv-time">${lastMsgTime}</div>
                        </div>
                    </div>
                `;
                    })
                    .join("")
            }
        </div>
    </div>
`;
  mainContent.innerHTML = inboxHTML;

  // Delegated click handler for conversation items
  const convList = mainContent.querySelector(".conversation-list");
  if (convList) {
    convList.onclick = (e) => {
      const item = e.target.closest(".conversation-item");
      if (item) {
        const username = item.dataset.username;
        const receiverId = parseInt(item.dataset.receiverId) || 0;
        renderChat(username, receiverId);
      }
    };
  }
  

}

// Clean up on page unload
window.addEventListener('beforeunload', stopInboxRefresh);

// Render chat with selected user
async function renderChat(partner, receiverId) {
  startInboxRefresh();
    currentPartner = partner;
    currentReceiverId = receiverId;
  if (!currentUser || !receiverId) {
    console.error(
      "renderChat called with invalid partner/ID",
      // partnerUsername,
      partner
    );
    renderInbox(); // Go back to inbox if details are missing
    return;
  }
  await fetchMessages(receiverId);
  if (!Array.isArray(messages)) {
    console.error("Messages is not an array", messages);
    messages = [];
  }
  const chatMessages = messages.filter((msg) => {
    if (!msg) return false;
    return (
      (msg.sender === currentUser && msg.receiver === receiverId) ||
      (msg.sender === partner && msg.receiver === currentUserId)
    );
  });
  console.log(chatMessages);
  let chatHTML = `
        <div class="chat-section">
            <button id="back-to-inbox" class="back-btn">← Back to Inbox</button>
            <h2>Chat with ${partner}</h2>
            <div class="chat-messages">
                ${
                  chatMessages.length === 0
                    ? '<div class="no-messages">No messages yet. Start the conversation!</div>'
                    : chatMessages
                        .map(
                          (msg) => `
                    <div class="chat-msg ${
                      msg.sender === currentUser ? "sent" : "received"
                    }">
                        <div class="msg-content">
                            <span class="msg-text">${msg.data}</span>
                            <span class="msg-time">${msg.time}</span>
                        </div>
                    </div>
                `
                        )
                        .join("")
                }
            </div>
            <form id="send-msg-form" class="send-msg-form">
                <input type="text" id="msg-input" placeholder="Type a message..." autocomplete="off" required />
                <button type="submit">Send</button>
            </form>
        </div>
    `;
  mainContent.innerHTML = chatHTML;

  const messagesContainer = document.querySelector(".chat-messages");
  if (messagesContainer) {
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  document.getElementById("back-to-inbox").onclick = () => renderInbox();

  document.getElementById("send-msg-form").onsubmit = async (e) => {
    e.preventDefault();
    const input = document.getElementById("msg-input");
    const text = input.value.trim();
    console.log("sending message: ", text);
    if (!text) return;

    if (!receiverId) {
      console.error("Missing receiver ID for", partner);
      return;
    }
    // Send via WebSocket
    ws.send(
      JSON.stringify({
        sender: currentUser,
        receiver: receiverId,
        data: text,
      })
    );

    // Also store via REST API
    try {
      const response = await fetch("/api/protected/api/messages", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          receiver: receiverId,
          data: text,
        }),
      });

      if (response.ok) {
        // Add the sent message to the local messages array
        messages.push({
          sender: currentUser,
          receiver: receiverId,
          data: text,
          time: new Date().toLocaleTimeString(),
        });

        // Re-render the chat to show the new message
        renderChat(partner, receiverId);
      } else {
        console.error("Failed to save message");
      }
    } catch (error) {
      console.error("Error saving message:", error);
    }

    input.value = "";
  };


}


// Initialize WebSocket connection
function initWebSocket() {
  ws = new WebSocket("/api/messaging"); // Adjust URL as needed

  ws.onopen = () => {
    console.log("WebSocket connected");
    // Optionally authenticate or identify user
    ws.send(JSON.stringify({ type: "auth", user: currentUser }));
  };

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log("Received message:", data);

      if (data.type === "message" && data.message) {
        // Ensure message has all required fields
        const message = {
          sender: data.message.sender,
          data: data.message.data,
          time: data.message.time || new Date().toLocaleTimeString(),
          receiver: data.message.receiver,
        };

        if (
          (message.sender === currentUser &&
            message.receiver === currentReceiverId) ||
          (message.sender === currentPartner &&
            message.receiver === currentUserId)
        ) {
          messages.push(message);

          // Re-render if viewing this conversation
          const currentChat = document.querySelector('.chat-section h2');
          if (currentChat && currentPartner && 
              (message.sender === currentPartner || message.receiver === currentPartner)) {
              // Re-render the chat with updated messages
              renderChat(currentPartner, currentReceiverId);
          }
          
        }
      }
    } catch (error) {
      console.error("Error processing message:", error, event.data);
    }
  };

  ws.onclose = () => {
    console.log("WebSocket disconnected, retrying...");
    setTimeout(initWebSocket, 2000); // Reconnect
  };
}

// Send message via WebSocket
async function sendMessage(receiverId, text) {
  if (!receiverId || receiverId === 0) {
    console.error("Invalid receiver ID:", receiverId);
    return;
  }
  if (ws && ws.readyState === WebSocket.OPEN) {
    const messageData = {
      type: "message",
      sender: currentUser,
      receiver: receiverId,
      data: text,
      time: new Date().toLocaleTimeString(),
      status: "Sent",
    };

    ws.send(JSON.stringify(messageData));

    // Also store via REST API
    try {
      const response = await fetch("api/protected/api/messages", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          receiver: receiverId,
          data: text,
        }),
      });

      if (response.ok) {
        console.log("Message saved successfully");
      } else {
        console.error("Failed to save message");
      }
    } catch (error) {
      console.error("Error saving message:", error);
    }
  } else {
    console.error("WebSocket not connected");
  }
}

// Attach event listener for inbox button
if (inboxBtn && mainContent) {
  inboxBtn.addEventListener("click", () => renderInbox());
}

// Expose function globally for status.js and conversation items
window.openInboxWithUser = function (username, receiverId) {
  renderChat(username, receiverId);
};

// On page load, fetch user info and initialize WebSocket
window.addEventListener("DOMContentLoaded", async () => {
  await initInbox();
  initWebSocket();
});

// Add minimal CSS for WhatsApp-like style (should be moved to a CSS file)
if (!document.getElementById("inbox-style")) {
  const inboxStyle = document.createElement("style");
  inboxStyle.id = "inbox-style";
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
    .chat-section { width: 100%; max-width: 500px; margin: 0 auto; background: #f0f0f0; border-radius: 8px; padding: 16px; display: flex; flex-direction: column; height: 80vh; }
    .chat-messages { flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 8px; padding: 10px 0; }
    .chat-msg { display: flex; margin-bottom: 8px; }
    .chat-msg.sent { justify-content: flex-end; }
    .chat-msg.received { justify-content: flex-start; }
    .msg-content { padding: 8px 12px; border-radius: 18px; max-width: 70%; }
    .sent .msg-content { background: #dcf8c6; }
    .received .msg-content { background: #fff; }
    .msg-text { display: block; margin-bottom: 4px; }
    .msg-time { font-size: 0.7em; color: #999; display: block; text-align: right; }
    .send-msg-form { display: flex; gap: 8px; margin-top: 12px; }
    .send-msg-form input { flex: 1; padding: 8px; border-radius: 20px; border: 1px solid #ddd; }
    .send-msg-form button { background: #128C7E; color: white; border: none; border-radius: 50%; width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; }
    .no-messages { text-align: center; color: #999; padding: 20px; }
    `;
  document.head.appendChild(inboxStyle);
}
