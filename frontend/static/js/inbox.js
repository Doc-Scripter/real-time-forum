const inboxBtn = document.getElementById("inbox-btn");
const mainContent = document.getElementById("main-content");

let messages = [];
let lastMessages = [];

let currentUser = null;
let currentUserId = null;
let currentReceiverId = null;
let currentPartner = null;
let unreadCheckInterval;
let messageCache = {};
let messagesPerPage = 10; // Number of messages to display per page
let currPage = 1; // Current page number
let ws;
let endIndex = 0;
let lastLoadMoreCall = 0; // Add this variable near the top with other global variables



async function markRead(receiverId) {
  try {
    await fetch("/api/protected/api/messages", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        senderId: receiverId,
      }),
    });
    
    // Update unread count after marking messages as read
    checkUnreadMessages();
  } catch (error) {
  }
  
}

async function checkUnreadMessages() {
  try {
    // Build the URL with optional exclude_sender parameter
    let url = "/api/protected/api/unread";
    if (currentReceiverId) {
      url += `?exclude_sender=${currentReceiverId}`;
    }
    
    const response = await fetch(url);
    if (response.ok) {
      const count = await response.json();
      const badge = document.getElementById("unread-badge");

      if (count > 0) {
        badge.style.display = "block";
        badge.textContent = count;
      } else {
        badge.style.display = "none";
      }
    }
  } catch (error) {
  }
}

// Start checking for unread messages
function startUnreadCheck() {
  checkUnreadMessages(); 
  unreadCheckInterval = setInterval(checkUnreadMessages, 5000);
}

// Stop checking for unread messages
function stopUnreadCheck() {
  if (unreadCheckInterval) {
    clearInterval(unreadCheckInterval);
  }
}

