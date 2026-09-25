# Gordle API Documentation

Gordle uses a REST style API to manage game state between the client and the server.

## Games

Each game of Gordle follows the normal Wordle rules: the player has six guesses to find a hidden five letter word. Games:
 - Have a unique ID which can be used to manipulate them
 - Store all the information about a game including the hidden word and what guesses have been made.

Until a game is complete (the player has won or lost the game) the hidden word is kept exclusively on the server side to prevent cheating, a notable improvement over traditional Wordle

## Conventions

All endpoints expect requests to:
 - Use the HTTP method matching the operation: `GET` to retrieve, `POST` to create.
 - Identify the game being operated on in the path, e.g. `/api/games/{game_id}`.
 - Where a payload is sent, be of the `application/json` mimetype.
 - Where a payload is sent, have a payload made up of valid JSON.

All responses from endpoints will:
 - Give an HTTP status code which may be used to determine if the operation was successful: `2xx` on success, `4xx` or `5xx` on failure. Each endpoint lists the codes it may return.
 - Be of the `application/json` mimetype.
 - Have a payload made up of valid JSON which, on failure:
   - May (or may not) contain a string `reason`, which is defined per-endpoint. Clients may use this for error handling.
   - May (or may not) contain a string `msg`, which will contain a human readable message with more information on the status of the request. Useful for debugging.

## API Endpoints

 - [`create_game`](./create_game.md) (`POST /api/games`): Create a new game.
 - [`get_game`](./get_game.md) (`GET /api/games/{game_id}`): Obtain information about an existing game.
 - [`make_guess`](./make_guess.md) (`POST /api/games/{game_id}/guess`): Make a guess in an existing game.
