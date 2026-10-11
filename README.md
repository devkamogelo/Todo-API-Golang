# A Golang Todo REST API

A full-stack REST API built with the Gin framework. Features PostgreSQL, Redis, JWT Authentication. 

## Features
- User Authentication: Secure registration and login with JWT tokens aswell as refresh token for session persistance
- Password Security: bcrypt hashing for secure password storage
- Redis: access and refresh token storage
- Protected Routes: Middleware-based route protection
- User-Specific Todos: Each user has their own private todo collection
- CRUD Operations: Create, read, update, and delete todos
- Database Migrations: Version-controlled database schema changes
- Hot Reloading: Air integration for development

## 10 API endpoints
## Todo Management
| Method | Endpoint     |       Description       |
|--------|--------------|-------------------------|
|  GET   |  /todos      | getting all users todos |
|  GET   |  /todos/:id  | getting a todo by id    |
|  POST  |  /todos      | creating todos          |
|  PUT   |  /todos/:id  | updating a todo         |
| Delete | /todos/:id   | delete a todo           |

## User Management
| Method | Endpoint     |       Description       |
|--------|--------------|-------------------------|
|  POST  |  /register   | create a user           |
|  POST  |  /login      | generate access token   |
|  POST  |   /refresh   | refresh access token    |

## Test/Heath Check
| Method | Endpoint           |       Description       |
|--------|--------------------|-------------------------|
|   GET  |  /                 |       Health Check      |
|   GET  |  /protected-route  |   Test Auth-Middleware  |
