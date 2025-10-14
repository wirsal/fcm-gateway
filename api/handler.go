package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wirsal/fcm-gateway/fcm"
)

type RequestPayload struct {
	Tokens       []string          `json:"tokens" binding:"required"`
	Notification fcm.Notification  `json:"notification" binding:"required"`
	Android      fcm.AndroidConfig `json:"android,omitempty"`
	Apns         fcm.ApnsConfig    `json:"apns,omitempty"`
}

type Handler struct {
	fcmService *fcm.Service
}

func NewHandler(fcmService *fcm.Service) *Handler {
	return &Handler{fcmService: fcmService}
}

func (h *Handler) Welcome(c *gin.Context) {

}

func (h *Handler) SendNotification(c *gin.Context) {
	var payload RequestPayload
	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body" + err.Error()})
		return
	}

	successCount := 0
	failureCount := 0
	var failedTokens []map[string]string

	for _, token := range payload.Tokens {
		_, err := h.fcmService.SendNotification(token, payload.Notification, payload.Android.Priority, payload.Apns.Headers)
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
