type LetterResult = "correct" | "present" | "absent";

interface Guess {
  word: string;
  result: LetterResult[];
}

interface Session {
  session_id: string;
  guesses: Guess[];
  max_guesses: number;
  word_length: number;
  // only sent once the game is over
  hidden_word?: string;
}

const sessionStorageKey = "gordle-session-id";

// must match the flip animation in index.astro
const flipDelayMilliseconds = 250;
const flipDurationMilliseconds = 500;

const prefersReducedMotion = window.matchMedia(
  "(prefers-reduced-motion: reduce)",
).matches;

const board = document.querySelector<HTMLDivElement>("#board")!;
const keys = document.querySelectorAll<HTMLButtonElement>(".key");
const toast = document.querySelector<HTMLDivElement>("#toast")!;
const newGameButton = document.querySelector<HTMLButtonElement>("#new-game")!;
const newGameDialog =
  document.querySelector<HTMLDialogElement>("#new-game-dialog")!;
const newGameForm = document.querySelector<HTMLFormElement>("#new-game-form")!;
const newGameCancelButton =
  document.querySelector<HTMLButtonElement>("#new-game-cancel")!;
const wordLengthInput =
  document.querySelector<HTMLInputElement>("#word-length")!;
const wordLengthOutput =
  document.querySelector<HTMLOutputElement>("#word-length-value")!;
const maxGuessesInput =
  document.querySelector<HTMLInputElement>("#max-guesses")!;
const maxGuessesOutput =
  document.querySelector<HTMLOutputElement>("#max-guesses-value")!;

let session: Session | null = null;
let currentGuess = "";
// true while a guess is being checked and revealed, so key presses are ignored
let busy = false;
let toastTimeout: number | undefined;

function isGameOver(game: Session): boolean {
  return game.hidden_word !== undefined;
}

function wait(milliseconds: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, milliseconds));
}

// the toast is hidden by css whenever it is empty
function showToast(message: string, { stayVisible = false } = {}) {
  window.clearTimeout(toastTimeout);
  toast.textContent = message;
  if (!stayVisible) {
    toastTimeout = window.setTimeout(hideToast, 1500);
  }
}

function hideToast() {
  window.clearTimeout(toastTimeout);
  toast.textContent = "";
}

function showEndMessage(game: Session) {
  const lastGuess = game.guesses[game.guesses.length - 1];
  if (lastGuess.word === game.hidden_word) {
    showToast(`Solved in ${game.guesses.length}/${game.max_guesses}`, {
      stayVisible: true,
    });
  } else {
    showToast(game.hidden_word!.toUpperCase(), { stayVisible: true });
  }
}

// an empty grid with one row of tiles per allowed guess
function buildBoard(wordLength: number, maxGuesses: number) {
  board.style.setProperty("--columns", String(wordLength));
  board.style.setProperty("--rows", String(maxGuesses));
  board.replaceChildren();

  for (let rowIndex = 0; rowIndex < maxGuesses; rowIndex++) {
    const row = document.createElement("div");
    row.className = "row";
    for (let columnIndex = 0; columnIndex < wordLength; columnIndex++) {
      const tile = document.createElement("div");
      tile.className = "tile";
      tile.style.setProperty("--index", String(columnIndex));
      row.append(tile);
    }
    board.append(row);
  }
}

// fill the tiles with the submitted guesses, followed by the guess being typed
function renderBoard(guesses: Guess[]) {
  board.querySelectorAll(".row").forEach((row, rowIndex) => {
    const guess = guesses[rowIndex];
    const isCurrentRow = rowIndex === guesses.length;

    row.querySelectorAll<HTMLElement>(".tile").forEach((tile, columnIndex) => {
      if (guess) {
        tile.textContent = guess.word[columnIndex];
        tile.dataset.state = guess.result[columnIndex];
      } else if (isCurrentRow && columnIndex < currentGuess.length) {
        tile.textContent = currentGuess[columnIndex];
        tile.dataset.state = "filled";
      } else {
        tile.textContent = "";
        tile.dataset.state = "empty";
      }
    });
  });
}

// colour each key by the best result its letter has had so far
function renderKeyboard(guesses: Guess[]) {
  const resultRank = { absent: 1, present: 2, correct: 3 };
  const bestResults = new Map<string, LetterResult>();

  for (const guess of guesses) {
    guess.result.forEach((result, index) => {
      const letter = guess.word[index];
      const bestSoFar = bestResults.get(letter);
      if (!bestSoFar || resultRank[result] > resultRank[bestSoFar]) {
        bestResults.set(letter, result);
      }
    });
  }

  for (const key of keys) {
    const result = bestResults.get(key.dataset.key!);
    if (result) {
      key.dataset.state = result;
    } else {
      delete key.dataset.state;
    }
  }
}

