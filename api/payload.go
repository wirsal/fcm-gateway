package api

import "github.com/wirsal/fcm-gateway/fcm"

type BroadcastPayload struct {
	Condition    string            `json:"condition" binding:"required"`
	Notification fcm.Notification  `json:"notification" binding:"required"`
	Data         map[string]string `json:"data,omitempty"`
	Android      fcm.AndroidConfig `json:"android,omitempty"`
	Apns         fcm.ApnsConfig    `json:"apns,omitempty"`
}

type RequestPayload struct {
	Tokens       []string          `json:"tokens" binding:"required"`
	Notification fcm.Notification  `json:"notification" binding:"required"`
	Android      fcm.AndroidConfig `json:"android,omitempty"`
	Apns         fcm.ApnsConfig    `json:"apns,omitempty"`
}
