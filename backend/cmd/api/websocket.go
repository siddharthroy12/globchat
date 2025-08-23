package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/coder/websocket"
)

type WebsocketConnectionMessage struct {
	Type   string `json:"type"`
	RoomID int    `json:"room_id"`
	Data   any    `json:"data"`
}

type WebSocketRoomManager struct {
	rooms map[int]map[*websocket.Conn]bool // room_id -> set of connections
	mu    sync.RWMutex
}

func NewWebSocketRoomManager() *WebSocketRoomManager {
	return &WebSocketRoomManager{
		rooms: make(map[int]map[*websocket.Conn]bool),
	}
}

func (m *WebSocketRoomManager) joinRoom(c *websocket.Conn, roomId int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.rooms[roomId] == nil {
		m.rooms[roomId] = make(map[*websocket.Conn]bool)
	}

	m.rooms[roomId][c] = true
}

func (m *WebSocketRoomManager) leaveRoom(c *websocket.Conn, roomId int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if room, exists := m.rooms[roomId]; exists {
		delete(room, c)

		// Clean up empty rooms
		if len(room) == 0 {
			delete(m.rooms, roomId)
		}
	}
}

func (m *WebSocketRoomManager) notifyRoom(roomId int, message WebsocketConnectionMessage) {
	m.mu.RLock()
	room, exists := m.rooms[roomId]
	if !exists {
		m.mu.RUnlock()
		return
	}

	// Create a copy of connections to avoid holding the lock while sending
	connections := make([]*websocket.Conn, 0, len(room))
	for conn := range room {
		connections = append(connections, conn)
	}
	m.mu.RUnlock()

	js, err := json.Marshal(message)

	if err != nil {
		return
	}

	// Send message to all connections in the room (except sender)
	for _, conn := range connections {
		err := conn.Write(context.Background(), websocket.MessageText, js)
		if err != nil {
			m.leaveRoom(conn, roomId)
		}
	}
}
func (m *WebSocketRoomManager) leaveAllRooms(c *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for roomId, room := range m.rooms {
		if room[c] {
			delete(room, c)
			// Clean up empty rooms
			if len(room) == 0 {
				delete(m.rooms, roomId)
			}
		}
	}
}

func (app *application) websocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	println("got hit!")
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"wompwomp"},
	})

	if err != nil {
		app.serverErrorResponse(w, r, err, "websocket connection accept")
		return
	}

	defer func() {
		app.roomManager.leaveAllRooms(c)
		c.CloseNow()
	}()

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)

	for {
		err := app.handleMessage(c, l)

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			app.logError(r, err, "websocket message handling")
			return
		}
	}
}

func (app *application) handleMessage(c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*5)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, data, err := c.Read(ctx)
	if err != nil {
		return err
	}

	if typ != websocket.MessageText {
		return fmt.Errorf("expected text message, got %v", typ)
	}

	var msg WebsocketConnectionMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("invalid message format: %w", err)
	}

	switch msg.Type {
	case "join":
		app.roomManager.joinRoom(c, msg.RoomID)
	case "leave":
		app.roomManager.leaveRoom(c, msg.RoomID)

	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}

	return nil
}