function shake(row: HTMLElement) {
  row.classList.remove("shake");
  // force a reflow so the animation restarts if the row is already shaking
  void row.offsetWidth;
  row.classList.add("shake");
}

function loadSession(newSession: Session) {
  session = newSession;
  currentGuess = "";
  try {
    localStorage.setItem(sessionStorageKey, newSession.session_id);
  } catch {
    // storage is blocked, so the game works but won't survive a reload
  }

  hideToast();
  buildBoard(newSession.word_length, newSession.max_guesses);
  renderBoard(newSession.guesses);
  renderKeyboard(newSession.guesses);
  if (isGameOver(newSession)) {
    showEndMessage(newSession);
  }
}

async function startNewGame(wordLength: number, maxGuesses: number) {
  try {
    const response = await fetch("/api/session", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        word_length: wordLength,
        max_guesses: maxGuesses,
      }),
    });
    if (!response.ok) {
      throw new Error(`create session failed with ${response.status}`);
    }
    loadSession(await response.json());
  } catch {
    showToast("Could not start a game");
  }
}

async function submitGuess() {
  if (!session) {
    return;
  }
  const row = board.children[session.guesses.length] as HTMLElement;

  if (currentGuess.length < session.word_length) {
    shake(row);
    showToast("Not enough letters");
    return;
  }

  busy = true;
  try {
    const response = await fetch(`/api/session/${session.session_id}/guess`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ guess: currentGuess }),
    });
    const body = await response.json();

    if (response.ok) {
      await revealGuess(row, body);
    } else if (body.reason === "invalid_word") {
      shake(row);
      showToast("Not in word list");
    } else if (body.reason === "session_not_found") {
      showToast("Game expired");
      session = null;
      newGameDialog.showModal();
    } else {
      showToast("Something went wrong");
    }
  } catch {
    showToast("Something went wrong");
  } finally {
    busy = false;
  }
}

async function revealGuess(row: HTMLElement, updatedSession: Session) {
  session = updatedSession;
  currentGuess = "";
  row.classList.add("reveal");
  renderBoard(updatedSession.guesses);

  // like wordle, keep the keyboard and result hidden until every tile has flipped
  if (!prefersReducedMotion) {
    const lastTileDelay =
      (updatedSession.word_length - 1) * flipDelayMilliseconds;
    await wait(lastTileDelay + flipDurationMilliseconds);
  }

  renderKeyboard(updatedSession.guesses);
  if (isGameOver(updatedSession)) {
    showEndMessage(updatedSession);
  }
}

function handleKey(key: string) {
  if (!session || busy || isGameOver(session)) {
    return;
  }

  if (key === "enter") {
    submitGuess();
  } else if (key === "backspace") {
    currentGuess = currentGuess.slice(0, -1);
    renderBoard(session.guesses);
  } else if (/^[a-z]$/.test(key) && currentGuess.length < session.word_length) {
    currentGuess += key;
    renderBoard(session.guesses);
  }
}

document.addEventListener("keydown", (event) => {
  if (newGameDialog.open || event.ctrlKey || event.metaKey || event.altKey) {
    return;
  }
  // otherwise enter would also click whichever button has focus
  if (event.key === "Enter") {
    event.preventDefault();
  }
  handleKey(event.key.toLowerCase());
});

for (const key of keys) {
  key.addEventListener("click", () => handleKey(key.dataset.key!));
}

// wait for any guess to finish revealing, so its result can't land on the new game
newGameButton.addEventListener("click", () => {
  if (!busy) {
    newGameDialog.showModal();
  }
});

newGameCancelButton.addEventListener("click", () => newGameDialog.close());

// closing the dialog hands focus back to the new game button, where space would reopen it
newGameDialog.addEventListener("close", () => newGameButton.blur());

wordLengthInput.addEventListener("input", () => {
  wordLengthOutput.value = wordLengthInput.value;
});

maxGuessesInput.addEventListener("input", () => {
  maxGuessesOutput.value = maxGuessesInput.value;
});

// method="dialog" closes the dialog when the form is submitted
newGameForm.addEventListener("submit", () => {
  startNewGame(wordLengthInput.valueAsNumber, maxGuessesInput.valueAsNumber);
});

// carry on with the game saved in this browser, or ask for a new one
async function resumeOrAskForNewGame() {
  try {
    const savedSessionID = localStorage.getItem(sessionStorageKey);
    if (savedSessionID) {
      const response = await fetch(`/api/session/${savedSessionID}`);
      if (response.ok) {
        loadSession(await response.json());
        return;
      }
    }
  } catch {
    // storage is blocked or the server is unreachable, so start afresh
  }
  newGameDialog.showModal();
}

resumeOrAskForNewGame();
