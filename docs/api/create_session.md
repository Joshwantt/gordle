# `create_session`

```
POST /api/session
```

The `create_session` endpoint is used to create a session, and should be used to start a new game.

## Request

```http
POST /api/session HTTP/1.1
Content-Type: application/json

{
    "word_length": 5,
    "max_guesses": 6
}
```

**Request Body:**
| Field | Type | Description |
|-------|------|-------------|
| `word_length` | Integer | How many letters are in the hidden word. Must be between `3` and `16` inclusive. |
| `max_guesses` | Integer | How many guesses the player may make before losing. Must be between `1` and `20` inclusive. |

## Response

### Success

On success, the server responds with `201 Created` and the full initial game state is returned. The response is identical in shape to [`get_session`](./get_session.md). Refer to that page for the complete field descriptions and examples.

### Failure

The status code will always be `4xx` on failure, as listed below.

```json
{
    "reason": "invalid_request"
}
```

| Status | `reason` | Meaning |
|--------|----------|---------|
| `400 Bad Request` | `invalid_request` | The request was malformed, missing required arguments, or `word_length` or `max_guesses` is out of range. |