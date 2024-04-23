
# Online Gaming Leaderboard API

The Online Gaming Leaderboard API is a RESTful API for managing highscores in an online gaming site. It provides endpoints to submit scores, retrieve user ranks, and list top N scores based on different scopes such as country, state, or globally.

## Routes

### Submit Score

- **URL:** `/submit`
- **Method:** POST
- **Description:** Submit a score to the system.
- **Payload:**
  ```json
  {
      "user_name": "string",
      "country": "string",
      "state": "string",
      "score": int
  }
  ```
- **Response:**
  - 200 OK: Score submitted successfully.
  - 400 Bad Request: Invalid request.

### Get Rank

- **URL:** `/get_rank`
- **Method:** GET
- **Description:** Get rank of a user.
- **Query Parameters:**
  - `user_name` (string, required): User name.
  - `scope` (string, required): Scope of ranking (state, country, or globally).
- **Response:**
  - 200 OK: User rank.
  - 400 Bad Request: Invalid request.

### List Top N

- **URL:** `/list_top_n`
- **Method:** GET
- **Description:** List top N ranks.
- **Query Parameters:**
  - `n` (integer, required): Number of ranks to list.
  - `scope` (string, required): Scope of ranking (state, country, or globally).
- **Response:**
  - 200 OK: Top N scores.
  - 400 Bad Request: Invalid request.

## Payload Structure

### Score Object

```json
{
    "user_name": "string",
    "country": "string",
    "state": "string",
    "score": int
}
```

- `user_name` (string): User name.
- `country` (string): Country of the user.
- `state` (string): State of the user.
- `score` (int): Score achieved by the user.

---
