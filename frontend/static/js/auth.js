// Authentication related functions
let isLoginMode = true;

// Simple input validation helper
function hasInvalidChars(value) {
    return /<[^>]*>/.test(value) || /[<>"'&]/.test(value);
}

function sanitizeInput(value) {
    return value.replace(/<[^>]*>/g, '').replace(/[<>"'&]/g, '').trim();
}

// Special password sanitization that removes HTML tags and special characters
function sanitizePassword(value) {
    // Remove HTML tags
    value = value.replace(/<[^>]*>/g, '');
    // Remove special characters except allowed ones
    value = value.replace(/[<>"'&]/g, '');
    // Remove control characters and whitespace
    value = value.replace(/\s/g, '');
    return value;
}

function hasInvalidPasswordChars(value) {
    // Check for HTML tags
    if (/<[^>]*>/.test(value)) return true;
    // Check for special characters
    if (/[<>"]'&]/.test(value)) return true;
    // Check for control characters
    if (/[\x00-\x1F\x7F-\x9F]/.test(value)) return true;
    return false;
}

async function checkAuth() {
    try {
        const response = await fetch('/api/protected/api/auth/status');
        const data = await response.json();
        const authButtons = document.getElementById('auth-buttons');
        const userFilters = document.getElementById('userFilters');
        
        if (data.authenticated) {
            authButtons.innerHTML = `
                <div class="auth-status">-
                    <span class="welcome-text">Welcome, ${data.nickname}</span>
                    <button onclick="logout()" class="auth-btn logout-btn">Logout</button>
                </div>
            `;
            userFilters.style.display = 'flex';
        } else {
            authButtons.innerHTML = `
                <button onclick="openAuthModal('login')" class="auth-btn">Login</button>
                <button onclick="openAuthModal('register')" class="auth-btn">Register</button>
            `;
            userFilters.style.display = 'none';
            openAuthModal('login');
        }
    } catch (error) {
        console.error('Error checking auth status:', error);
        // Handle error state
        const authButtons = document.getElementById('auth-buttons');
        authButtons.innerHTML = `
            <button onclick="openAuthModal('login')" class="auth-btn">Login</button>
            <button onclick="openAuthModal('register')" class="auth-btn">Register</button>
        `;
        document.getElementById('userFilters').style.display = 'none';
    }
}


function validateForm() {
    const messageDiv = document.getElementById('authMessage');
    messageDiv.style.display = 'block';
    messageDiv.className = 'auth-message error';

    if (isLoginMode) {
        // Login validation
        const loginIdentifier = document.getElementById('loginIdentifier').value.trim();
        const password = document.getElementById('password').value.trim();
        
        if (!loginIdentifier || !password) {
            messageDiv.textContent = 'Please fill in all fields';
            return false;
        }
        
        // Check for HTML/invalid characters in login identifier
        if (hasInvalidChars(loginIdentifier)) {
            messageDiv.textContent = 'Login field contains invalid characters';
            return false;
        }
        
        // Check for HTML/invalid characters in password
        if (hasInvalidPasswordChars(password)) {
            messageDiv.textContent = 'Password contains invalid characters';
            return false;
        }
    } else {
        // Registration validation
        
        const nickname = document.getElementById('nickname').value.trim();
        const firstName = document.getElementById('firstName').value.trim();
        const lastName = document.getElementById('lastName').value.trim();
        const email = document.getElementById('email').value.trim();
        const age = document.getElementById('age').value;
        const gender = document.getElementById('gender').value;
        const password = document.getElementById('password').value.trim();

        if ( !nickname || !firstName || !lastName || !email || !age || !gender || !password) {
            messageDiv.textContent = 'Please fill in all fields';
            return false;
        }

        // Check for HTML/invalid characters in text fields
        if (hasInvalidChars(nickname)) {
            messageDiv.textContent = 'Nickname contains invalid characters';
            return false;
        }
        
        // Check for HTML/invalid characters in password
        if (hasInvalidPasswordChars(password)) {
            messageDiv.textContent = 'Password contains invalid characters';
            return false;
        }
        if (hasInvalidChars(firstName)) {
            messageDiv.textContent = 'First name contains invalid characters';
            return false;
        }
        
        if (hasInvalidChars(lastName)) {
            messageDiv.textContent = 'Last name contains invalid characters';
            return false;
        }
        
        if (hasInvalidChars(email)) {
            messageDiv.textContent = 'Email contains invalid characters';
            return false;
        }

        // Age validation
        if (age < 13 || age > 200) {
            messageDiv.textContent = 'Age must be between 13 and 200 years';
            return false;
        }

        // Email validation
        if (!email.includes('@') || !email.includes('.')) {
            messageDiv.textContent = 'Please enter a valid email address';
            return false;
        }

        // Check for HTML/invalid characters in password
        if (hasInvalidPasswordChars(password)) {
            messageDiv.textContent = 'Password contains invalid characters';
            return false;
        }

        // Password validation
        if (password.length < 8) {
            messageDiv.textContent = 'Password must be at least 8 characters long';
            return false;
        }
        
        if (!/\d/.test(password)) {
            messageDiv.textContent = 'Password must contain at least one number';
            return false;
        }
        
        if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
            messageDiv.textContent = 'Password must contain at least one special character';
            return false;
        }
        
        if (!/[A-Z]/.test(password)) {
            messageDiv.textContent = 'Password must contain at least one uppercase letter';
            return false;
        }
    }

    messageDiv.style.display = 'none';
    return true;
}

function showAuthError(message) {
    const messageDiv = document.getElementById('authMessage');
    messageDiv.textContent = message;
    messageDiv.style.display = 'block';
    messageDiv.className = 'auth-message error';
}

function showAuthSuccess(message) {
    const messageDiv = document.getElementById('authMessage');
    messageDiv.textContent = message;
    messageDiv.style.display = 'block';
    messageDiv.className = 'auth-message success';
}

function openAuthModal(mode, message = '') {
    isLoginMode = mode === 'login';
    const modalTitle = document.getElementById('modalTitle');
    const messageDiv = document.getElementById('authMessage');
    const strengthContainer = document.getElementById('password-strength-container');
    const requirementsDiv = document.getElementById('password-requirements');
    const loginFields = document.getElementById('loginFields');
    const registerFields = document.getElementById('registerFields');
    
    modalTitle.textContent = isLoginMode ? 'Login' : 'Register';
    
    if (message) {
        showAuthSuccess(message);
    } else {
        messageDiv.style.display = 'none';
    }
    
    // Show/hide password strength elements and requirements based on mode
    strengthContainer.style.display = isLoginMode ? 'none' : 'block';
    requirementsDiv.style.display = isLoginMode ? 'none' : 'block';
    loginFields.style.display = isLoginMode ? 'block' : 'none';
    registerFields.style.display = isLoginMode ? 'none' : 'block';
    
    document.getElementById('authModal').classList.add('active');
    document.querySelector('.modal-switch').textContent = isLoginMode ? 'Register Instead' : 'Login Instead';
    document.querySelector('.modal-submit').textContent = isLoginMode ? 'Login' : 'Register';
    
    // Reset password strength meter and requirements
    if (!isLoginMode) {
        checkPasswordStrength('');
    }
    
    // Setup input sanitization when modal opens
    setupInputSanitization();
}

function closeAuthModal() {
    document.getElementById('authModal').classList.remove('active');
    document.getElementById('authForm').reset();
}

function toggleAuthMode() {
    isLoginMode = !isLoginMode;
    openAuthModal(isLoginMode ? 'login' : 'register');
}

// Auto-sanitize inputs as user types and show feedback
function setupInputSanitization() {
    const messageDiv = document.getElementById('authMessage');
    const textInputIds = ['nickname', 'firstName', 'lastName', 'email', 'loginIdentifier'];
    const passwordInputIds = ['password'];
    
    // Function to show temporary message in red
    function showMessage(message, duration = 2000) {
        messageDiv.textContent = message;
        messageDiv.style.display = 'block';
        messageDiv.className = 'auth-message error';
        messageDiv.style.color = 'red';
        
        // Clear any existing timeout
        if (messageDiv.dataset.timeoutId) {
            clearTimeout(parseInt(messageDiv.dataset.timeoutId));
        }
        
        // Set new timeout to hide message
        const timeoutId = setTimeout(() => {
            messageDiv.style.display = 'none';
        }, duration);
        
        messageDiv.dataset.timeoutId = timeoutId;
    }
    
    // Function to handle input sanitization
    function handleInput(input, sanitizeFunction, errorMessage) {
        if (input && !input.hasAttribute('data-sanitized')) {
            input.setAttribute('data-sanitized', 'true');
            input.addEventListener('input', (e) => {
                const original = e.target.value;
                const sanitized = sanitizeFunction(original);
                
                if (original !== sanitized) {
                    e.target.value = sanitized;
                    showMessage(errorMessage);
                }
            });
        }
    }
    
    // Handle text inputs
    textInputIds.forEach(id => {
        const input = document.getElementById(id);
        handleInput(input, sanitizeInput, 'Invalid character not allowed');
    });
    
    // Handle password inputs
    passwordInputIds.forEach(id => {
        const input = document.getElementById(id);
        handleInput(input, sanitizePassword, 'Invalid characters not allowed.');
    });
}

// periodically check auth status
function startAuthStatusCheck() {
    setInterval(async () => {
        try {
            const response = await fetch('/api/protected/api/auth/status');
            if (!response.ok) {
                // If logged out
                handleError('Your session has ended. Please login again.');
                window.location.href = '/';
            }
        } catch (error) {
            console.error('Auth check failed:', error);
        }
    }, 30000); // Check every 30 seconds
}

// Called after logging in successfully
async function handleAuth(event) {
    event.preventDefault();
    event.stopPropagation(); // Prevent event bubbling

    if (!validateForm()) {
        return;
    }

    const messageDiv = document.getElementById('authMessage');
    messageDiv.style.display = 'none'; // Clear any previous messages

    const password = document.getElementById('password').value.trim();
    let requestBody;

    if (isLoginMode) {
        const loginIdentifier = document.getElementById('loginIdentifier').value.trim();
        requestBody = {
            loginIdentifier: sanitizeInput(loginIdentifier),
            password
        };
    } else {
        requestBody = {
            nickname: sanitizeInput(document.getElementById('nickname').value.trim()),
            firstName: sanitizeInput(document.getElementById('firstName').value.trim()),
            lastName: sanitizeInput(document.getElementById('lastName').value.trim()),
            email: sanitizeInput(document.getElementById('email').value.trim()),
            age: parseInt(document.getElementById('age').value),
            gender: document.getElementById('gender').value,
            password: sanitizePassword(password)
        };
    }

    try {
        const response = await fetch(`/api/${isLoginMode ? 'login' : 'register'}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
        });

        const data = await response.json();
        console.log('Auth response:', data);

        if (response.ok && data.success) {
            if (isLoginMode) {
                // Login success
                showAuthSuccess('Login successful');
                startAuthStatusCheck();
                setTimeout(() => {
                    window.location.reload();
                }, 1000);
            } else {
                // Registration success
                showAuthSuccess('Registration successful! Please login.');
                setTimeout(() => {
                    isLoginMode = true;
                    openAuthModal('login', 'Registration successful! Please login.');
                }, 1000);
            }
        } else {
            // Handle login or registration failure
            if (data.message) {
                showAuthError(data.message);
            } else if (isLoginMode) {
                showAuthError('Invalid credentials. Please check your email/nickname and password.');
            } else {
                showAuthError('Registration failed. Please try again.');
            }
        }
    } catch (error) {
        console.error('Auth error:', error);
        showAuthError('An error occurred. Please try again later.');
    }
}

async function logout() {
    try {
        const response = await fetch('/api/protected/api/logout', { 
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        });

        if (!response.ok) {
            if (response.status === 405) {
                handleError("Invalid request method");
                return;
            }
            console.log("Logout failed: ", response.statusText);
            throw new Error('Logout failed');
        }

        window.location.reload();
    } catch (error) {
        console.error('Error logging out:', error);
        handleError('Failed to logout. Please try again.');
    }
}

// Update the checkPasswordStrength function
function checkPasswordStrength(password) {
    const strengthMeter = document.getElementById('password-strength');
    const strengthText = document.getElementById('strength-text');
    
    // Get requirement elements
    const lengthCheck = document.getElementById('length-check');
    const numberCheck = document.getElementById('number-check');
    const specialCheck = document.getElementById('special-check');
    const uppercaseCheck = document.getElementById('uppercase-check');
    
    if (!password) {
        strengthMeter.style.width = '0%';
        strengthMeter.style.backgroundColor = '#e0e0e0';
        strengthText.textContent = '';
        
        // Reset all requirement checks
        [lengthCheck, numberCheck, specialCheck, uppercaseCheck].forEach(check => {
            check.classList.remove('valid');
        });
        return;
    }

    let strength = 0;
    
    // Length check
    if (password.length >= 8) {
        strength += 25;
        lengthCheck.classList.add('valid');
    } else {
        lengthCheck.classList.remove('valid');
    }

    // Contains number
    if (/\d/.test(password)) {
        strength += 25;
        numberCheck.classList.add('valid');
    } else {
        numberCheck.classList.remove('valid');
    }

    // Contains special char
    if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
        strength += 25;
        specialCheck.classList.add('valid');
    } else {
        specialCheck.classList.remove('valid');
    }

    // Contains uppercase
    if (/[A-Z]/.test(password)) {
        strength += 25;
        uppercaseCheck.classList.add('valid');
    } else {
        uppercaseCheck.classList.remove('valid');
    }

    // Update strength meter
    strengthMeter.style.width = `${strength}%`;
    
    // Update color based on strength
    if (strength <= 25) {
        strengthMeter.style.backgroundColor = '#ff4d4d';
        strengthText.textContent = 'Weak';
    } else if (strength <= 50) {
        strengthMeter.style.backgroundColor = '#ffd700';
        strengthText.textContent = 'Fair';
    } else if (strength <= 75) {
        strengthMeter.style.backgroundColor = '#2ecc71';
        strengthText.textContent = 'Good';
    } else {
        strengthMeter.style.backgroundColor = '#27ae60';
        strengthText.textContent = 'Strong';
    }
}