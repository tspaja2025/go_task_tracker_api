# Go Task Tracker API

A simple task tracker API built with Go.

## Roadmap.sh beginner project
This project was created as a part of Todo List API beginner project.
Check out the project details [roadmap.sh](https://roadmap.sh/projects/todo-list-api)

## TODO

* User registration to create a new user
* Login endpoint to authenticate the user and generate a token
* CRUD operations for managing the task list
* User authentication to allow only authorized users to access the task list
* Error handling and security measures
* Database for storing the user and task list
* Data validation
* pagination, filtering and sorting for the task list
* Unit tests for the API
* Rate limiting and throttling for the API
* Refresh token mechanism for the authentication

## List on the endpoints and the details of the request and response

### User Registration

```javascript
POST /register
{
	"name": "John Doe",
	"email": "john@doe.com",
	"password": "password" // Password needs to be hashed before storing it in the database
}
```

```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
}
```
The token can be either a JWT or a random string that can be used for authentication.

### User Login

```javascript
POST /login
{
	"email": "john@doe.com",
	"password": "password"
}
```

```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
}
```
Email and password is validated and respond with a token if the authentication is successful.

### Create a Task Item

```javascript
POST /tasks
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}
```
Use must send the token received from the login endpoint in the header to authenticate the request.

```json
{
	"message": "Unauthorized"
}
```
In case the token is missing or invalid respond with an error and status code 401.

```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}
```
Upon successful creation of the to-do item, respond with the details of the created item.

### Update a Task Item

```javascript
PUT /tasks/1
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, bread, and cheese"
}
```
Use must send the token received from the login endpoint in the header to authenticate the request.

```json
{
	"message": "Forbidden"
}
```
Respond with an error and status code 403 if the user is not authorized to update the item.

```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Buy milk, eggs, bread, and cheese"
}
```
Upon success update of the task item, respond with the updated details of the item.

### Delete a Task Item

```javascript
DELETE /tasks/1
```
User must be authenticated and authorized to delete the task item. Upon successful deletion, respond with the status code 204.

### Get Task Items

```javascript
GET /tasks?page=1&limit=10
```
User must be authenticated to access the tasks and the response should be paginated.

```json
{
  "data": [
    {
      "id": 1,
      "title": "Buy groceries",
      "description": "Buy milk, eggs, bread"
    },
    {
      "id": 2,
      "title": "Pay bills",
      "description": "Pay electricity and water bills"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 2
}
```
Respond with the list of task items along with the pagination details.

## Features

* -

---

## Installation

Clone the repository:

```bash
git clone https://github.com/tspaja2025/go_task_tracker_api.git
cd go_task_tracker_api
```

Run the application:

```bash
go run main.go
```

---

## Usage

```bash
-
```
---

## Example Output

```text
-
```

---

## Data Storage

-

Example:

```json
-
```

---

## Technologies Used

* Go
* JSON file storage
* Standard library packages:

  * `-`

## Learning Goals

This project was built to practice:

* User authentication
* Schema design and Databases
* RESTful API design
* CRUD operations
* Error handling
* Security

---

## License

This project is open source and available under the MIT License.
