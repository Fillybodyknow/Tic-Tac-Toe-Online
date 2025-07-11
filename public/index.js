var http = "https://wren-super-cobra.ngrok-free.app";

var submitbtn = document.getElementById("submit-button");
var player_name = document.getElementById("player-name");

submitbtn.addEventListener("click", () => {
    if (player_name.value == "") {
        alert("Please enter your name");
        return;
    }

    fetch(http + "/api/player", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            username: player_name.value
        })
    })
    .then(res => res.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            localStorage.setItem("player_id", data.player.player_id);  // ✅ FIX ตรงนี้
            window.location.href = "lobby.html";
        }
    })
    .catch(err => {
        console.error("Request failed:", err);
        alert("เกิดข้อผิดพลาดขณะส่งข้อมูลไปยังเซิร์ฟเวอร์");
    });
});
