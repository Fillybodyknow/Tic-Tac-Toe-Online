// public/lobby.js

const HTTP_SERVER_URL = "https://tic-tac-toe-online-backend-20dy.onrender.com";
const PLAYER_ID = localStorage.getItem("player_id");

// Check if player_id exists, otherwise redirect to home
if (!PLAYER_ID) {
    window.location.href = "/";
}

const socket = io(HTTP_SERVER_URL, {
    query: {
        player_id: PLAYER_ID
    },
    transports: ["websocket", "polling"],
    reconnection: true,         // ✅ เปิด auto reconnect
    reconnectionAttempts: 5,
    reconnectionDelay: 1000,
    withCredentials: true,
});

// DOM Elements
const messageBox = document.getElementById("message");
const roomList = document.getElementById("room-list");
const createRoomBtn = document.getElementById("create-room-btn");

// Player Data Display Elements
const playerIdSpan = document.getElementById("player-id");
const playerNameSpan = document.getElementById("player-name");
const playerWonSpan = document.getElementById("player-won");
const playerLostSpan = document.getElementById("player-lost");
const playerDrawSpan = document.getElementById("player-draw");

// --- Socket.IO Event Listeners ---

// On successful connection and initial room list update
socket.on("connection", (rooms) => {
    getPlayerData(); // Fetch and display player stats
    updateRoomList(rooms);
});

// After creating a room successfully
socket.on("create-room-successfuly", (data) => {
    console.log("Room created:", data.room.room_id, "by player:", data.player_id);
    // When a new room is created, just update the whole list to simplify logic
    // This ensures correct ordering and "no rooms" message handling
    socket.emit("get-all-rooms"); // Request updated list from server
    showMessage(`Room ${data.room.room_id} created successfully!`, "success");

    // If this player created the room, automatically join it
    if (data.player_id === PLAYER_ID) {
        joinRoom(data.room.room_id);
    }
});

// When a room is updated (e.g., player joins/leaves, game starts)
socket.on("room-updated", (updatedRoom) => {
    // Re-fetch all rooms to ensure the list is always fresh and correctly ordered.
    // This simplifies client-side logic for individual room updates.
    socket.emit("get-all-rooms");
});

// When a room is removed (e.g., both players leave)
socket.on("room-removed", (roomId) => {
    // Re-fetch all rooms to ensure the list is always fresh and correctly ordered.
    socket.emit("get-all-rooms");
    showMessage(`Room ${roomId} has been closed.`, "info");
});

// Assuming your server also emits an event for "all-rooms-list" when requested or on initial connection
// If your server sends 'connection' event with rooms, this 'all-rooms-list' might be redundant
// or could be used for specific updates. Adjust based on your backend.
socket.on("all-rooms-list", (rooms) => {
    updateRoomList(rooms);
});


// --- UI Event Listeners ---

// Create Room Button click
createRoomBtn.addEventListener("click", () => {
    socket.emit("create-room");
});

// --- Helper Functions ---

// Fetches and displays player data
function getPlayerData() {
    fetch(`${HTTP_SERVER_URL}/api/player/${PLAYER_ID}`)
        .then(res => {
            if (!res.ok) {
                throw new Error(`HTTP error! status: ${res.status}`);
            }
            return res.json();
        })
        .then(data => {
            if (data && data.player) {
                playerIdSpan.textContent = data.player.player_id;
                playerNameSpan.textContent = data.player.player_name;
                playerWonSpan.textContent = data.player.win || 0;
                playerLostSpan.textContent = data.player.lose || 0;
                playerDrawSpan.textContent = data.player.draw || 0;
            } else {
                throw new Error("Invalid player data received.");
            }
        })
        .catch(err => {
            console.error("Failed to fetch player data:", err);
            showMessage("Could not load player data. Please log in again.", "danger");
            setTimeout(() => {
                window.location.href = "/";
            }, 2000);
        });
}

// Updates the entire room list by clearing and re-populating
function updateRoomList(rooms) {
    roomList.innerHTML = ''; // Clear existing rooms

    if (rooms.length === 0) {
        // If no rooms, add the "no rooms" message
        const noRoomsLi = document.createElement("li");
        noRoomsLi.className = "list-group-item text-center text-muted py-4";
        noRoomsLi.id = "no-rooms-message";
        noRoomsLi.textContent = "No active rooms found. Be the first to create one! 🤔";
        roomList.appendChild(noRoomsLi);
    } else {
        // If rooms exist, add each room
        rooms.forEach(room => {
            addRoomToList(room); // Use appendChild here
        });
    }
}

// Adds a single room to the list (always appends)
function addRoomToList(room) {
    const li = document.createElement("li");
    li.className = "list-group-item d-flex justify-content-between align-items-center";
    li.dataset.roomId = room.room_id; // Custom data attribute

    const playerXName = room.player_x?.player_name || 'Waiting...';
    const playerOName = room.player_o?.player_name || 'Waiting...';
    const isGameStarted = room.player_x?.player_name && room.player_o?.player_name;
    const isRoomFull = room.player_x?.player_id && room.player_o?.player_id;

    let actionButton = '';
    let statusBadge = '';

    if (isGameStarted) {
        statusBadge = `<span class="badge bg-danger rounded-pill">Game Started ⚔️</span>`;
        actionButton = '';
    } else if (isRoomFull) {
        statusBadge = `<span class="badge bg-info rounded-pill">Room Full 🔒</span>`;
        actionButton = '';
    } else if (room.player_x?.player_id === PLAYER_ID || room.player_o?.player_id === PLAYER_ID) {
        statusBadge = `<span class="badge bg-warning rounded-pill">Waiting for Opponent...</span>`;
        actionButton = '';
    } else {
        actionButton = `<button class="btn btn-info" onclick="joinRoom('${room.room_id}')">Join Game ▶️</button>`;
    }

    li.innerHTML = `
        <div class="room-info">
            Room ID: <strong>${room.room_id}</strong> <br>
            <small class="text-muted">
                ${playerXName} <span class="text-secondary mx-1">vs</span> ${playerOName}
            </small>
        </div>
        <div class="room-actions">
            ${statusBadge}
            ${actionButton}
        </div>
    `;

    roomList.appendChild(li); // Always append, as the list is cleared and rebuilt
}

// Redirects to game room
function joinRoom(roomId) {
    localStorage.setItem("room_id", roomId);
    window.location.href = "game_room.html";
}

// Shows a temporary message at the bottom of the screen
function showMessage(msg, type = "info", duration = 3000) {
    messageBox.textContent = msg;
    messageBox.className = `alert alert-${type} text-center d-block`;
    setTimeout(() => {
        messageBox.classList.add("d-none");
    }, duration);
}

// Initial fetch of player data when the page loads
getPlayerData();