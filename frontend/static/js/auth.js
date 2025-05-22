// Authentication related functions
let isLoginMode = true;

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
        showAuthSuccess('Login successful');
        startAuthStatusCheck(); // Start the 30-second interval check
        setTimeout(() => {
            window.location.reload();
        }, 1000);
        // Login validation
        const loginIdentifier = document.getElementById('loginIdentifier').value.trim();
        const password = document.getElementById('password').value.trim();
        
        if (!loginIdentifier || !password) {
            messageDiv.textContent = 'Please fill in all fields';
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

        // Age validation
        if (age < 13) {
            messageDiv.textContent = 'You must be at least 13 years old to register';
            return false;
        }

        // Email validation
        if (!email.includes('@') || !email.includes('.')) {
            messageDiv.textContent = 'Please enter a valid email address';
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
}

function closeAuthModal() {
    document.getElementById('authModal').classList.remove('active');
    document.getElementById('authForm').reset();
}

function toggleAuthMode() {
    isLoginMode = !isLoginMode;
    openAuthModal(isLoginMode ? 'login' : 'register');
}

// periodically check auth status
function startAuthStatusCheck() {
    setInterval(async () => {
        try {
            const response = await fetch('/api/protected/api/auth/status');
            if (!response.ok) {
                // If logged out
                handleError('Your session has ended. Please login again.');
                window.location.href = '/login';
            }
        } catch (error) {
            console.error('Auth check failed:', error);
        }
    }, 30000); // Check every 30 seconds
}

// Called after logging in successfully
async function handleAuth(event) {
    event.preventDefault();

    if (!validateForm()) {
        return;
    }

    const password = document.getElementById('password').value.trim();
    let requestBody;

    if (isLoginMode) {
        const loginIdentifier = document.getElementById('loginIdentifier').value.trim();
        requestBody = {
            loginIdentifier,
            password
        };
    } else {
        requestBody = {
            nickname: document.getElementById('nickname').value.trim(),
            firstName: document.getElementById('firstName').value.trim(),
            lastName: document.getElementById('lastName').value.trim(),
            email: document.getElementById('email').value.trim(),
            age: parseInt(document.getElementById('age').value),
            gender: document.getElementById('gender').value,
            password: password
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

        if (!data.success) {
            throw new Error(data.message || 'Authentication failed');
        }

        if (isLoginMode) {
            // Login success
            showAuthSuccess('Login successful');
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
    } catch (error) {
        console.error('Auth error:', error);
        // Show more detailed error message
        if (error.message) {
            showAuthError(error.message);
        } else if (error.response) {
            showAuthError(error.response.statusText || 'An error occurred');
        } else {
            showAuthError('An error occurred. Please check the console for details.');
        }
    }
}

async function logout() {
    try {
        const response = await fetch('/api/protected/api/logout', { 
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
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