# `get_game`

```
GET /api/games/{game_id}
```

The `get_game` endpoint is used to retrieve information about an existing game, including guesses that have been made.

## Request

```http
GET /api/games/c0a95416-a39e-4adb-a7fc-b05f90e61ec4 HTTP/1.1
```

**Path Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| `game_id` | String | The UUID of the game to retrieve. |

## Response

On success, the server responds with `200 OK`.

### Success - Fresh Game

Example response from a fresh game, with no guesses.

```json
{
    "game_id": "c0a95416-a39e-4adb-a7fc-b05f90e61ec4",
    "guesses": [],
    "max_guesses": 6,
    "word_length": 5
}
```

**Response Body:**
| Field | Type | Description |
|-------|------|-------------|
| `game_id` | String | The UUID of the game. |
| `guesses` | `Guess[]` | The guesses made so far, in the order they were made. On a fresh game, this will be empty. |
| `max_guesses` | Integer | How many guesses the player may make before losing. |
| `word_length` | Integer | How many letters are in the hidden word. Every guess must be this length. |

### Success - Game In Progress

Example response from a game in progress with two guesses made. The hidden word is "crumb", although this information is only known to the server.

```json
{
    "game_id": "c0a95416-a39e-4adb-a7fc-b05f90e61ec4",
    "guesses": [
        {
            "word": "zebra",
            "result": ["absent","absent","present","present","absent"]
        },
        {
            "word": "brink",
            "result": ["present","correct","absent","absent","absent"]
        }
    ],
    "max_guesses": 6,
    "word_length": 5
}
```

**Response Body:**
| Field | Type | Description |
|-------|------|-------------|
| `game_id` | String | The UUID of the game. |
| `guesses` | `Guess[]` | See below. |
| `max_guesses` | Integer | How many guesses the player may make before losing. |
| `word_length` | Integer | How many letters are in the hidden word. Every guess must be this length. |


### Success - Completed Game

Example response from a completed game where the player won by guessing the hidden word "crumb". When a game is complete, `hidden_word` is included in the response. A player wins by guessing the hidden word, and loses by making `max_guesses` guesses without doing so. A winning game's final guess will always show all `"correct"` results.

```json
{
    "game_id": "c0a95416-a39e-4adb-a7fc-b05f90e61ec4",
    "guesses": [
        {
            "word": "zebra",
            "result": ["absent","absent","present","present","absent"]
        },
        {
            "word": "brink",
            "result": ["present","correct","absent","absent","absent"]
        },
        {
            "word": "crumb",
            "result": ["correct","correct","correct","correct","correct"]
        }
    ],
    "max_guesses": 6,
    "word_length": 5,
    "hidden_word": "crumb"
}
```

**Response Body:**
| Field | Type | Description |
|-------|------|-------------|
| `game_id` | String | The UUID of the game. |
| `guesses` | `Guess[]` | See below. |
| `max_guesses` | Integer | How many guesses the player may make before losing. |
| `word_length` | Integer | How many letters are in the hidden word. Every guess must be this length. |
| `hidden_word` | String | The word the player was trying to guess. Only present when the game is complete. |


### `Guess`

Each entry in the `guesses` array represents one guess made by the player, in order.

```json
{
    "word": "zebra",
    "result": ["absent", "absent", "present", "present", "absent"]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `word` | String | The word that was guessed. |
| `result` | `LetterResult[]` | The result for each letter in `word`, in order. See below. |


### `LetterResult`

Each entry in a guess's `result` array is a string enum describing one letter of the guess.

| Value | Colour | Meaning |
|-------|--------|---------|
| `"correct"` | Green | This letter is in the hidden word at this exact position. |
| `"present"` | Yellow | This letter is in the hidden word, but not at this position. |
| `"absent"` | Grey | This letter does not appear in the hidden word. |


### Failure

The status code will always be `4xx` on failure, as listed below.

```json
{
    "reason": "game_not_found"
}
```

| Status | `reason` | Meaning |
|--------|----------|---------|
| `400 Bad Request` | `invalid_request` | The request was malformed, such as a `game_id` that is not a valid UUID. |
| `404 Not Found` | `game_not_found` | No game exists with the provided `game_id`. |
