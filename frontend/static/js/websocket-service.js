class WebSocketService {
    constructor() {
        this.connected = false;
        this.connectAttempts = 0;
        this.maxReconnectDelay = 5000;
        this.handlers = {
            post: new Set(),
            comment: new Set(),
            like: new Set()
        };
        this.messageQueue = [];
        this.connect();
    }

    connect() {
        try {
            const userId = localStorage.getItem('userId');
            if (!userId) {
                console.warn('No userId found, delaying WebSocket connection');
                setTimeout(() => this.connect(), 1000);
                return;
            }

            this.ws = new WebSocket(`ws://localhost:9111/ws?userId=${userId}`);
            
            this.ws.onopen = () => {
                ('WebSocket connection established');
                this.connected = true;
                this.connectAttempts = 0;
                
                // Process any queued messages
                while (this.messageQueue.length > 0) {
                    const msg = this.messageQueue.shift();
                    this.send(msg.type, msg.action, msg.data);
                }
            };
            
            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    ('WebSocket message received:', message);
                    this.handleWebSocketMessage(message);
                } catch (error) {
                    console.error('Error processing WebSocket message:', error);
                }
            };

            this.ws.onclose = () => {
                ('WebSocket connection closed');
                this.connected = false;
                this.reconnect();
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.connected = false;
            };
        } catch (error) {
            console.error('Error establishing WebSocket connection:', error);
            this.reconnect();
        }
    }

    reconnect() {
        if (!this.connected) {
            const delay = Math.min(1000 * Math.pow(2, this.connectAttempts), this.maxReconnectDelay);
            (`Attempting to reconnect in ${delay}ms...`);
            setTimeout(() => {
                this.connectAttempts++;
                this.connect();
            }, delay);
        }
    }

    handleWebSocketMessage(message) {
        try {
            ('Processing message:', message);
            const handlers = this.handlers[message.type];
            if (handlers && handlers.size > 0) {
                (`Found ${handlers.size} handlers for type ${message.type}`);
                handlers.forEach(handler => {
                    try {
                        handler(message);
                    } catch (error) {
                        console.error('Error in message handler:', error);
                    }
                });
            } else {
                ('No handlers found for message type:', message.type);
            }
        } catch (error) {
            console.error('Error handling WebSocket message:', error);
        }
    }

    subscribe(type, handler) {
        (`Subscribing to ${type} events`);
        if (this.handlers[type]) {
            this.handlers[type].add(handler);
            (`Handlers for ${type}: ${this.handlers[type].size}`);
        }
    }

    unsubscribe(type, handler) {
        if (this.handlers[type]) {
            this.handlers[type].delete(handler);
        }
    }

    send(type, action, data) {
        const message = {
            type,
            action,
            data,
            timestamp: new Date()
        };

        if (!this.connected) {
            ('WebSocket not connected, queueing message:', message);
            this.messageQueue.push(message);
            return;
        }

        try {
            if (this.ws.readyState === WebSocket.OPEN) {
                ('Sending WebSocket message:', message);
                this.ws.send(JSON.stringify(message));
            } else {
                ('WebSocket not ready, queueing message:', message);
                this.messageQueue.push(message);
            }
        } catch (error) {
            console.error('Error sending WebSocket message:', error);
            this.messageQueue.push(message);
        }
    }
}

// Create global instance
window.wsService = new WebSocketService();