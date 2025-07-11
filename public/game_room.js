const statusRoom = document.getElementById("status-room");
const turnStatus = document.getElementById("turn-status");
const player_id = localStorage.getItem("player_id");
const room_id = localStorage.getItem("room_id");

const http = "http://localhost:3000";

let boardElement = null; // ✅ เพิ่มตัวแปร global

const socket = io(http + "/game-room", {
  query: {
    player_id,
    room_id
  }
});


socket.on("join-room-failed", (err) => {
  statusRoom.innerHTML = `<div class="alert alert-danger">${err}</div>`;
});

socket.on("game-ready", (msg) => {
  let countdown = 3;

  const status = document.createElement("h4");
  status.className = "text-success";
  status.innerHTML = "All Players Ready!";

  const countdownText = document.createElement("div");
  countdownText.className = "fs-4 mt-2 text-info fw-bold";
  countdownText.innerHTML = `Game starts in ${countdown}...`;

  statusRoom.innerHTML = "";
  statusRoom.appendChild(status);
  statusRoom.appendChild(countdownText);

  const countdownInterval = setInterval(() => {
    countdown--;
    if (countdown >= 0) {
      countdownText.innerHTML = `Game starts in ${countdown}...`;
    }

    if (countdown < 0) {
      clearInterval(countdownInterval);
      renderGameBoard();
      socket.emit("request-game-state");
    }
  }, 1000);
});



socket.on("game-not-ready", (msg) => {
  statusRoom.innerHTML = `<h4 class="text-warning">${msg}</h4>`;
});

// ✅ ฟังก์ชันสร้างกระดานเกม 3x3
function renderGameBoard() {
  boardElement = document.createElement("div");
  boardElement.className = "d-grid mx-auto mt-3";
  boardElement.style.gridTemplateColumns = "repeat(3, 100px)";
  boardElement.style.gap = "5px";
  boardElement.style.width = "max-content";

  for (let i = 0; i < 9; i++) {
    const cell = document.createElement("button");
    cell.className = "btn btn-outline-dark";
    cell.style.width = "100px";
    cell.style.height = "100px";
    cell.textContent = "";
    cell.dataset.index = i;

    // ✅ ผูกการคลิกให้ส่งตำแหน่ง
    cell.addEventListener("click", () => {
      const row = Math.floor(i / 3);
      const col = i % 3;
      socket.emit("make-move", { row, col });
    });

    boardElement.appendChild(cell);
  }

  statusRoom.innerHTML = "";
  statusRoom.appendChild(boardElement);
}

  socket.on("update-board", (gameRoom) => {
    if (!boardElement) return;

    MySymbol = gameRoom.player_x.player_id === player_id ? "X" : "O";

    for (let i = 0; i < 9; i++) {
      const row = Math.floor(i / 3);
      const col = i % 3;
      const cell = boardElement.children[i];
      const symbol = gameRoom.board[row][col];
      cell.textContent = symbol;
      cell.disabled = symbol !== "";
    }

    if (gameRoom.turn === player_id) {
      turnStatus.innerHTML = `<span class="text-success">👉 It's your turn! | Your Symbol : ${MySymbol}</span>`;
    } else {
      turnStatus.innerHTML = `<span class="text-secondary">⏳ Waiting for opponent...</span>`;
    }
  });



socket.on("make-move-failed", (msg) => {
  alert(msg);
});

// ✅ ฟัง event ผู้ชนะ
socket.on("game-winner", (winnerID) => {
  alert(`🏆 Player ${winnerID} is the Winner!`);
  disableBoard();
});

// ✅ ฟัง event เสมอ
socket.on("game-draw", (msg) => {
  alert("🤝 Draw! " + msg);
  disableBoard();
});

// ✅ ปิดการใช้งานกระดานทั้งหมดเมื่อจบเกม
function disableBoard() {
  if (!boardElement) return;
  Array.from(boardElement.children).forEach(cell => {
    cell.disabled = true;
  });
}



