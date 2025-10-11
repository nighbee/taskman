package handlers

import (
	"encoding/json"
	"log"
	"taskman-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	clients map[uuid.UUID]map[*websocket.Conn]bool // orgID -> connections
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients: make(map[uuid.UUID]map[*websocket.Conn]bool),
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *fiber.Ctx) error {
	// Get user ID from context (should be set by auth middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "User not authenticated"})
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Get organization ID from query parameter
	orgIDStr := c.Query("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	// Upgrade connection
	return websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		// Add client to organization
		h.addClient(orgID, c)

		// Send welcome message
		welcomeMsg := models.WebSocketMessage{
			Type:      models.MessageTypeUserJoined,
			Data:      []byte(`{"user_id":"` + userIDUUID.String() + `","org_id":"` + orgID.String() + `"}`),
			Timestamp: time.Now(),
			UserID:    userIDUUID,
		}
		h.sendMessageToOrg(orgID, welcomeMsg)

		// Handle messages
		for {
			var msg models.WebSocketMessage
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				h.removeClient(orgID, c)
				break
			}

			// Set user ID for the message
			msg.UserID = userIDUUID
			msg.Timestamp = time.Now()

			// Handle different message types
			switch msg.Type {
			case models.MessageTypeTaskMoved:
				h.handleTaskMoved(orgID, msg)
			case models.MessageTypeProjectMoved:
				h.handleProjectMoved(orgID, msg)
			default:
				log.Printf("Unknown message type: %s", msg.Type)
			}
		}
	})(c)
}

// addClient adds a client to an organization
func (h *WebSocketHandler) addClient(orgID uuid.UUID, conn *websocket.Conn) {
	if h.clients[orgID] == nil {
		h.clients[orgID] = make(map[*websocket.Conn]bool)
	}
	h.clients[orgID][conn] = true
}

// removeClient removes a client from an organization
func (h *WebSocketHandler) removeClient(orgID uuid.UUID, conn *websocket.Conn) {
	if h.clients[orgID] != nil {
		delete(h.clients[orgID], conn)
		if len(h.clients[orgID]) == 0 {
			delete(h.clients, orgID)
		}
	}
}

// sendMessageToOrg sends a message to all clients in an organization
func (h *WebSocketHandler) sendMessageToOrg(orgID uuid.UUID, msg models.WebSocketMessage) {
	if h.clients[orgID] != nil {
		for conn := range h.clients[orgID] {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
				conn.Close()
				delete(h.clients[orgID], conn)
			}
		}
	}
}

// BroadcastTaskMoved broadcasts a task moved event
func (h *WebSocketHandler) BroadcastTaskMoved(orgID uuid.UUID, taskID uuid.UUID, projectID uuid.UUID, oldStatus, newStatus string, userID uuid.UUID) {
	data := models.TaskMovedData{
		TaskID:    taskID,
		ProjectID: projectID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		UserID:    userID,
	}

	dataBytes, _ := json.Marshal(data)
	msg := models.WebSocketMessage{
		Type:      models.MessageTypeTaskMoved,
		Data:      dataBytes,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	h.sendMessageToOrg(orgID, msg)
}

// BroadcastProjectMoved broadcasts a project moved event
func (h *WebSocketHandler) BroadcastProjectMoved(orgID uuid.UUID, projectID uuid.UUID, oldStatus, newStatus string, userID uuid.UUID) {
	data := models.ProjectMovedData{
		ProjectID: projectID,
		OrgID:     orgID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		UserID:    userID,
	}

	dataBytes, _ := json.Marshal(data)
	msg := models.WebSocketMessage{
		Type:      models.MessageTypeProjectMoved,
		Data:      dataBytes,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	h.sendMessageToOrg(orgID, msg)
}

// handleTaskMoved handles task moved messages
func (h *WebSocketHandler) handleTaskMoved(orgID uuid.UUID, msg models.WebSocketMessage) {
	// Broadcast to all clients in the organization
	h.sendMessageToOrg(orgID, msg)
}

// handleProjectMoved handles project moved messages
func (h *WebSocketHandler) handleProjectMoved(orgID uuid.UUID, msg models.WebSocketMessage) {
	// Broadcast to all clients in the organization
	h.sendMessageToOrg(orgID, msg)
}
