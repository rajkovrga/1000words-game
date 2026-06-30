# 1000words-game

CLI and API application for learning languages through the 1000 most common words.

The goal of the application is to allow users to learn **any language that exists in the system**.  
The system is designed so that languages, words, levels, progress, and game attempts can be expanded over time without changing the main application logic.

The application has two modes:

- **Practice** — learning mode without login
- **Game** — authenticated mode with user progress, attempts, and level progression

Backend uses:

- Go
- SQLite
- goose migrations
- chi router
- token authentication
- Blake2b token hashing
- roles and permissions

---

## Starting the CLI application

The CLI application is started with:

```bash
bash scripts/start.sh
```

The script:

1. checks Go installation
2. checks/installs goose
3. downloads Go dependencies
4. runs database migrations
5. builds the CLI application
6. starts the application

---

## Starting the API application

The API application is started with:

```bash
bash scripts/api-start.sh
```

The script:

1. checks Go installation
2. checks/installs goose
3. downloads Go dependencies
4. runs database migrations
5. builds the API application from `./cmd/api`
6. starts the API server

Default API address:

```text
http://localhost:8080
```

---

## .env example

The `.env` file should contain the following basic values:

```env
APP_NAME=1000words-game
DB_PATH=words.db
WORDS_PER_LEVEL=100

API_PORT=8080
API_READ_TIMEOUT_SECONDS=10
API_WRITE_TIMEOUT_SECONDS=10
API_IDLE_TIMEOUT_SECONDS=60
API_TOKEN_DAYS=30
```

---

## Default admin

Seeded admin account:

```text
email: admin@1000words.local
password: 123456789
```

The admin user receives the role:

```text
admin
```

---

## API routes

Base URL:

```text
http://localhost:8080/api/v1
```

---

## Health

### GET `/api/health`

Public health route.

### GET `/api/v1/health`

Versioned health route.

Example:

```bash
curl http://localhost:8080/api/v1/health
```

---

## Auth

### POST `/api/v1/auth/register`

Registers a user and returns a token.

Body:

```json
{
  "email": "user@test.com",
  "password": "123456789"
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"123456789"}'
```

---

### POST `/api/v1/auth/login`

Logs in a user.

Body:

```json
{
  "email": "user@test.com",
  "password": "123456789"
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"123456789"}'
```

Response returns a token:

```json
{
  "success": true,
  "data": {
    "token": "...",
    "expires_at": "...",
    "user": {
      "id": 1,
      "email": "user@test.com"
    }
  }
}
```

---

### POST `/api/v1/auth/logout`

Logs out the user.

Requires token:

```text
Authorization: Bearer TOKEN
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer TOKEN"
```

---

## Me

All `me` routes require a token.

### GET `/api/v1/me`

Returns the currently authenticated user, including roles and permissions.

```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer TOKEN"
```

---

### GET `/api/v1/me/roles`

Returns the roles of the currently authenticated user.

```bash
curl http://localhost:8080/api/v1/me/roles \
  -H "Authorization: Bearer TOKEN"
```

---

### GET `/api/v1/me/permissions`

Returns the permissions of the currently authenticated user.

```bash
curl http://localhost:8080/api/v1/me/permissions \
  -H "Authorization: Bearer TOKEN"
```

---

## Languages

### GET `/api/v1/languages`

Public route. Returns available languages.

```bash
curl http://localhost:8080/api/v1/languages
```

---

## Levels

### GET `/api/v1/levels`

Public route. Returns available levels.

```bash
curl http://localhost:8080/api/v1/levels
```

---

## Practice

Practice is a learning mode without login.

### POST `/api/v1/practice/start`

Starts practice mode.

Body:

```json
{
  "level_number": 1,
  "target_language_code": "en",
  "native_language_code": "sr"
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/practice/start \
  -H "Content-Type: application/json" \
  -d '{"level_number":1,"target_language_code":"en","native_language_code":"sr"}'
```

Practice returns questions and answers because it is not a competitive mode.

---

## Progress

Progress routes require a token.

### GET `/api/v1/progress`

Returns the user's progress.

```bash
curl http://localhost:8080/api/v1/progress \
  -H "Authorization: Bearer TOKEN"
```

---

### POST `/api/v1/progress`

Adds a new language learning progress for the user.

Body:

```json
{
  "target_language_code": "en",
  "native_language_code": "sr"
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/progress \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"target_language_code":"en","native_language_code":"sr"}'
```

---

## Game

Game routes require a token.

Game uses the user's existing progress.

### POST `/api/v1/game/start`

Starts the game for a selected progress.

Body:

```json
{
  "progress_id": 1
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/game/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"progress_id":1}'
```

The response returns:

- `attempt_id`
- `progress_id`
- level
- allowed number of wrong answers
- questions
- answers
- total questions

The frontend checks answers and counts mistakes.

---

### POST `/api/v1/game/finish`

Finishes the game and saves the result.

Body example for a failed attempt:

```json
{
  "attempt_id": 1,
  "progress_id": 1,
  "correct_count": 0,
  "wrong_count": 3,
  "total_questions": 3
}
```

Body example for a passed attempt:

```json
{
  "attempt_id": 2,
  "progress_id": 1,
  "correct_count": 98,
  "wrong_count": 2,
  "total_questions": 100
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/v1/game/finish \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"attempt_id":1,"progress_id":1,"correct_count":98,"wrong_count":2,"total_questions":100}'
```

The backend validates that:

- the attempt belongs to the user
- the progress belongs to the user
- the attempt matches the progress
- the result values are not negative
- `correct_count + wrong_count = total_questions`
- the full attempt must be completed in order to pass
- if the user passes, progress moves to the next level

---

## Roles and permissions

There are three roles:

```text
user
moderator
admin
```

`user` can:

- log in
- use practice
- create progress
- play the game
- view own profile
- view own permissions

`moderator` is planned for a future dashboard and word/content management.

`admin` has all permissions.

---

## Running migrations manually

Run all migrations:

```bash
goose -dir database/migrations sqlite3 words.db up
```

Check migration status:

```bash
goose -dir database/migrations sqlite3 words.db status
```

Rollback the latest migration:

```bash
goose -dir database/migrations sqlite3 words.db down
```

---

## Useful Game API test flow

1. Login or register
2. Save the token
3. Create progress
4. Save `progress_id`
5. Start the game with `game/start`
6. Save `attempt_id`
7. Send result with `game/finish`

Example flow:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"123456789"}'
```

```bash
curl -X POST http://localhost:8080/api/v1/progress \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"target_language_code":"en","native_language_code":"sr"}'
```

```bash
curl -X POST http://localhost:8080/api/v1/game/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"progress_id":1}'
```

```bash
curl -X POST http://localhost:8080/api/v1/game/finish \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"attempt_id":1,"progress_id":1,"correct_count":98,"wrong_count":2,"total_questions":100}'
```

---

## Build folder

Built files are located in:

```text
build/
```

CLI build:

```text
build/1000words-game
```

API build:

```text
build/1000words-api
```

---

## Running the applications

CLI:

```bash
bash scripts/start.sh
```

API:

```bash
bash scripts/api-start.sh
```
