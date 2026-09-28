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
  hidden_word?: string;
}

const sessionStorageKey = "gordle-session-id";

// used to keep only the strongest result per letter on the keyboard
const resultRank: Record<LetterResult, number> = {
  absent: 1,
  present: 2,
  correct: 3,
};

const flipStaggerMilliseconds = 250;
const flipDurationMilliseconds = 500;

const board = document.querySelector<HTMLDivElement>("#board")!;
const keyboard = document.querySelector<HTMLDivElement>("#keyboard")!;
const toast = document.querySelector<HTMLDivElement>("#toast")!;
const newGameButton = document.querySelector<HTMLButtonElement>("#new-game")!;
const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");

let session: Session | null = null;
let currentGuess = "";
let busy = false;
let toastTimeout: number | undefined;

function isComplete(game: Session): boolean {
  return game.hidden_word !== undefined;
}

function readStoredSessionID(): string | null {
  try {
    return localStorage.getItem(sessionStorageKey);
  } catch {
    return null;
  }
}

function storeSessionID(sessionID: string) {
  try {
    localStorage.setItem(sessionStorageKey, sessionID);
  } catch {
    // the game still works without persistence, it just won't survive a refresh
  }
}

function showToast(message: string, persistent = false) {
  window.clearTimeout(toastTimeout);
  toast.textContent = message;
  toast.classList.add("visible");
  if (!persistent) {
    toastTimeout = window.setTimeout(
      () => toast.classList.remove("visible"),
      1500,
    );
  }
}

function hideToast() {
  window.clearTimeout(toastTimeout);
  toast.classList.remove("visible");
}

function wait(milliseconds: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, milliseconds));
}

// create an empty board sized to the session's word length and guess count
function buildBoard(game: Session) {
  board.style.setProperty("--columns", String(game.word_length));
  board.style.setProperty("--rows", String(game.max_guesses));

  const rows = [];
  for (let rowIndex = 0; rowIndex < game.max_guesses; rowIndex++) {
    const row = document.createElement("div");
    row.className = "row";
    row.addEventListener("animationend", () => row.classList.remove("shake"));
    for (let columnIndex = 0; columnIndex < game.word_length; columnIndex++) {
      const tile = document.createElement("div");
      tile.className = "tile";
      tile.style.setProperty("--index", String(columnIndex));
      row.append(tile);
    }
    rows.push(row);
  }
  board.replaceChildren(...rows);
}

function renderBoard(game: Session) {
  Array.from(board.children).forEach((row, rowIndex) => {
    const guess = game.guesses[rowIndex];
    const isCurrentRow = rowIndex === game.guesses.length;

    Array.from(row.children).forEach((tileNode, columnIndex) => {
      const tile = tileNode as HTMLDivElement;
      let letter = "";
      let state = "empty";
      if (guess) {
        letter = guess.word[columnIndex];
        state = guess.result[columnIndex];
      } else if (isCurrentRow && columnIndex < currentGuess.length) {
        letter = currentGuess[columnIndex];
        state = "filled";
      }
      tile.textContent = letter;
      tile.dataset.state = state;
    });
  });
}

function renderKeyboard(game: Session) {
  const bestResults = new Map<string, LetterResult>();
  for (const guess of game.guesses) {
    guess.result.forEach((result, index) => {
      const letter = guess.word[index];
      const known = bestResults.get(letter);
      if (!known || resultRank[result] > resultRank[known]) {
        bestResults.set(letter, result);
      }
    });
  }

  for (const key of keyboard.querySelectorAll<HTMLButtonElement>(".key")) {
    const result = bestResults.get(key.dataset.key!);
    if (result) {
      key.dataset.state = result;
    } else {
      delete key.dataset.state;
    }
  }
}

function showEndMessage(game: Session) {
  const lastGuess = game.guesses[game.guesses.length - 1];
  if (lastGuess?.word === game.hidden_word) {
    showToast(`Solved in ${game.guesses.length}/${game.max_guesses}`, true);
  } else {
    showToast(game.hidden_word!.toUpperCase(), true);
  }
}

