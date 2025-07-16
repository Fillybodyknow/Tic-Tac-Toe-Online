// Your existing JavaScript code, ensuring IDs match
        var HTTP_SERVER_URL = "https://tic-tac-toe-online-backend-20dy.onrender.com";
        var submitbtn = document.getElementById("submit-button");
        var player_name_input = document.getElementById("player-name-input"); 

        submitbtn.addEventListener("click", (event) => {
            event.preventDefault();

            if (player_name_input.value.trim() === "") {
                alert("Please enter your name to play!");
                return;
            }

            fetch(HTTP_SERVER_URL + "/api/player", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "ngrok-skip-browser-warning": "true"
                },
                body: JSON.stringify({
                    username: player_name_input.value.trim()
                })
            })
            .then(res => {
                if (res.status === 204) {
                    return {};
                }
                return res.json();
            })
            .then(data => {
                if (data.error) {
                    alert(data.error);
                } else {
                    localStorage.setItem("player_id", data.player.player_id);
                    console.log("Player ID:", data.player.player_id);
                    
                    window.location.href = "lobby.html";
                }
            })
            .catch(err => {
                console.error("Request failed:", err);
                alert("An error occurred while connecting to the server. Please try again.");
            });
        });