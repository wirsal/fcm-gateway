package fcm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// mockRoundTripper adalah custom http.RoundTripper untuk mocking panggilan HTTP.
type mockRoundTripper struct {
	handlers map[string]http.HandlerFunc
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Temukan handler yang cocok berdasarkan URL
	for urlPrefix, handler := range m.handlers {
		if strings.HasPrefix(req.URL.String(), urlPrefix) {
			rec := httptest.NewRecorder()
			handler(rec, req)
			return rec.Result(), nil
		}
	}
	// Jika tidak ada handler yang cocok, kembalikan error
	return nil, fmt.Errorf("no mock handler for %s", req.URL.String())
}

// createTestCredentialsFile membuat file kredensial JSON palsu untuk pengujian.
func createTestCredentialsFile(t *testing.T, tokenURL, projectID string) string {
	t.Helper()
	credsContent := fmt.Sprintf(`{
		"type": "service_account",
		"project_id": "%s",
		"private_key_id": "test_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC/p9J+t3T4\ndyT/O6EHRgACgQEA7A==\n-----END PRIVATE KEY-----\n",
		"client_email": "test@test-project.iam.gserviceaccount.com",
		"client_id": "123456789",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "%s",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/test-project.iam.gserviceaccount.com"
	}`, projectID, tokenURL)

	tmpFile, err := os.CreateTemp(t.TempDir(), "test-creds-*.json")
	require.NoError(t, err)

	_, err = tmpFile.WriteString(credsContent)
	require.NoError(t, err)

	err = tmpFile.Close()
	require.NoError(t, err)

	return tmpFile.Name()
}

func TestNewService(t *testing.T) {
	t.Run("success - create new service", func(t *testing.T) {
		// --- Setup ---
		// Gunakan mock HTTP client untuk mencegah panggilan keluar saat mengambil token
		mockClient := &http.Client{
			Transport: &mockRoundTripper{
				handlers: map[string]http.HandlerFunc{
					"https://oauth2.googleapis.com/token": func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
					},
				},
			},
		}
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)
		credsFile := createTestCredentialsFile(t, "https://oauth2.googleapis.com/token", "test-project-id")

		// --- Execute ---
		service, err := NewService(ctx, credsFile, []string{"test-scope"}, "http://test.fcm/endpoint")

		// --- Assert ---
		require.NoError(t, err)
		require.NotNil(t, service)
		assert.Equal(t, "test-project-id", service.projectID)
		assert.Equal(t, "http://test.fcm/endpoint", service.endpointURL)
	})

	t.Run("error - credentials file not found", func(t *testing.T) {
		// --- Execute ---
		_, err := NewService(context.Background(), "non-existent-file.json", nil, "")

		// --- Assert ---
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error read credentials file")
	})

	t.Run("error - invalid json in credentials file", func(t *testing.T) {
		// --- Setup ---
		tmpFile, err := os.CreateTemp(t.TempDir(), "bad-creds-*.json")
		require.NoError(t, err)
		_, err = tmpFile.WriteString("this is not json")
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// --- Execute ---
		_, err = NewService(context.Background(), tmpFile.Name(), nil, "")

		// --- Assert ---
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gagal load credentials")
	})
}

func TestService_SendNotification(t *testing.T) {
	projectID := "test-project-123"
	fcmEndpoint := "https://fcm.googleapis.com/v1/projects/"
	tokenEndpoint := "https://oauth2.googleapis.com/token"

	t.Run("success - send notification", func(t *testing.T) {
		// --- Setup ---
		// Buat mock transport untuk menangani permintaan token dan FCM
		mockTransport := &mockRoundTripper{
			handlers: map[string]http.HandlerFunc{
				tokenEndpoint: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"access_token":"mock-access-token","token_type":"Bearer","expires_in":3600}`))
				},
				fcmEndpoint: func(w http.ResponseWriter, r *http.Request) {
					// Verifikasi header authorization
					authHeader := r.Header.Get("Authorization")
					assert.Equal(t, "Bearer mock-access-token", authHeader)

					// Verifikasi body request
					bodyBytes, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					var reqBody FCMRequest
					err = json.Unmarshal(bodyBytes, &reqBody)
					require.NoError(t, err)

					assert.Equal(t, "test-device-token", reqBody.Message.Token)
					assert.Equal(t, "Test Title", reqBody.Message.Notification.Title)

					// Kembalikan respons sukses
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"name":"projects/test-project/messages/12345"}`))
				},
			},
		}

		mockClient := &http.Client{Transport: mockTransport}
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)
		credsFile := createTestCredentialsFile(t, tokenEndpoint, projectID)

		// Buat service dengan kredensial dan URL endpoint yang sudah di-mock
		data, err := os.ReadFile(credsFile)
		require.NoError(t, err)
		// FIX: Scopes harus diberikan ke google.CredentialsFromJSON untuk mencegah
		// penggunaan logika kredensial default yang mengabaikan mock client di context.
		creds, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/firebase.messaging")
		require.NoError(t, err)

		service := &Service{
			creds:       creds,
			projectID:   projectID,
			endpointURL: fcmEndpoint + "%s/messages:send", // Format URL sesuai implementasi
		}

		// --- Execute ---
		notification := Notification{Title: "Test Title", Body: "Test Body"}
		apnsPayload := ApnsPayload{}
		respBody, err := service.SendNotification("test-device-token", notification, "high", nil, apnsPayload)

		// --- Assert ---
		require.NoError(t, err)
		assert.Contains(t, respBody, "projects/test-project/messages/12345")
	})

	t.Run("error - fcm returns non-200 status", func(t *testing.T) {
		// --- Setup ---
		mockTransport := &mockRoundTripper{
			handlers: map[string]http.HandlerFunc{
				tokenEndpoint: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"access_token":"mock-access-token","token_type":"Bearer","expires_in":3600}`))
				},
				fcmEndpoint: func(w http.ResponseWriter, r *http.Request) {
					// Kembalikan respons error
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"internal server error"}`))
				},
			},
		}
		mockClient := &http.Client{Transport: mockTransport}
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)
		credsFile := createTestCredentialsFile(t, tokenEndpoint, projectID)
		data, err := os.ReadFile(credsFile)
		require.NoError(t, err)

		creds, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/firebase.messaging")
		require.NoError(t, err)

		service := &Service{
			creds:       creds,
			projectID:   projectID,
			endpointURL: fcmEndpoint + "%s/messages:send",
		}

		// --- Execute ---
		_, err = service.SendNotification("any-token", Notification{}, "", nil, ApnsPayload{})

		// --- Assert ---
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FCM error 500")
		assert.Contains(t, err.Error(), "internal server error")
	})
}