// Fetch current user from auth status endpoint
async function fetchCurrentUser() {
  try {
    const response = await fetch("/api/protected/api/auth/status");
    const data = await response.json();
    ("INFO : fetched current user:",data)
    if (data.authenticated) {
      currentUser = data.nickname;
      currentUserId = data.user_id;
    } else {
      currentUser = null;
      currentUserId = null;
    }
  } catch (e) {
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
    if (mainContent && mainContent.querySelector(".inbox-section")) {
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

  ("DEBUG - Last Messages:", lastMessages);
  ("DEBUG - Current User:", currentUser);

  return lastMessages.map((msg) => {
    ("DEBUG - Processing message:", msg);

    // Log values before deciding partner
    ("DEBUG - Sender:", msg.sender);
    ("DEBUG - Receiver:", msg.receiver);
    let receiverName = msg.receiver;
    const onlineUsers = document.querySelectorAll(".online-user");
    onlineUsers.forEach((user) => {
      if (user.dataset.receiverId == msg.receiver) {
        receiverName = user.querySelector(".receiver").textContent;
      }
    });

    const isCurrentUserSender = msg.sender === currentUser;
    const partner = isCurrentUserSender ? receiverName : msg.sender;

    ("DEBUG - Determined partner:", partner);

    return {
      partner: partner,
      lastMsg: {
        data: msg.data || "",
        time: msg.time || "",
        receiver: receiverName, // Changed from msg.sender to msg.receiver
        receiverId: msg.receiver,
      },
    };
  });
}
// Render inbox: conversation list or chat view
async function renderInbox() {
  // Reset current chat variables when returning to inbox
  currentPartner = null;
  currentReceiverId = null;
  
  // initWebSocket();

  document.getElementById("toggleRight").style.display = "flex";
  document.getElementById("toggleLeft").style.display = "flex";
  if (!currentUser) {
    mainContent.innerHTML = "<div>Please log in to view your inbox.</div>";
    return;
  }

  await fetchLastMessages();
  ("DEBUG - After fetchLastMessages:", lastMessages);

  startInboxRefresh();

  // Show conversation list
  const conversations = getConversations();
  ("DEBUG - Final conversations:", conversations);
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

  // Update unread messages count when returning to inbox
  checkUnreadMessages();
}

// Clean up on page unload
window.addEventListener("beforeunload", () => {
  stopUnreadCheck();
  stopInboxRefresh();
});
// Render chat with selected user
async function renderChat(partner, receiverId) {
  ("DEBUG: Opening chat with:", partner);

  const leftSidebar = document.querySelector(".sidebar-left");
  const rightSidebar = document.querySelector(".sidebar-right");

  [leftSidebar, rightSidebar].forEach((sb) => {
    sb.classList.remove("active");
  });

  document.getElementById("toggleRight").style.display = "none";
  document.getElementById("toggleLeft").style.display = "none";

  document.getElementById("unread-badge").style.display = "none";

  startInboxRefresh();

  currentPartner = partner;
  currentReceiverId = receiverId;
  if (!currentUser || !receiverId) {
    console.error("renderChat called with invalid partner/ID", partner);
    renderInbox(); // Go back to inbox if details are missing
    return;
  }

  await fetchMessages(receiverId);

  if (!Array.isArray(messages)) {
    console.error("Messages is not an array", messages);
    messages = [];
  }
  (messages)
  
  
    messageCache[receiverId] = {
      messages: messages,
      currPage: 1,
    };
    messages = messages.slice(-10);

  const chatMessages = messages.filter((msg) => {
    if (!msg) return false;
    return (
      (msg.sender === currentUser && msg.receiver === receiverId) ||
      (msg.sender === partner && msg.receiver === currentUserId)
    );
  });

  // In the renderChat function, after the try-catch block that marks messages as read:
  try {
    await fetch("/api/protected/api/messages", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        senderId: receiverId,
      }),
    });
    
    // Update unread count after marking messages as read
    checkUnreadMessages();
  } catch (error) {
    console.error("Error marking messages as read:", error);
  }


  let chatHTML = `
        <div class="chat-section">
            <button id="back-to-inbox" class="back-btn">← Back to Inbox</button>
            <h2>Chat with ${partner}</h2>
            <button id="load-more-btn" class="load-more-btn">Load Previous Messages</button>
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
                        <div class="msg-details">
                        <span class="msg-sender">${msg.sender === currentUser ? currentUser : msg.sender}</span>
                            <span class="msg-time">${msg.time}</span>
                        </div>
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


  
  messagesContainer = document.querySelector(".chat-messages");

  if (messagesContainer) {
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  document.getElementById("back-to-inbox").onclick = () => {
    markRead(receiverId);
    renderInbox()
  };
  const loadMoreBtn = document.getElementById("load-more-btn");
  if (loadMoreBtn) {
    loadMoreBtn.onclick = () => loadMoreMessages(receiverId);
  }
  document.getElementById("send-msg-form").onsubmit = async (e) => {
    // endIndex + 1;
    e.preventDefault();
    const input = document.getElementById("msg-input");
    const text = input.value.trim();
    const sanitizedText=sanitizeHTML(text);
    ("sending message: ", sanitizedText);
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
        data: sanitizedText,
        time: new Date().toLocaleTimeString(),
      })
    );
    
    messageCache[receiverId].messages.push({
      sender: currentUser,
      receiver: receiverId,
      data: sanitizedText,
      time: new Date().toLocaleTimeString(),
    });

    // Also store via REST API
    try {
      const response = await fetch("/api/protected/api/messages", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          receiver: receiverId,
          data: sanitizedText,
        }),
      });
      const newMessage={
        sender: currentUser,
        receiver: receiverId,
        data: sanitizedText,
        time: new Date().toLocaleTimeString(),
      }
      if (response.ok) {
        // Add the sent message to the local messages array
        messages.push(newMessage);

        // Re-render the chat to show the new message
        // renderChat(partner, receiverId);

        appendNewMessage(newMessage)
      } else {
        console.error("Failed to save message");
      }
    } catch (error) {
      console.error("Error saving message:", error);
    }

    input.value = "";
  };
  // Add scroll event listener to load more messages
  const messagesContainerScroll = document.querySelector(".chat-messages");

  messagesContainerScroll.addEventListener("scroll", () => {
    if (messagesContainerScroll.scrollTop <=3) {
      loadMoreMessages(receiverId);
    }
  });
}

function loadMoreMessages(receiverId) {
  // Throttle to once per second
  const now = Date.now();
  if (now - lastLoadMoreCall < 1000) {
    return; // Exit early if called within 1 second
  }
  lastLoadMoreCall = now;
  
  const cache = messageCache[receiverId];
  const totalMessages = cache.messages.length;
  const alreadyDisplayed = cache.currPage * messagesPerPage;
  
  if (alreadyDisplayed < totalMessages) {
    cache.currPage++;
    
    // Calculate backwards from the end
    const endIndex = totalMessages - alreadyDisplayed;
    const startIndex = Math.max(0, endIndex - messagesPerPage);
    const messagesToPrepend = cache.messages.slice(startIndex, endIndex);

    const messagesContainer = document.querySelector(".chat-messages");
    const messagesHTML = messagesToPrepend
      .map(
        (msg) => `
      <div class="chat-msg ${msg.sender === currentUser ? "sent" : "received"}">
        <div class="msg-content">
        <span class="msg-text">${sanitizeHTML(msg.data)}</span>
         <div class="msg-details">
        <span class="msg-sender">${msg.sender === currentUser ? currentUser : msg.sender}</span>
          <span class="msg-time">${msg.time}</span>
          </div>
          </div>
      </div>
    `
      )
      .join("");

    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = messagesHTML;
    
    // Store scroll position to maintain view
    const currentScrollHeight = messagesContainer.scrollHeight;
    messagesContainer.prepend(...tempDiv.childNodes);
    
    // Adjust scroll to maintain position
    const newScrollHeight = messagesContainer.scrollHeight;
    messagesContainer.scrollTop = newScrollHeight - currentScrollHeight;
  }
}

function sanitizeHTML(str) {
  if (typeof str !== 'string') return '';
  const temp = document.createElement('div');
  temp.textContent = str; 
  return temp.innerHTML;  
}
// Initialize WebSocket connection
function initWebSocket() {
  if (ws && (ws.readyState === WebSocket.CONNECTING || ws.readyState === WebSocket.OPEN)) {
    ("WebSocket already connected or connecting.");
    return; 
  }
  ws = new WebSocket("api/protected/api/messaging"); 

  ws.onopen = () => {
    ("WebSocket connected");
    ws.send(JSON.stringify({ type: "auth", user: currentUser }));
  };

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      ("Received message:", data);

      if (data.type === "message" && data.message) {
        const message = {
          sender: data.message.sender || "",
          data: data.message.data || "",
          time: data.message.time || new Date().toLocaleTimeString(),
          receiver: data.message.receiver || 0,
        };

        if (
          (message.sender === currentUser &&
            message.receiver === currentReceiverId) ||
          (message.sender === currentPartner &&
            message.receiver === currentUserId)
        ) {
          messages.push(message);
          messageCache[currentReceiverId].messages.push(message);

          const currentChat = document.querySelector(".chat-section h2");
          if (
            currentChat &&
            currentPartner &&
            (message.sender === currentPartner ||
              message.receiver === currentPartner)
          ) {
            appendNewMessage(message)
          }else{

            checkUnreadMessages();
          }
        }
      }
    } catch (error) {
      console.error("Error processing message:", error, event.data);
    }
  };

  ws.onclose = () => {
    ("WebSocket disconnected, retrying...");
    setTimeout(initWebSocket, 2000); 
  };
}

function appendNewMessage(message) {
  const messagesContainer = document.querySelector(".chat-messages");
  if (!messagesContainer) return;
  
  // Create message element
  const messageDiv = document.createElement('div');
  messageDiv.className = `chat-msg ${message.sender === currentUser ? "sent" : "received"}`;
  
  messageDiv.innerHTML = `
    <div class="msg-content">
                        <span class="msg-text">${sanitizeHTML(message.data)}</span>
                        <div class="msg-details">
                        <span class="msg-sender">${message.sender === currentUser ? currentUser : message.sender}</span>
                            <span class="msg-time">${message.time}</span>
                        </div>
    </div>
  `;
  
  // Append to container
  messagesContainer.appendChild(messageDiv);
  
  // Scroll to bottom to show new message
  messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

// Attach event listener for inbox button
if (inboxBtn && mainContent) {
  inboxBtn.addEventListener("click", () => renderInbox()
);

}

// Expose function globally for status.js and conversation items
window.openInboxWithUser = function (nickname, receiverId) {
  initWebSocket();
  renderChat(nickname, receiverId);
};

