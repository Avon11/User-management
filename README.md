# User-management

This project is a task management API built with Go and Fiber, featuring user authentication and task management capabilities. The API uses MongoDB as the database, and the Go Fiber web framework is used to handle the routes. It also incorporates middleware for protected routes that require user authentication.

## Features

- **User authentication**: Sign up, sign in, and sign out.
- **Task management**: Create, read, update, and delete tasks.
- **Protected routes**: Task management routes require authentication via middleware.

## Project Structure

```
task-manager-api/
├── main.go
├── go.mod
├── go.sum
├── config/
│ └── config.go
├── models/
│ ├── user.go
│ └── task.go
├── handlers/
│ ├── userHandler.go
│ └── taskHandler.go
├── middleware/
│ └── middleware.go
├── database/
│ └── mongodb.go
├── utils/
│ └── jwt.go
└── README.md
```

## Environment Variables

Create a `.env` file in the project root with the following content:

```bash
MONGO_URI=mongodb://localhost:27017/kenshi
PORT=8080
JWT_SECRET=your_jwt_secret
```

## Installation and Setup

```
git clone https://github.com/Avon11/User-management.git
go mod tidy
go run main.go
```

The server will start on http://localhost:8080

## API Endpoints

### POST /api/signup – Create a new user account.

```request
{
    "username": "Avon11",
    "email": "avon11@email.com",
    "password":"987654321"
}
```

```response
{
    "id": "6705570247768420e03d865c",
    "username": "Avon11",
    "email": "avon11@email.com",
    "password": ""
}
```

### POST /api/signin – Log in an existing user.

```request
{
    "email": "avon11@email.com",
    "password": "987654321"
}
```

```response
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg0MDY4MzMsInVzZXJfaWQiOiI2NzA1NTcwMjQ3NzY4NDIwZTAzZDg2NWMifQ.Dp-4zjdaFfJRR0T8kGE8G_E1p9DjFt_WrEBdWzyPOWc"
}
```

### POST /api/signout – Log out an existing user.

```response
OK
```

### POST /api/tasks/ – Create a new task.

```request
{
    "title": "Abhinav5 task 2",
    "description": "this is task 2 for abhinav 5",
    "status": "active"
}
```

```response
{
    "id": "67055651a163999f34b41efd",
    "userId": "670407ac3a334c01fe44ea69",
    "title": "Abhinav5 task 2",
    "description": "this is task 2 for abhinav 5",
    "status": "active",
    "createdAt": "2024-10-08T21:27:05.8226426+05:30",
    "updatedAt": "2024-10-08T21:27:05.8226426+05:30"
}
```

### GET /api/tasks/ – Get all tasks for the authenticated user.

```response
[
    {
        "id": "6704091123dffaf480ad48cb",
        "userId": "670407ac3a334c01fe44ea69",
        "title": "Abhinav5 task 1",
        "description": "this is task 1 for abhinav 5",
        "status": "active",
        "createdAt": "2024-10-07T16:15:13.511Z",
        "updatedAt": "2024-10-07T16:15:13.511Z"
    },
    {
        "id": "670409443a766927b969f004",
        "userId": "670407ac3a334c01fe44ea69",
        "title": "Abhinav5 task 2",
        "description": "this is task 2 for abhinav 5",
        "status": "active",
        "createdAt": "2024-10-07T16:16:04.695Z",
        "updatedAt": "2024-10-07T16:16:04.695Z"
    }
]
```

### GET /api/tasks/:id – Get a specific task by ID.

```response
{
        "id": "670409443a766927b969f004",
        "userId": "670407ac3a334c01fe44ea69",
        "title": "Abhinav5 task 2",
        "description": "this is task 2 for abhinav 5",
        "status": "active",
        "createdAt": "2024-10-07T16:16:04.695Z",
        "updatedAt": "2024-10-07T16:16:04.695Z"
}
```

### PUT /api/tasks/:id – Update a task by ID

```request
{
    "title": "Abhinav5 3",
    "description": "this is 3",
    "status": "active"
}
```

```response
OK
```

### DELETE /api/tasks/:id – Delete a task by ID

```response
OK
```

### GET /api/token – Get new access token

```response
{
    "token":"tertyuirdtfyguh2345678"
}
```
