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
let isInitialized = true;
let messagesPerPage = 10; // Number of messages to display per page
let currPage = 1; // Current page number
let ws;
let endIndex = 0;
async function checkUnreadMessages() {
  console.log("DEBUG: Checking unread messages");
  try {
    const response = await fetch("/api/protected/api/unread");
    if (response.ok) {
      const count = await response.json();
      console.log("DEBUG: Unread count:", count);
      const badge = document.getElementById("unread-badge");

      if (count > 0) {
        badge.style.display = "block";
        badge.textContent = count;
      } else {
        badge.style.display = "none";
      }
    }
  } catch (error) {
    console.error("DEBUG: Error checking unread messages:", error);
  }
}

// Start checking for unread messages
function startUnreadCheck() {
  checkUnreadMessages(); // Initial check
  unreadCheckInterval = setInterval(checkUnreadMessages, 10000); // Check every 10 seconds
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

  console.log("DEBUG - Last Messages:", lastMessages);
  console.log("DEBUG - Current User:", currentUser);

  return lastMessages.map((msg) => {
    console.log("DEBUG - Processing message:", msg);

    // Log values before deciding partner
    console.log("DEBUG - Sender:", msg.sender);
    console.log("DEBUG - Receiver:", msg.receiver);
    let receiverName = msg.receiver;
    const onlineUsers = document.querySelectorAll(".online-user");
    onlineUsers.forEach((user) => {
      if (user.dataset.receiverId == msg.receiver) {
        receiverName = user.querySelector(".receiver").textContent;
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
        receiver: receiverName, // Changed from msg.sender to msg.receiver
        receiverId: msg.receiver,
      },
    };
  });
}
// Render inbox: conversation list or chat view
async function renderInbox() {
  document.getElementById("toggleRight").style.display = "flex";
  document.getElementById("toggleLeft").style.display = "flex";
  if (!currentUser) {
    mainContent.innerHTML = "<div>Please log in to view your inbox.</div>";
    return;
  }

  await fetchLastMessages();
  console.log("DEBUG - After fetchLastMessages:", lastMessages);

  startInboxRefresh();

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
window.addEventListener("beforeunload", () => {
  stopUnreadCheck();
  stopInboxRefresh();
});
// Render chat with selected user
async function renderChat(partner, receiverId) {
  console.log("DEBUG: Opening chat with:", partner);

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
  
  if (isInitialized) {
    messageCache[receiverId] = {
      messages: messages,
      currPage: 1,
    };
    messages = messages.slice(-10);
    
    isInitialized = false;
  }
  console.log("DEBUG (info) :\n ", messages.length);

  const chatMessages = messages.filter((msg) => {
    if (!msg) return false;
    return (
      (msg.sender === currentUser && msg.receiver === receiverId) ||
      (msg.sender === partner && msg.receiver === currentUserId)
    );
  });

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
    // Refresh unread count
    await checkUnreadMessages();
  } catch (error) {
    console.error("Error marking messages as read:", error);
  }


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

  if (isInitialized) {
 

    const messagesContainer = document.querySelector(".chat-messages");
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  } else {
    const messagesContainer = document.querySelector(".chat-messages");
    const newMessageElements =
      messagesContainer.querySelectorAll(".chat-msg.new");
    if (newMessageElements.length > 0) {
      newMessageElements[0].scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
    }
  }


  // if (messagesContainer) {
  //   messagesContainer.scrollTop = messagesContainer.scrollHeight;
  // }

  document.getElementById("back-to-inbox").onclick = () => renderInbox();

  document.getElementById("send-msg-form").onsubmit = async (e) => {
    endIndex + 1;
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
        time: new Date().toLocaleTimeString(),
      })
    );
    messageCache[receiverId].messages.push({
      sender: currentUser,
      receiver: receiverId,
      data: text,
      time: new Date().toLocaleTimeString(),
    });

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
  // Add scroll event listener to load more messages
  const messagesContainer = document.querySelector(".chat-messages");

  messagesContainer.addEventListener("scroll", () => {
    if (messagesContainer.scrollTop === 0) {
      loadMoreMessages(receiverId);
    }
  });
}

function loadMoreMessages(receiverId) {
  const cache = messageCache[receiverId];
  if (cache.currPage * messagesPerPage < cache.messages.length) {
    cache.currPage++;
    const startIndex = (cache.currPage - 1) * messagesPerPage;
    const endIndex = startIndex + messagesPerPage;
    const messagesToPrepend = cache.messages.slice(startIndex, endIndex);

    const messagesContainer = document.querySelector(".chat-messages");
    const messagesHTML = messagesToPrepend
      .map(
        (msg) => `
      <div class="chat-msg ${msg.sender === currentUser ? "sent" : "received"}">
        <div class="msg-content">
          <span class="msg-text">${msg.data}</span>
          <span class="msg-time">${msg.time}</span>
        </div>
      </div>
    `
      )
      .join("");

    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = messagesHTML;
    messagesContainer.prepend(...tempDiv.childNodes);
  }
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
          messageCache[currentReceiverId].messages.push(message);
          endIndex + 1;

          // Re-render if viewing this conversation
          const currentChat = document.querySelector(".chat-section h2");
          if (
            currentChat &&
            currentPartner &&
            (message.sender === currentPartner ||
              message.receiver === currentPartner)
          ) {
            // Re-render the chat with updated messages
            renderChat(currentPartner, currentReceiverId);
          }
        }
        checkUnreadMessages();
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
  startUnreadCheck();
});