function loadSession(game: Session) {
  session = game;
  currentGuess = "";
  storeSessionID(game.session_id);
  hideToast();
  buildBoard(game);
  renderBoard(game);
  renderKeyboard(game);
  if (isComplete(game)) {
    showEndMessage(game);
  }
}

async function createSession(): Promise<Session> {
  const response = await fetch("/api/session", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: "{}",
  });
  if (!response.ok) {
    throw new Error(`create session failed with ${response.status}`);
  }
  return response.json();
}

async function fetchSession(sessionID: string): Promise<Session | null> {
  const response = await fetch(`/api/session/${encodeURIComponent(sessionID)}`);
  return response.ok ? response.json() : null;
}

async function startNewGame() {
  if (busy) {
    return;
  }
  busy = true;
  try {
    loadSession(await createSession());
  } catch {
    showToast("Could not start a game");
  } finally {
    busy = false;
  }
}

function shakeCurrentRow() {
  if (!session) {
    return;
  }
  const row = board.children[session.guesses.length];
  row.classList.remove("shake");
  // force a reflow so the animation restarts when shaking twice in a row
  void (row as HTMLElement).offsetWidth;
  row.classList.add("shake");
}

async function submitGuess() {
  if (!session) {
    return;
  }
  if (currentGuess.length < session.word_length) {
    shakeCurrentRow();
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

    if (!response.ok) {
      if (body.reason === "invalid_word") {
        shakeCurrentRow();
        showToast("Not in word list");
      } else if (body.reason === "session_not_found") {
        showToast("Game expired, starting a new one");
        loadSession(await createSession());
      } else {
        showToast("Something went wrong");
      }
      return;
    }

    const updatedSession: Session = body;
    const revealedRow = board.children[session.guesses.length];
    session = updatedSession;
    currentGuess = "";
    revealedRow.classList.add("reveal");
    renderBoard(updatedSession);

    // keep the keyboard and result hidden until the tiles have flipped, like wordle
    if (!reducedMotion.matches) {
      await wait(
        (updatedSession.word_length - 1) * flipStaggerMilliseconds +
          flipDurationMilliseconds,
      );
    }
    renderKeyboard(updatedSession);
    if (isComplete(updatedSession)) {
      showEndMessage(updatedSession);
    }
  } catch {
    showToast("Something went wrong");
  } finally {
    busy = false;
  }
}

function handleKey(key: string) {
  if (!session || busy || isComplete(session)) {
    return;
  }

  if (key === "enter") {
    submitGuess();
  } else if (key === "backspace") {
    currentGuess = currentGuess.slice(0, -1);
    renderBoard(session);
  } else if (/^[a-z]$/.test(key) && currentGuess.length < session.word_length) {
    currentGuess += key;
    renderBoard(session);
  }
}

document.addEventListener("keydown", (event) => {
  if (event.ctrlKey || event.metaKey || event.altKey) {
    return;
  }
  const key = event.key.toLowerCase();
  if (key === "enter" || key === "backspace" || /^[a-z]$/.test(key)) {
    // stops enter from also clicking whichever button has focus
    event.preventDefault();
    handleKey(key);
  }
});

keyboard.addEventListener("click", (event) => {
  const key = (event.target as HTMLElement).closest<HTMLButtonElement>(".key");
  if (key) {
    handleKey(key.dataset.key!);
  }
});

newGameButton.addEventListener("click", () => {
  newGameButton.blur();
  startNewGame();
});

async function initialise() {
  const storedSessionID = readStoredSessionID();
  const storedSession = storedSessionID
    ? await fetchSession(storedSessionID).catch(() => null)
    : null;
  if (storedSession) {
    loadSession(storedSession);
  } else {
    await startNewGame();
  }
}

initialise();
