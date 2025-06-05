// Enhanced input validation and sanitization functions
let isLoginMode = true;

// Input sanitization function

// Strict alphanumeric validation for nicknames
function validateNickname(nickname) {
  // Only allow letters, numbers, and underscores
  const nicknameRegex = /^[0-9a-zA-Z_<>!@#$%^&*()`~]+$/;
  return (
    nicknameRegex.test(nickname) &&
    nickname.length >= 3 &&
    nickname.length <= 32
  );
}

// Name validation (only letters and spaces)
function validateName(name) {
  const nameRegex = /^[a-zA-Z_<>!@#$%^&*()`~]+$/;
  return nameRegex.test(name) && name.length >= 3 && name.length <= 32;
}

// Enhanced email validation
function validateEmail(email) {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(email) && email.length <= 254;
}

// Real-time input sanitization and validation feedback
function setupInputSanitization() {
  // Get all input fields with their validation rules
  const inputConfigs = {
    nickname: {
      sanitize: true,
      validate: validateNickname,
      errorMsg:
        "Nickname can only contain letters, numbers, and underscores (2 - 32 characters)",
    },
    firstName: {
      sanitize: true,
      validate: validateName,
      errorMsg: "First name can only contain letters (3- 32 characters)",
    },
    lastName: {
      sanitize: true,
      validate: validateName,
      errorMsg: "Last name can only contain letters (3-32 characters)",
    },
    email: {
      sanitize: true,
      validate: validateEmail,
      errorMsg: "Please enter a valid email address",
    },
    loginIdentifier: {
      sanitize: true,
      validate: validateLoginIdentifier,
      errorMsg: "Please enter a valid email or nickname",
    },
  };

  Object.keys(inputConfigs).forEach((inputId) => {
    const input = document.getElementById(inputId);
    const config = inputConfigs[inputId];

    if (input) {
      // Create or get error message element
      let errorElement = document.getElementById(`${inputId}-error`);
      if (!errorElement) {
        errorElement = document.createElement("div");
        errorElement.id = `${inputId}-error`;
        errorElement.className = "field-error";
        errorElement.style.color = "#ff4d4d";
        errorElement.style.fontSize = "12px";
        errorElement.style.marginTop = "4px";
        errorElement.style.display = "none";
        input.parentNode.appendChild(errorElement);
      }

      // Clear validation styling on focus
      input.addEventListener("focus", function (e) {
        hideFieldError(inputId);
      });

      // Final validation on blur
      input.addEventListener("blur", function (e) {
        const value = e.target.value.trim();
        if (value.length > 0 && !config.validate(value)) {
          showFieldError(inputId, config.errorMsg);
          e.target.classList.add("error");
        }
      });

      // Prevent paste of dangerous content
      input.addEventListener("paste", function (e) {
        setTimeout(() => {
          if (config.sanitize) {
            const cursorPos = e.target.selectionStart;
            const originalLength = e.target.value.length;
            const sanitized = e.target.value;

            if (e.target.value !== sanitized) {
              e.target.value = sanitized;
              const lengthDiff = originalLength - sanitized.length;
              const newCursorPos = Math.max(0, cursorPos - lengthDiff);
              e.target.setSelectionRange(newCursorPos, newCursorPos);
            }
          }

          // Trigger validation after paste
          input.dispatchEvent(new Event("input"));
        }, 0);
      });
    }
  });

  

    

}

function sanitizeInput(input) {
  const sanitized = input.replace(/[<>!@#$%^&*()`~]/g, "");
  return sanitized;
}

// Helper functions for field-specific error messages
function showFieldError(fieldId, message) {
  const errorElement = document.getElementById(`${fieldId}-error`);
  if (errorElement) {
    errorElement.textContent = message;
    errorElement.style.display = "block";
  }
}

function hideFieldError(fieldId) {
  const errorElement = document.getElementById(`${fieldId}-error`);
  if (errorElement) {
    errorElement.style.display = "none";
  }
}

// Validation function for login identifier (email or nickname)
function validateLoginIdentifier(identifier) {
  if (identifier.includes("@")) {
    return validateEmail(identifier);
  } else {
    return validateNickname(identifier);
  }
}

// Enhanced form validation - only validates on form submission
function validateForm() {
  const messageDiv = document.getElementById("authMessage");
  let hasErrors = false;
  let firstErrorField = null;

  // Clear previous field errors
  document
    .querySelectorAll(".field-error")
    .forEach((error) => (error.style.display = "none"));
  document
    .querySelectorAll("input")
    .forEach((input) => input.classList.remove("error"));

  if (isLoginMode) {
    // Login validation
    const loginIdentifier = document.getElementById("loginIdentifier").value;

    const password = document.getElementById("password").value;

    if (!loginIdentifier) {
      showFieldError("loginIdentifier", "Please enter your email or nickname");
      document.getElementById("loginIdentifier").classList.add("error");
      hasErrors = true;
      if (!firstErrorField) firstErrorField = "loginIdentifier";
    } else if (!validateLoginIdentifier(loginIdentifier)) {
      showFieldError(
        "loginIdentifier",
        "Please enter a valid email or nickname"
      );
      document.getElementById("loginIdentifier").classList.add("error");
      hasErrors = true;
      if (!firstErrorField) firstErrorField = "loginIdentifier";
    }

    if (!password) {
      showFieldError("password", "Please enter your password");
      document.getElementById("password").classList.add("error");
      hasErrors = true;
      if (!firstErrorField) firstErrorField = "password";
    }
  } else {
    // Registration validation
    const fields = {
      nickname: {
        value: document.getElementById("nickname").value,
        validator: validateNickname,
        message:
          "Nickname can only contain letters, numbers, and underscores (3-20 characters)",
      },
      firstName: {
        value: document.getElementById("firstName").value,
        validator: validateName,
        message: "First name can only contain letters (2-32 characters)",
      },
      lastName: {
        value: document.getElementById("lastName").value,
        validator: validateName,
        message:
          "Last name can only contain letters and spaces (2-50 characters)",
      },
      email: {
        value: document.getElementById("email").value,
        validator: validateEmail,
        message: "Please enter a valid email address",
      },
      age: {
        value: parseInt(document.getElementById("age").value),
        validator: validateAge,
        message: "Age must be between 13 and 200 years",
      },
      gender: {
        value: document.getElementById("gender").value,
        validator: (gender) =>
          ["male", "female", "other", "prefer-not-to-say"].includes(
            gender.toLowerCase()
          ),
        message: "Please select a gender",
      },
      password: {
        value: document.getElementById("password").value,
        validator: validatePassword,
        message: "Password must meet all requirements",
      },
    };

    Object.keys(fields).forEach((fieldId) => {
      const field = fields[fieldId];
      const inputElement = document.getElementById(fieldId);

      const fieldValue = typeof field.value === 'string' ? field.value.trim() : field.value;


      if (!fieldValue || fieldValue === "") {
        showFieldError(
          fieldId,
          `Please enter ${
            fieldId === "firstName"
              ? "first name"
              : fieldId === "lastName"
              ? "last name"
              : fieldId
          }`
        );
        inputElement.classList.add("error");
        hasErrors = true;
        if (!firstErrorField) firstErrorField = fieldId;
      } else if (!field.validator(field.value)) {
        showFieldError(fieldId, field.message);
        inputElement.classList.add("error");
        hasErrors = true;
        if (!firstErrorField) firstErrorField = fieldId;
      }
    });
  }

  if (hasErrors) {
    // Focus on the first field with an error
    if (firstErrorField) {
      const firstErrorElement = document.getElementById(firstErrorField);
      firstErrorElement.focus();
    }

    messageDiv.textContent = "Please correct the highlighted fields";
    messageDiv.style.display = "block";
    messageDiv.className = "auth-message error";
    return false;
  }

  messageDiv.style.display = "none";
  return true;
}

// Separate password validation function
function validatePassword(password) {
  if (password.length < 8) return false;
  if (!/\d/.test(password)) return false;
  if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) return false;
  if (!/[A-Z]/.test(password)) return false;

  // Check for common weak passwords
  const commonPasswords = ["password", "123456789", "qwerty123", "password123"];
  if (commonPasswords.some((common) => password.toLowerCase().includes(common)))
    return false;

  return true;
}

function validateAge(age){
  if (isNaN(age)||age < 13 || age > 200) {
    return false
  } 
  return true
}

let Initialized = false;
// Authentication check function
async function checkAuth() {
  try {
    const response = await fetch("/api/protected/api/auth/status");
    const data = await response.json();
    const authButtons = document.getElementById("auth-buttons");
    const userFilters = document.getElementById("userFilters");

    if (data.authenticated) {
      if (!Initialized) {
        Initialized = true;
        fetchAndDisplayOnlineUsers();
        setInterval(fetchAndDisplayOnlineUsers, 5000);
        resetPosts();
        setupInfiniteScroll();
        loadFilterCategories();
        await initInbox();
        startUnreadCheck();
      }
      // Sanitize the nickname before displaying
      const sanitizedNickname = data.nickname || "User";
      authButtons.innerHTML = `
                <div class="auth-status">
                    <span class="welcome-text">Welcome, ${sanitizedNickname}</span>
                    <button onclick="logout()" class="auth-btn logout-btn">Logout</button>
                </div>
            `;
      userFilters.style.display = "flex";
    } else {
      clearInterval(fetchAndDisplayOnlineUsers);
      authButtons.innerHTML = `
                <button onclick="openAuthModal('login')" class="auth-btn">Login</button>
                <button onclick="openAuthModal('register')" class="auth-btn">Register</button>
            `;
      userFilters.style.display = "none";
      openAuthModal("login");
    }
  } catch (error) {
    const authButtons = document.getElementById("auth-buttons");
    authButtons.innerHTML = `
            <button onclick="openAuthModal('login')" class="auth-btn">Login</button>
            <button onclick="openAuthModal('register')" class="auth-btn">Register</button>
        `;
    document.getElementById("userFilters").style.display = "none";
  }
}
let authInterval;
function startAuthCheck() {
  authInterval = setInterval(checkAuth, 5000);
}
function stopAuthCheck() {
  if (authInterval) {
    clearInterval(authInterval);
  }
}

// Enhanced authentication handler
async function handleAuth(event) {
  event.preventDefault();
  event.stopPropagation();

  if (!validateForm()) {
    return;
  }

  const messageDiv = document.getElementById("authMessage");
  messageDiv.style.display = "none";

  const password = document.getElementById("password").value;
  let requestBody;

  if (isLoginMode) {
    const loginIdentifier = document.getElementById("loginIdentifier").value;

    requestBody = {
      loginIdentifier,
      password,
    };
  } else {
    // Sanitize all inputs before sending
    requestBody = {
      nickname: sanitizeInput(document.getElementById("nickname").value),
      firstName: sanitizeInput(document.getElementById("firstName").value),
      lastName: sanitizeInput(document.getElementById("lastName").value),
      email: document.getElementById("email").value,
      age: parseInt(document.getElementById("age").value),
      gender: sanitizeInput(document.getElementById("gender").value),
      password,
    };
  }

  try {
    const response = await fetch(`/api/${isLoginMode ? "login" : "register"}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(requestBody),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      if (isLoginMode) {
        showAuthSuccess("Login successful");
        startAuthStatusCheck();
        setTimeout(() => {
          window.location.reload();
        }, 1000);
      } else {
        showAuthSuccess("Registration successful! Please login.");
        setTimeout(() => {
          isLoginMode = true;
          openAuthModal("login", "Registration successful! Please login.");
        }, 1000);
      }
    } else {
      const errorMessage = data.message || "An error occurred";
      showAuthError(errorMessage);
    }
  } catch (error) {
    console.error("Auth error:", error);
    showAuthError("An error occurred. Please try again later.");
  }
}

// Initialize input sanitization when the page loads
document.addEventListener("DOMContentLoaded", function () {
  startAuthCheck();
  setupInputSanitization();
});

function showAuthError(message) {
  const messageDiv = document.getElementById("authMessage");
  messageDiv.textContent = message;
  messageDiv.style.display = "block";
  messageDiv.className = "auth-message error";
}

function showAuthSuccess(message) {
  const messageDiv = document.getElementById("authMessage");
  messageDiv.textContent = message;
  messageDiv.style.display = "block";
  messageDiv.className = "auth-message success";
}

function openAuthModal(mode, message = "") {
  isLoginMode = mode === "login";
  const modalTitle = document.getElementById("modalTitle");
  const messageDiv = document.getElementById("authMessage");
  const strengthContainer = document.getElementById(
    "password-strength-container"
  );
  const requirementsDiv = document.getElementById("password-requirements");
  const loginFields = document.getElementById("loginFields");
  const registerFields = document.getElementById("registerFields");

  modalTitle.textContent = isLoginMode ? "Login" : "Register";

  if (message) {
    showAuthSuccess(message);
  } else {
    messageDiv.style.display = "none";
  }

  strengthContainer.style.display = isLoginMode ? "none" : "block";
  requirementsDiv.style.display = isLoginMode ? "none" : "block";
  loginFields.style.display = isLoginMode ? "block" : "none";
  registerFields.style.display = isLoginMode ? "none" : "block";

  document.getElementById("authModal").classList.add("active");
  document.querySelector(".modal-switch").textContent = isLoginMode
    ? "Register Instead"
    : "Login Instead";
  document.querySelector(".modal-submit").textContent = isLoginMode
    ? "Login"
    : "Register";

  if (!isLoginMode) {
    checkPasswordStrength("");
    stopAuthCheck();
  }

  // Clear any existing field errors and styling
  document
    .querySelectorAll(".field-error")
    .forEach((error) => (error.style.display = "none"));
  document.querySelectorAll("input").forEach((input) => {
    input.classList.remove("error", "valid");
  });

  // Setup input sanitization for the modal (call this after the modal is open)
  setTimeout(() => {
    setupInputSanitization();
  }, 100);
}

function closeAuthModal() {
  document.getElementById("authModal").classList.remove("active");
  document.getElementById("authForm").reset();
}

function toggleAuthMode() {
  isLoginMode = !isLoginMode;
  openAuthModal(isLoginMode ? "login" : "register");
}

function startAuthStatusCheck() {
  setInterval(async () => {
    try {
      const response = await fetch("/api/protected/api/auth/status");
      if (!response.ok) {
        handleError("Your session has ended. Please login again.");
        window.location.href = "/";
      }
    } catch (error) {
      console.error("Auth check failed:", error);
    }
  }, 30000);
}

async function logout() {
  try {
    const response = await fetch("/api/protected/api/logout", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    if (!response.ok) {
      if (response.status === 405) {
        handleError("Invalid request method");
        return;
      }
      throw new Error("Logout failed");
    }

    window.location.reload();
  } catch (error) {
    console.error("Error logging out:", error);
    handleError("Failed to logout. Please try again.");
  }
}

function checkPasswordStrength(password) {
  const strengthMeter = document.getElementById("password-strength");
  const strengthText = document.getElementById("strength-text");

  const lengthCheck = document.getElementById("length-check");
  const numberCheck = document.getElementById("number-check");
  const specialCheck = document.getElementById("special-check");
  const uppercaseCheck = document.getElementById("uppercase-check");

  if (!password) {
    strengthMeter.style.width = "0%";
    strengthMeter.style.backgroundColor = "#e0e0e0";
    strengthText.textContent = "";

    [lengthCheck, numberCheck, specialCheck, uppercaseCheck].forEach(
      (check) => {
        if (check) check.classList.remove("valid");
      }
    );
    return;
  }

  let strength = 0;

  if (password.length >= 8) {
    strength += 25;
    if (lengthCheck) lengthCheck.classList.add("valid");
  } else {
    if (lengthCheck) lengthCheck.classList.remove("valid");
  }

  if (/\d/.test(password)) {
    strength += 25;
    if (numberCheck) numberCheck.classList.add("valid");
  } else {
    if (numberCheck) numberCheck.classList.remove("valid");
  }

  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
    strength += 25;
    if (specialCheck) specialCheck.classList.add("valid");
  } else {
    if (specialCheck) specialCheck.classList.remove("valid");
  }

  if (/[A-Z]/.test(password)) {
    strength += 25;
    if (uppercaseCheck) uppercaseCheck.classList.add("valid");
  } else {
    if (uppercaseCheck) uppercaseCheck.classList.remove("valid");
  }

  strengthMeter.style.width = `${strength}%`;

  if (strength <= 25) {
    strengthMeter.style.backgroundColor = "#ff4d4d";
    strengthText.textContent = "Weak";
  } else if (strength <= 50) {
    strengthMeter.style.backgroundColor = "#ffd700";
    strengthText.textContent = "Fair";
  } else if (strength <= 75) {
    strengthMeter.style.backgroundColor = "#2ecc71";
    strengthText.textContent = "Good";
  } else {
    strengthMeter.style.backgroundColor = "#27ae60";
    strengthText.textContent = "Strong";
  }
}
