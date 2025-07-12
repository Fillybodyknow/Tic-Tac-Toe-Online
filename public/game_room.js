const HTTP_SERVER_URL = "http://101.51.120.113:3000";
const PLAYER_ID = localStorage.getItem("player_id");
const ROOM_ID = localStorage.getItem("room_id");

let roomStatusDiv;
let countdownStatusDiv;
let turnStatusDiv;
let gameBoardWrapper;
let socket;

let boardElement = null;
let mySymbol = "";

let selectedSpecialPawn = "";

document.addEventListener('DOMContentLoaded', () => {
    roomStatusDiv = document.getElementById("room-status");
    countdownStatusDiv = document.getElementById("countdown-status");
    turnStatusDiv = document.getElementById("turn-status");
    gameBoardWrapper = document.getElementById("game-board-wrapper");
    specialPawnSelect = document.getElementById("special-pawn");

    if (!PLAYER_ID || !ROOM_ID) {
        alert("Player ID or Room ID is missing. Please go back to the lobby.");
        window.location.href = "lobby.html";
        return;
    }

    socket = io(HTTP_SERVER_URL + "/game-room", {
        query: {
            player_id: PLAYER_ID,
            room_id: ROOM_ID
        },
        transports: ["websocket", "polling"],
        reconnection: true,         // ✅ เปิด auto reconnect
    reconnectionAttempts: 5,
    reconnectionDelay: 1000,
        withCredentials: true,
    });

    socket.on("connect", () => {
        console.log(`Connected to game room ${ROOM_ID}`);
        socket.emit("request-game-state");
    });

    socket.on("join-room-failed", (err) => {
        roomStatusDiv.innerHTML = `<div class="alert alert-danger">${err}</div>`;
        setTimeout(() => {
            window.location.href = "lobby.html";
        }, 3000);
    });

    socket.on("game-ready", (msg) => {
        let countdown = 3;

        roomStatusDiv.innerHTML = `<h4 class="text-success">✅ ${msg}</h4>`;
        countdownStatusDiv.className = "fs-4 mt-2 text-info fw-bold";
        countdownStatusDiv.innerHTML = `Game starts in ${countdown}...`;

        const countdownInterval = setInterval(() => {
            countdown--;
            if (countdown >= 0) {
                countdownStatusDiv.innerHTML = `Game starts in ${countdown}...`;
            }

            if (countdown < 0) {
                clearInterval(countdownInterval);
                renderGameBoard();
                socket.emit("request-game-state");
                countdownStatusDiv.innerHTML = "";
            }
        }, 1000);
    });

    socket.on("game-not-ready", (msg) => {
        // This line should now be safe as roomStatusDiv is guaranteed to be assigned
        roomStatusDiv.innerHTML = `<h4 class="text-warning">⏳ ${msg}</h4>`;
        turnStatusDiv.innerHTML = "";
    });

    socket.on("update-board", (gameRoom) => {
    if (!boardElement) return;

    specialPawnSelect.innerHTML = "";
    if (gameRoom.player_x && gameRoom.player_x.player_id === PLAYER_ID) {
        mySymbol = "X";
        specialPawnSelect.innerHTML += `<option value="">X Small - ∞</option>`;
        if (gameRoom.special_pawn_x["X|medium|2"] > 0) {
            specialPawnSelect.innerHTML += `<option value="X|medium|2">X Medium - ${gameRoom.special_pawn_x["X|medium|2"]}</option>`;
        }
        if (gameRoom.special_pawn_x["X|large|3"] > 0) {
            specialPawnSelect.innerHTML += `<option value="X|large|3">X Large - ${gameRoom.special_pawn_x["X|large|3"]}</option>`;
        }
    } else if (gameRoom.player_o && gameRoom.player_o.player_id === PLAYER_ID) {
        mySymbol = "O";
        specialPawnSelect.innerHTML += `<option value="">O Small - ∞</option>`;
        if (gameRoom.special_pawn_o["O|medium|2"] > 0) {
            specialPawnSelect.innerHTML += `<option value="O|medium|2">O Medium - ${gameRoom.special_pawn_o["O|medium|2"]}</option>`;
        }
        if (gameRoom.special_pawn_o["O|large|3"] > 0) {
            specialPawnSelect.innerHTML += `<option value="O|large|3">O Large - ${gameRoom.special_pawn_o["O|large|3"]}</option>`;
        }
    } else {
        mySymbol = "";
    }

    // Update Board
    for (let i = 0; i < 9; i++) {
        const row = Math.floor(i / 3);
        const col = i % 3;
        const cell = boardElement.children[i];
        const symbol = gameRoom.board[row][col];

        // Reset styles
        cell.classList.remove('player-x', 'player-o', 'pawn-small', 'pawn-medium', 'pawn-large');

        // Extract baseSymbol and apply text
        let baseSymbol = symbol;
        if (symbol.includes("|")) {
            baseSymbol = symbol.split("|")[0];
        }
        cell.textContent = baseSymbol;

        // Apply color
        if (baseSymbol === "X") {
            cell.classList.add("player-x");
        } else if (baseSymbol === "O") {
            cell.classList.add("player-o");
        }

        // Apply size class
        if (symbol.includes("|medium|2")) {
            cell.classList.add("pawn-medium");
        } else if (symbol.includes("|large|3")) {
            cell.classList.add("pawn-large");
        } else if (symbol !== "") {
            cell.classList.add("pawn-small");
        }

        const isMyTurn = gameRoom.turn === PLAYER_ID;
    }

    // Status text
    if (gameRoom.turn === PLAYER_ID) {
        turnStatusDiv.innerHTML = `<span class="text-success">👉 It's your turn! Your Symbol: <span class="fw-bold">${mySymbol}</span></span>`;
    } else {
        turnStatusDiv.innerHTML = `<span class="text-secondary">⏳ Waiting for opponent's move...</span>`;
    }

    roomStatusDiv.innerHTML = "";
    countdownStatusDiv.innerHTML = "";
});


    socket.on("make-move-failed", (msg) => {
        alert(`Move failed: ${msg}`);
        socket.emit("request-game-state");
    });

    socket.on("game-winner", (winnerData) => {
        let winnerMessage = "";
        if (winnerData === PLAYER_ID) {
            winnerMessage = "🎉 You Won! Congratulations! 🏆";
            turnStatusDiv.innerHTML = `<span class="text-success fw-bold">${winnerMessage}</span>`;
        } else if (winnerData === "draw") {
            winnerMessage = "🤝 It's a Draw! Good Game!";
            turnStatusDiv.innerHTML = `<span class="text-info fw-bold">${winnerMessage}</span>`;
        } else if (winnerData === null || winnerData === undefined) {
            winnerMessage = "Game Ended!";
            turnStatusDiv.innerHTML = `<span class="text-info fw-bold">${winnerMessage}</span>`;
        } else {
            winnerMessage = `😞 Player ${winnerData} is the Winner! Better luck next time.`;
            turnStatusDiv.innerHTML = `<span class="text-danger fw-bold">${winnerMessage}</span>`;
        }

        alert(winnerMessage);
        disableBoard();
        showGoToLobbyButton();
    });

    socket.on("game-draw", () => {
        const message = "🤝 It's a Draw! Good Game!";
        alert(message);
        turnStatusDiv.innerHTML = `<span class="text-info fw-bold">${message}</span>`;
        disableBoard();
        showGoToLobbyButton();
    });

    socket.on("player-disconnected", (msg) => {
        alert(`⚠️ Opponent disconnected: ${msg}. Returning to lobby.`);
        turnStatusDiv.innerHTML = `<span class="text-warning fw-bold">Opponent Disconnected!</span>`;
        disableBoard();
        setTimeout(() => {
            window.location.href = "lobby.html";
        }, 3000);
    });

    function renderGameBoard() {
        if (!boardElement) {
            boardElement = document.createElement("div");
            boardElement.id = "game-board-container";
            gameBoardWrapper.appendChild(boardElement);
        } else {
            boardElement.innerHTML = '';
        }

        for (let i = 0; i < 9; i++) {
            const cell = document.createElement("button");
            cell.className = "game-cell";
            cell.dataset.index = i;

            cell.addEventListener("click", () => {
                const row = Math.floor(i / 3);
                const col = i % 3;

                const special_pawn = specialPawnSelect?.value || "";
                socket.emit("make-move", { row, col, special_pawn});
            });

            boardElement.appendChild(cell);
        }
    }

    function disableBoard() {
        if (!boardElement) return;
        Array.from(boardElement.children).forEach(cell => {
            cell.disabled = true;
        });
    }

    function showGoToLobbyButton() {
        let backButton = document.querySelector('.go-back-btn');
        if (!backButton) {
            const containerDiv = document.querySelector('.container');
            backButton = document.createElement('button');
            backButton.className = 'btn go-back-btn mt-4';
            backButton.textContent = '⬅️ Back to Lobby';
            backButton.onclick = () => { window.location.href = 'lobby.html'; };
            containerDiv.appendChild(backButton);
        }
        backButton.style.display = 'block';
    }

});