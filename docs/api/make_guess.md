# `make_guess`

```
POST /api/session/{session_id}/guess
```

The `make_guess` endpoint is used to make a guess in an existing session.

## Request

```http
POST /api/session/c0a95416-a39e-4adb-a7fc-b05f90e61ec4/guess HTTP/1.1
Content-Type: application/json

{
    "guess": "crane"
}
```

**Path Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| `session_id` | String | The UUID of the session to make a guess in. |

**Request Body:**
| Field | Type | Description |
|-------|------|-------------|
| `guess` | String | The word being guessed. |

## Response

### Success

On success, the server responds with `200 OK` and the full updated game state is returned. The response is identical in shape to [`get_session`](./get_session.md). Refer to that page for the complete field descriptions and examples.

### Failure

The status code will always be `4xx` on failure, as listed below.

```json
{
    "reason": "invalid_word"
}
```

| Status | `reason` | Meaning |
|--------|----------|---------|
| `400 Bad Request` | `invalid_request` | The request was malformed or missing required arguments. |
| `404 Not Found` | `session_not_found` | No session exists with the provided `session_id`. |
| `409 Conflict` | `session_complete` | The session is already over and no more guesses can be made. |
| `422 Unprocessable Content` | `invalid_word` | The guess is not a valid word, or is not `word_length` letters long. |
