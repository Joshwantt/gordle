# Gordle API Documentation

Gordle uses a REST style API to manage game state between the client and the server.

## Sessions

Each session is one game of Gordle, following the normal Wordle rules: the player has a set number of guesses to find a hidden word of a set length, both chosen when the session is created. Sessions:
 - Have a unique ID which can be used to manipulate them
 - Store all the information about the game including the hidden word and what guesses have been made.

Until a session is complete (the player has won or lost the game) the hidden word is kept exclusively on the server side to prevent cheating, a notable improvement over traditional Wordle

## Conventions

All endpoints expect requests to:
 - Use the HTTP method matching the operation: `GET` to retrieve, `POST` to create.
 - Identify the session being operated on in the path, e.g. `/api/session/{session_id}`.
 - Where a payload is sent, be of the `application/json` mimetype.
 - Where a payload is sent, have a payload made up of valid JSON.

All responses from endpoints will:
 - Give an HTTP status code which may be used to determine if the operation was successful: `2xx` on success, `4xx` or `5xx` on failure. Each endpoint lists the codes it may return.
 - Be of the `application/json` mimetype.
 - Have a payload made up of valid JSON which, on failure:
   - May (or may not) contain a string `reason`, which is defined per-endpoint. Clients may use this for error handling.
   - May (or may not) contain a string `msg`, which will contain a human readable message with more information on the status of the request. Useful for debugging.

## API Endpoints

 - [`create_session`](./create_session.md) (`POST /api/session`): Create a new session.
 - [`get_session`](./get_session.md) (`GET /api/session/{session_id}`): Obtain information about an existing session.
 - [`make_guess`](./make_guess.md) (`POST /api/session/{session_id}/guess`): Make a guess in an existing session.
