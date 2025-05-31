function showErrorPage(errorData) {
  // const error_page = document.querySelector("body");
  document.open();
  document.writeln(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${errorData.Title} - Forum</title>
    <title>${errorData.Title} - Forum</title>
    <link rel="stylesheet" href="/static/css/styles.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
    <style>
        main {
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: calc(100vh - 70px);
            padding: 1rem;
        }

        .error-container {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            text-align: center;
            padding: 2rem;
            background-color: var(--bg-card);
            border-radius: var(--border-radius);
            width: 100%;
            max-width: 600px;
        }

        .error-code {
            font-size: 6rem;
            font-weight: bold;
            color: var(--primary-color);
            margin-bottom: 1rem;
            line-height: 1;
        }

        .error-message {
            font-size: 1.5rem;
            color: var(--text-primary);
            margin-bottom: 2rem;
        }

        .error-description {
            color: var(--text-secondary);
            margin-bottom: 2rem;
            line-height: 1.5;
        }

        .home-button {
            background-color: var(--primary-color);
            color: white;
            padding: 0.75rem 1.5rem;
            border-radius: var(--border-radius);
            text-decoration: none;
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            transition: var(--transition);
        }

        .home-button:hover {
            background-color: var(--primary-hover);
            transform: translateY(-2px);
        }

        .error-icon {
            font-size: 4rem;
            color: var(--primary-color);
            margin-bottom: 1rem;
        }
    </style>
</head>
<body>
    <header>
        <nav>
            <div class="nav-left">
                <a href="/" class="logo">
                    <i class="fas fa-comments"></i>
                    Forum
                </a>
            </div>
        </nav>
    </header>

    <main>
        <div class="error-container">
            <i class="fas fa-exclamation-circle error-icon"></i>
            <div class="error-code">${errorData.Code}</div>
            <h1 class="error-message">${errorData.Title}</h1>
            <p class="error-description">${errorData.Description}</p>
            <div class="error-code">${errorData.Code}</div>
            <h1 class="error-message">${errorData.Title}</h1>
            <p class="error-description">${errorData.Description}</p>
            <a href="/" class="home-button">
                <i class="fas fa-home"></i>
                Return to Home
            </a>
        </div>
    </main>
</body>
</html> `);
  document.close();
}

// Intercept all fetch requests
if (!window.originalFetch) {
    const originalFetch = window.fetch;
    window.originalFetch = originalFetch; 
  }
window.fetch = async function (...args) {
  const errorData = {
    400: {
      Code: 400,
      Title: "Bad Request",
      Description:
        "The server cannot process your request due to invalid syntax.",
    },
    401: {
      Code: 401,
      Title: "Unauthorized",
      Description: "You need to be logged in to access this resource.",
    },
    403: {
      Code: 403,
      Title: "Forbidden",
      Description: "You don't have permission to access this resource.",
    },
    405: {
      Code: 405,
      Title: "Page Not Found",
      Description:
        "The page you're looking for doesn't exist or has been moved.",
    },
    500: {
      Code: 500,
      Title: "Internal Server Error",
      Description: "Something went wrong on our end. Please try again later.",
    },
  };
  try {
    const response = await originalFetch(...args);
    if (response.headers.get("X-Status-Code") === "404") {
      let error = errorData[405];
      showErrorPage(error);
      throw new Error("Page not found");
    }

    if (response.status === 403) {
      let error = errorData[403];
      showErrorPage(error);
      throw new Error("Forbidden");
    }
    if (response.status === 400) {
      let error = errorData[400];
      showErrorPage(error);
      throw new Error("Bad Request");
    }
    ("code status", response.status);
    return response;
  } catch (error) {
    ("error fetching: ", error);
  }
};

if (
  window.location.pathname !== "/" &&
  !window.location.pathname.startsWith("/static/")
) {
  // Check if this is a valid route by making a HEAD request
  fetch(window.location.pathname, { method: "HEAD" }).catch((error) => {
    // The catch will be triggered by our intercepted fetch above if it's a 404
    ("Route not found, error page already shown: ", error);
  });
}
