# A Golang Todo REST API

A full-stack REST API built with the Gin framework. Features PostgreSQL, Redis, JWT Authentication. 

## Features
___
User Authentication: Secure registration and login with JWT tokens aswell as refresh token for session persistance
Password Security: bcrypt hashing for secure password storage
Redis: access and refresh token storage
Protected Routes: Middleware-based route protection
User-Specific Todos: Each user has their own private todo collection
CRUD Operations: Create, read, update, and delete todos
Database Migrations: Version-controlled database schema changes
Hot Reloading: Air integration for development
