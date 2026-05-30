1. First run

```bash
docker-compose up --build
http://localhost:8080
docker-compose down or Ctrl + C
```

2. Apply changes

```bash
docker-compose down -v
docker-compose up --build
```

3. Run and test

```bash
docker-compose down -v
docker-compose up --build
curl http://localhost:8080/health
```

4. Test Validation

```bash
docker-compose up --build
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test123", "email": "test@golang.org", "password": "securepassword123"}'
```

5. Access Token & Refresh Token

```bash
docker-compose up --build
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@golang.org", "password": "securepassword123"}'
```

Expected response:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

6. Update the repository layer

```bash
docker-compose up --build
curl -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "YOUR_COPIED_REFRESH_TOKEN_HERE"}'
```

7. Auth middleware

Step A:
```bash
docker-compose up --build
curl -i http://localhost:8080/tasks/test
```

Step B:
```bash
curl -X GET http://localhost:8080/tasks/test \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_STRING_HERE"
```

Expected response:

```json
{
  "message": "Access granted!",
  "your_user_id": 1
}
```

8. Verify CRUD system

```bash
docker-compose up --build
```

Create a task
```bash
curl -X POST http://localhost:8080/tasks/ \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Docker", "description": "Finish setting up the containerized api layers"}'
```

Read a task
```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" http://localhost:8080/tasks/1
```

Update a task status
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

Delete a task
```bash
curl -X DELETE -H "Authorization: Bearer YOUR_ACCESS_TOKEN" http://localhost:8080/tasks/1
```

9. Test query system

```bash
docker-compose up --build
```

Get page 1 with a limit of 5 tasks
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" "http://localhost:8080/tasks?page=1&limit=5"
```

Filter only "pending" tasks, sorted by priority in ascending order
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" "http://localhost:8080/tasks?status=pending&sort_by=priority&order=ASC"
```

Conbine filtering and pagination criteria
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" "http://localhost:8080/tasks?status=completed&priority=high&page=1&limit=2"
```
