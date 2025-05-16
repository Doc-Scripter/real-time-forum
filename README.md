# Real-Time-Forum

## Project Overview

A modern, single-page forum application featuring real-time interactions, private messaging, and dynamic content management. Built with Golang, WebSockets, and SQLite for efficient performance and seamless user experienc, and is containerized using Docker.

## Features

### User Authentication

- Secure registration (nickname, email, password, profile data)

- Flexible login (email or nickname  combined with the password)
- Session management with logout capability

### Posts & Comments

- Only registered users can create posts and comments.

- Posts are categorized.

- Visible to all users (registered or not).

### Likes & Dislikes

- Only registered users can like/dislike posts and comments.

- Like/dislike counts are visible to all users.

#### Filtering

- Filter posts by categories.

- Registered users can filter by created posts and liked posts.

### Database (SQLite)

- Stores users, posts, comments, likes/dislikes, and categories.

- Implements SELECT, CREATE, and INSERT queries.

## Setup & Installation

### Prerequisites

Docker installed on your system.

## Run the Project

1. Clone the repository:
   ```
   git clone https://learn.zone01kisumu.ke/git/cliffootieno/real-time-forum
   ```
2. Navigate to the Project Directory:
   ```
   cd real-time-forum
   ```
3. Build and run the Docker container:
   ```
   ./launch_docker_forum.sh
   ````
4. Access the forum at:
   ``` 
   http://localhost:9111

   ````

## Technologies Used

Go (Golang) – Backend API

SQLite – Database

HTML, CSS, JavaScript – Frontend

Docker – Containerization

## Contribution

We welcome contributions to project! To contribute, follow these steps:
1. Fork the repository.
2. Create a new branch (git checkout -b feature/your-feature-name).
3. Make your changes and commit them (git commit -m 'Add some feature').
4. Push to the branch (git push origin feature/your-feature-name).
5. Open a pull request.

## Authors

[STEPHEN OGINGA](https://learn.zone01kisumu.ke/git/steodhiambo)

[CLIFFORD OTIENO](https://learn.zone01kisumu.ke/git/cliffootieno)

## License

This project is licensed under the [MIT](https://opensource.org/license/mit) License.
