# 🎮 Tic Tac Toe Online — Realtime Multiplayer

นี่คือโปรเจกต์เกม Tic Tac Toe ออนไลน์ที่มาพร้อมฟีเจอร์การเล่นแบบเรียลไทม์ ผู้เล่นสามารถใช้หมากพิเศษและสื่อสารกับระบบผ่าน **Socket.IO** และ API ที่ถูกออกแบบมาอย่างเรียบง่าย

---

## 📌 ฟีเจอร์หลัก

* **Realtime Multiplayer:** เล่นพร้อมกันแบบเรียลไทม์ผ่าน Socket.IO
* **หมากพิเศษ 3 ขนาด:** หมากมี 3 ขนาด ได้แก่ เล็ก (1), กลาง (2), และใหญ่ (3)
* **กฎพิเศษ "ทับหมาก":** หมากขนาดใหญ่กว่าสามารถวางทับหมากของฝ่ายตรงข้ามได้
* **Backend พลัง Go:** พัฒนาด้วย **Go (Gin + Socket.IO)** เพื่อประสิทธิภาพและความเร็ว
* **Frontend ที่ใช้งานง่าย:** สร้างด้วย HTML, Bootstrap และ JavaScript
* **API จัดการผู้เล่น:** มี API สำหรับสร้างและดึงข้อมูลผู้เล่น

---

## 📂 เทคโนโลยีที่ใช้

| Layer     | Tech Stack                     |
| :-------- | :----------------------------- |
| Backend   | Go, Gin, Socket.IO             |
| Frontend  | HTML, JavaScript, Bootstrap    |
| Protocol  | WebSocket (Socket.IO), HTTP API |

---

## 🧠 Logic ของเกม

* ผู้เล่นแต่ละคนมีหมาก 3 ขนาด:
    * **เล็ก** (ธรรมดา)
    * **กลาง** (symbol|medium|2)
    * **ใหญ่** (symbol|large|3)
* สามารถทับหมากของฝ่ายตรงข้ามได้ก็ต่อเมื่อหมากที่วางมีขนาดใหญ่กว่าเท่านั้น
* เกมจะจบลงเมื่อมีผู้เล่นวางหมากเรียงกัน 3 ตัว ไม่ว่าจะเป็นแนวตั้ง แนวนอน หรือแนวทแยง โดยยึด **symbol หลัก** ของหมากที่ชนะ

---

## 🌐 Socket.IO Events (Realtime)

### Namespace: `/` (Lobby)

| Event                  | ทิศทาง           | อธิบาย                                     |
| :--------------------- | :--------------- | :---------------------------------------- |
| `connect-successfuly`  | Server → Client  | แจ้งว่าเชื่อมต่อสำเร็จ                      |
| `create-room`          | Client → Server  | สร้างห้องเกมใหม่                           |
| `create-room-successfuly` | Server → Client  | แจ้งผู้เล่นว่าสร้างห้องเรียบร้อยแล้ว        |
| `connection`           | Server → Client  | ส่งข้อมูลห้องทั้งหมดให้ Lobby               |

### Namespace: `/game-room` (In Game)

| Event                  | ทิศทาง            | อธิบาย                                                 |
| :--------------------- | :---------------- | :------------------------------------------------------ |
| `connect-successfuly`  | Server → Client   | ผู้เล่นเชื่อมต่อเข้าห้องเกมได้แล้ว                       |
| `join-room-failed`     | Server → Client   | เข้าห้องเกมล้มเหลว                                       |
| `game-ready`           | Server → Room     | เกมพร้อมเริ่ม (มีผู้เล่นครบ)                               |
| `game-not-ready`       | Server → Room     | รอผู้เล่นอีกคน                                           |
| `make-move`            | Client → Server   | ส่งข้อมูลการเดินหมาก                                    |
| `update-board`         | Server → Room     | กระจายข้อมูลกระดานใหม่หลังมีการเดินหมาก                 |
| `game-winner`          | Server → Room     | มีผู้ชนะ                                                |
| `game-draw`            | Server → Room     | เกมเสมอ                                                |
| `request-game-state`   | Client → Server   | ดึงสถานะกระดานล่าสุด                                     |

---

## 🔄 REST API

| Method | Endpoint          | Description            |
| :----- | :---------------- | :--------------------- |
| `POST` | `/api/player`     | สร้างชื่อผู้เล่นใหม่    |
| `GET`  | `/api/player/:id` | ดึงข้อมูลผู้เล่นตาม ID |

---

## 🚀 วิธี Deploy โปรเจกต์ Tic-Tac-Toe Online ขึ้น Render (Web Service)

โปรเจกต์นี้สามารถ Deploy ได้ง่าย ๆ บน Render โดยใช้เพียง **Web Service** ของ Render เท่านั้น เนื่องจาก Backend ที่พัฒนาด้วย Go มีการ Serve ไฟล์จาก `public/` อยู่แล้ว จึงไม่จำเป็นต้องใช้ Static Site เพิ่มเติม

### 🛠️ Deploy Backend (Go)

1.  ไปที่ [https://render.com](https://render.com) และเข้าสู่ระบบด้วย GitHub
2.  คลิก `New` → `Web Service`
3.  เลือก Repository ของคุณที่ Push โปรเจกต์นี้ขึ้นไป
4.  ตั้งค่าตามรายละเอียดด้านล่าง:
    * **Name:** `tic-tac-toe-online` (หรือชื่อที่คุณต้องการ)
    * **Runtime:** `Go`
    * **Build Command:** `go build -o main`
    * **Start Command:** `./main`
5.  คลิก `Create Web Service`

### 🛠️ ตั้งค่าหลัง Deploy

1.  หลังจาก Deploy สำเร็จ ให้คัดลอก **URL** ที่ได้จาก Render
2.  เปิดไฟล์ `index.js`, `game_room.js`, และ `lobby.js` ในโปรเจกต์ของคุณ
3.  แทนที่ `const HTTP_SERVER_URL = "...";` ด้วย URL ที่คุณคัดลอกมา
4.  จากนั้น **Push การเปลี่ยนแปลงเหล่านี้ขึ้น Git** และรอให้ Render ทำการ Deploy ใหม่

เมื่อการ Deploy เสร็จสมบูรณ์ คุณก็จะสามารถเข้าถึงเกมผ่าน URL ที่ได้รับและเชื่อมต่อกับ Backend ได้ทันที!

---

## ❤️ Created with by Surachat
