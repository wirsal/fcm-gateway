package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

type Service struct {
	creds       *google.Credentials
	projectID   string
	endpointURL string
	httpClient  *http.Client
}

func NewService(ctx context.Context, credentialsFile string, scopes []string, endpointURL string) (*Service, error) {
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("error read credentials file")
	}

	creds, err := google.CredentialsFromJSON(ctx, data, scopes...)
	if err != nil {
		return nil, fmt.Errorf("gagal load credentials: %w", err)
	}

	return &Service{
		creds:       creds,
		projectID:   creds.ProjectID,
		endpointURL: endpointURL,
		httpClient:  &http.Client{},
	}, nil
}

func (s *Service) SendNotification(
	token string,
	notification Notification,
	androidPriority string,
	apnsHeaders map[string]string,
	apnsPayload ApnsPayload,
) (string, error) {
	reqBody := FCMRequest{
		Message: Message{
			Token:        token,
			Notification: notification,
			Android:      AndroidConfig{Priority: androidPriority},
			Apns: ApnsConfig{
				Headers: apnsHeaders,
				Payload: apnsPayload},
		},
	}
	return sendToFirebase(s, reqBody)
}

func (s *Service) BroadcastNotification(
	condition string,
	notification Notification,
	data map[string]string, // Argument baru
	androidPriority string,
	apnsHeaders map[string]string,
	apnsPayload ApnsPayload,
) (string, error) {
	reqBody := FCMBroadcastRequest{ // Menggunakan struct request yang baru
		Message: BroadcastMessage{
			Condition:    condition,
			Notification: notification,
			Data:         data, // Sertakan Data
			Android:      AndroidConfig{Priority: androidPriority},
			Apns: ApnsConfig{
				Headers: apnsHeaders,
				Payload: apnsPayload,
			},
		},
	}

	return sendToFirebase(s, reqBody)
}

func sendToFirebase(s *Service, reqBody interface{}) (string, error) {
	tok, err := s.creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("create token failed")
	}
	url := fmt.Sprintf(s.endpointURL, s.projectID)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("gagal marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error kirim request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal baca response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("FCM error %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
