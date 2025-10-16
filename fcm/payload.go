package fcm

type ApnsAps struct {
	MutableContent int `json:"mutable-content,omitempty"`
	// Alert          *ApnsAlert `json:"alert,omitempty"`
	Badge int    `json:"badge,omitempty"`
	Sound string `json:"sound,omitempty"`
}

type ApnsPayload struct {
	Aps ApnsAps `json:"aps"`
}

type ApnsConfig struct {
	Headers map[string]string `json:"headers,omitempty"`
	Payload ApnsPayload       `json:"payload"`
}

type AndroidConfig struct {
	Priority string `json:"priority,omitempty"`
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	// Sound string `json:"sound"`
	Image string `json:"image"`
}

type Message struct {
	Token        string        `json:"token"`
	Notification Notification  `json:"notification"`
	Android      AndroidConfig `json:"android,omitempty"`
	Apns         ApnsConfig    `json:"apns,omitempty"`
}

type FCMRequest struct {
	Message Message `json:"message"`
}

// ### Here all broadcast ###

type BroadcastMessage struct {
	Condition    string            `json:"condition"`
	Notification Notification      `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
	Android      AndroidConfig     `json:"android,omitempty"`
	Apns         ApnsConfig        `json:"apns,omitempty"`
}

type FCMBroadcastRequest struct {
	Message BroadcastMessage `json:"message"`
}
