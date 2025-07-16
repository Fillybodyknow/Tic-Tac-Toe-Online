# 🎮 Tic Tac Toe Online — Realtime Multiplayer

โปรเจกต์เกม Tic Tac Toe แบบออนไลน์ ที่สามารถเล่นพร้อมกันแบบ Real-Time ได้ 2 คน โดยรองรับฟีเจอร์หมากพิเศษ และสื่อสารผ่าน **Socket.IO** กับระบบ API แบบเรียบง่าย

---

## 📌 ฟีเจอร์หลัก

- ✅ เล่นแบบ Realtime Multiplayer (ผ่าน Socket.IO)
- 🟢 หมากพิเศษ 3 ขนาด: เล็ก (1), กลาง (2), ใหญ่ (3)
- ⚔️ กฎพิเศษ: หมากขนาดใหญ่กว่าสามารถทับของศัตรูได้
- 🧠 Backend พัฒนาโดยใช้ **Go (Gin + Socket.IO)**
- 🎨 Frontend ใช้ HTML + Bootstrap + JavaScript
- 📦 API สำหรับสร้างและดึงข้อมูลผู้เล่น

---

## 📂 เทคโนโลยีที่ใช้

| Layer     | Tech Stack                     |
|-----------|--------------------------------|
| Backend   | Go, Gin, Socket.IO             |
| Frontend  | HTML, JavaScript, Bootstrap    |
| Protocol  | WebSocket (Socket.IO), HTTP API|

---

## 🧠 Logic ของเกม

- ผู้เล่นแต่ละคนมีหมากพิเศษชนิด:  
  - `เล็ก` (ธรรมดา)  
  - `กลาง` (`symbol|medium|2`)  
  - `ใหญ่` (`symbol|large|3`)  
- สามารถทับหมากของฝ่ายตรงข้ามได้ถ้าขนาดใหญ่กว่าเท่านั้น
- เกมชนะเมื่อวางหมากเรียง 3 แถวแนวตั้ง แนวนอน หรือแนวทแยง โดยพิจารณา **symbol หลัก**

---

## 🌐 Socket.IO Events (Realtime)

### Namespace: `/` (Lobby)
| Event                  | ทิศทาง       | อธิบาย                             |
|------------------------|--------------|-------------------------------------|
| `connect-successfuly` | Server → Client | แจ้งว่าเชื่อมต่อสำเร็จ              |
| `create-room`         | Client → Server | สร้างห้องเกมใหม่                    |
| `create-room-successfuly` | Server → Client | แจ้งผู้เล่นว่าได้สร้างห้องเรียบร้อย |
| `connection`          | Server → Client | ส่งข้อมูลห้องทั้งหมดให้ Lobby      |

### Namespace: `/game-room` (In Game)
| Event                | ทิศทาง         | อธิบาย                                     |
|----------------------|----------------|---------------------------------------------|
| `connect-successfuly` | Server → Client | ผู้เล่นเชื่อมต่อเข้าห้องเกมได้แล้ว         |
| `join-room-failed`   | Server → Client | เข้าห้องเกมล้มเหลว                          |
| `game-ready`         | Server → Room   | เกมพร้อมเริ่ม (มีผู้เล่นครบ)               |
| `game-not-ready`     | Server → Room   | รอผู้เล่นอีกคน                              |
| `make-move`          | Client → Server | ส่งข้อมูลการเดินหมาก                       |
| `update-board`       | Server → Room   | กระจายข้อมูลกระดานใหม่หลังมีการเดินหมาก   |
| `game-winner`        | Server → Room   | มีผู้ชนะ                                     |
| `game-draw`          | Server → Room   | เกมเสมอ                                      |
| `request-game-state` | Client → Server | ดึงสถานะกระดานล่าสุด                       |

---

## 🔄 REST API

| Method | Endpoint           | Description               |
|--------|--------------------|---------------------------|
| POST   | `/api/player`      | สร้างชื่อผู้เล่นใหม่       |
| GET    | `/api/player/:id`  | ดึงข้อมูลผู้เล่นตาม ID    |

---

## ▶️ วิธีใช้งาน (Local)

### 1. Clone โปรเจกต์

```bash
git clone https://github.com/yourusername/tic-tac-toe-online.git
cd tic-tac-toe-online
```
