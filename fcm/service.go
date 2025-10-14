package fcm

import "golang.org/x/oauth2/google"

type AndroidConfig struct {
	Priority string `json:"priority,omitempty"`
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Message struct {
	Token        string        `json:"token"`
	Notification Notification  `json:"notification"`
	Android      AndroidConfig `json:"android,omitempty"`
}

type FCMRequest struct {
	Message Message `json:"message"`
}

type Service struct {
	creds     *google.Credentials
	projectID string
}
