const http = "http://localhost:3000";
const playerData = localStorage.getItem("player_id");
const socket = io(http, {
  query: {
    player_id: playerData  // ต้องกำหนดไว้ก่อน
  }
});


    const messageBox = document.getElementById("message");
    const roomList = document.getElementById("room-list");

    // แสดงข้อความในกล่อง message
    // function showMessage(msg, type = 'info') {
    //   messageBox.textContent = msg;
    //   messageBox.className = `alert alert-${type} text-center`;
    //   messageBox.classList.remove("d-none");
    //   setTimeout(() => messageBox.classList.add("d-none"), 3000);
    // }

    // ฟัง event การเชื่อมต่อ
    // socket.on("connect-successfuly", (msg) => {
    //   showMessage(msg, "success");
    // });

    // ฟัง event แสดงรายชื่อห้องทั้งหมด
    socket.on("connection", (rooms) => {
      GetPlayerData()
      updateRoomList(rooms);
    });

    // ฟัง event หลังจากสร้างห้องสำเร็จ
    socket.on("create-room-successfuly", (data) => {
      // showMessage("Room created: " + newRoom.room_id, "success");
      console.log("Room created:", data.room.room_id, "by player:", data.player_id);
      addRoomToList(data.room);
      if (data.player_id == playerData) {
        joinRoom(data.room.room_id);
      }
    });

    // ฟังก์ชันสร้างห้อง
    function createRoom() {
      socket.emit("create-room");
    }

    // อัปเดตรายชื่อห้องทั้งหมด
    function updateRoomList(rooms) {
      roomList.innerHTML = '';
      rooms.forEach(room => {
        addRoomToList(room);
      });
    }

    // เพิ่มห้องหนึ่งห้องลงใน list
   function addRoomToList(room) {
  const li = document.createElement("li");
  li.className = "list-group-item d-flex justify-content-between align-items-center";

  const isGameStarted = room.player_x?.player_name && room.player_o?.player_name;

  li.innerHTML = `
    Room ID: <strong>${room.room_id}</strong>
    <span class="badge bg-primary rounded-pill">
      Players: ${room.player_x?.player_name || 'Waiting'} vs ${room.player_o?.player_name || 'Waiting'}
    </span>
    ${isGameStarted
      ? `<span class="badge bg-danger rounded-pill">Game Started</span>`
      : `<button class="btn btn-primary" onclick="joinRoom('${room.room_id}')">Join</button>`}
  `;

  roomList.appendChild(li);
}


    // ฟังก์ชันเชื่อมต่อห้อง
    function joinRoom(roomId) {
      localStorage.setItem("room_id", roomId);
      window.location.href = "game_room.html";
    }

    function GetPlayerData() {
      const playerData = localStorage.getItem("player_id");
      if (!playerData) {
        window.location.href = "index.html";
        return;
      }
      fetch(http +"/api/player/" + playerData).then(res => {
    if (!res.ok) throw new Error("Bad response");
    return res.json();
    }).then(data => {
        document.getElementById("player-name").textContent = data.player.player_name;
        document.getElementById("player-id").textContent = data.player.player_id;
        document.getElementById("player-won").textContent = data.player.win;
        document.getElementById("player-lost").textContent = data.player.lose;
        document.getElementById("player-draw").textContent = data.player.draw;
      })
      .catch(err => {
      window.location.href = "index.html";
      });
    }