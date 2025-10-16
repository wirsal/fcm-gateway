package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wirsal/fcm-gateway/fcm"
)

type Handler struct {
	fcmService *fcm.Service
}

func NewHandler(fcmService *fcm.Service) *Handler {
	return &Handler{fcmService: fcmService}
}

func (h *Handler) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome!",
	})
}

func (h *Handler) SendNotification(c *gin.Context) {
	var payload RequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body: " + err.Error()})
		return
	}

	if len(payload.Tokens) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tokens list cannot be empty"})
		return
	}

	successCount := 0
	failureCount := 0
	var failedTokens []map[string]string

	for _, token := range payload.Tokens {
		_, err := h.fcmService.SendNotification(token, payload.Notification, payload.Android.Priority, payload.Apns.Headers, payload.Apns.Payload)
		if err != nil {
			log.Printf("Gagal kirim ke token %s: %v", token, err)
			failureCount++
			failedTokens = append(failedTokens, map[string]string{"token": token, "error": err.Error()})
		} else {
			successCount++
		}
	}
	response := gin.H{
		"success_count": successCount,
		"failure_count": failureCount,
	}
	if failureCount > 0 {
		response["failed_tokens"] = failedTokens
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) SendBroadcast(c *gin.Context) {
	var payload BroadcastPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body: " + err.Error()})
		return
	}

	_, err := h.fcmService.BroadcastNotification(
		payload.Condition,
		payload.Notification,
		payload.Data,
		payload.Android.Priority,
		payload.Apns.Headers,
		payload.Apns.Payload,
	)
	if err != nil {
		log.Printf("Gagal broadcast ke condition %s: %v", payload.Condition, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send broadcast", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Broadcast message successfully sent to FCM for topic condition.",
	})
}
