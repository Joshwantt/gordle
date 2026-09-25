# `create_game`

```
POST /api/games
```

The `create_game` endpoint is used to create a game, and should be used to start a new game.

## Request

This endpoint does not take any parameters, the payload will be ignored. It is acceptable to have an empty payload.

```http
POST /api/games HTTP/1.1
Content-Type: application/json

{}
```

## Response

### Success

On success, the server responds with `201 Created` and the full initial game state is returned. The response is identical in shape to [`get_game`](./get_game.md). Refer to that page for the complete field descriptions and examples.

### Failure

No reason is supplied for failure. This operation is intended to always succeed.

The status code will always be `500 Internal Server Error` on failure.

```json
{}
```
