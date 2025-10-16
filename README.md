# 🚀 FCM Gateway: High-Performance Notification Service

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/)

FCM Gateway is a high-performance RESTful API built with Go and Gin, serving as a reliable gateway to send push notifications through the Firebase Cloud Messaging (FCM) v1 API.This project is designed to be a robust, efficient, and easily configurable solution by implementing direct HTTP calls to the FCM endpoint to ensure maximum reliability.✨ Key FeaturesHigh Performance: Built on Gin, one of the fastest web frameworks in the Go ecosystem.Reliable: Uses direct HTTP calls to the FCM API v1, avoiding potential issues or bugs that might exist in the SDKs.Efficient Authentication: Manages the OAuth2 token lifecycle automatically, including token caching and refreshing, thereby minimizing latency.External Configuration: All settings (port, scopes, file paths, endpoints) are fully managed by Viper through a config.yaml file, with no hardcoded values.Flexible: Supports sending notifications to one or many device tokens in a single API call.Dynamic Payloads: Allows clients to define platform-specific configurations like Android priority or APNS headers directly from the JSON payload.Clean Project Structure: The code is neatly organized by functionality (api, fcm, config, credentials) for easy maintenance.

## 📁 Project Structure

```text
fcm-gateway/
├── config/                  # For .yaml data files only
│   └── .config.yaml
├── credentials/             # For secret files ONLY
│   └── service-account.json
├── api/                     # HTTP Handlers (Gin)
│   └── handler.go
├── fcm/                     # Core logic for interacting with FCM
│   └── service.go
├── internal/                # Internal project packages
│   └── config/              # Go logic for loading configuration
│       └── config.go
├── cmd
|   └── main.go              # Application entry point
├── go.mod
└── .gitignore               # Important for security
```

### 🚀 Getting Started
**Prerequisites** 
- Go (version 1.22 or newer)
- A Google Cloud Project with Firebase enabled.
- Billing linked to your Google Cloud project.
- The Firebase Cloud Messaging API (v1) enabled in your Google Cloud Console.

***Installation and Configuration***

Follow these steps to run the server in your local environment.
1. Clone the Repository
```bash
git clone https://github.com/wirsal/fcm-gateway.git
cd fcm-gateway
```

2. Prepare the Configuration FileThis application uses a `.yaml` file for configuration.# Copy the sample configuration file
```bash 
cp config/sample.config.yaml config/.config.yaml
```

You can adjust the port number inside `config/config.yaml` if needed.

3. Prepare the Credentials File (Very Important)
The server requires a Service Account key in JSON format to authenticate with Google.

  a. Get your `service-account.json` file from the Google Cloud Console:
  Go to the Service Accounts page: `https://console.cloud.google.com/iam-admin/serviceaccounts` 
  Select the correct project.
  Choose the firebase-adminsdk service account.
  Go to the KEYS tab, click ADD KEY > Create new key, select JSON, and then click CREATE. 
  A file will be downloaded.

  b. Move and rename the downloaded file:# Copy the sample credentials file
cp credentials/sample.service-account.json credentials/service-account.json

  c. Open the newly copied `credentials/service-account.json`, delete its contents, and paste the entire content from the JSON file you downloaded from Google Cloud.
  
4. Install DependenciesThis command will automatically download all required libraries (Gin, Viper, etc.).go mod tidy

**Running the Server**

Once all configurations are complete, run the server with the command:go run main.go

If successful, you will see the following output:
```bash
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
    
    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
     - using env:   export GIN_MODE=release
     - using code:  gin.SetMode(gin.ReleaseMode)
    
    [GIN-debug] POST   /send                     --> fcm-gateway/api.(*Handler).SendNotifications (3 handlers)
    [GIN-debug] Listening and serving HTTP on :8080
    Server Gin is running at http://localhost:8080
```
## ⚙️ API Usage
**Endpoint:** 
```curl
POST /sendSends a notification to one or more devices.curl Example:curl -X POST http://localhost:8080/send \
-H "Content-Type: application/json" \
-d '{
    "tokens": [
        "YOUR_DEVICE_TOKEN_HERE"
    ],
    "notification": {
        "title": "Hello from FCM Gateway!",
        "body": "This notification was sent via a cool Go API."
    },
    "android": {
        "priority": "HIGH"
    },
    "apns": {
        "headers": {
            "apns-priority": "10"
        }
    }
}'
```
```json
Request Body| Key | Type | Required? | Description || tokens | []string | Yes | An array containing one or more device registration tokens. || notification | object | Yes | The object containing the title and body of the notification. || android | object | No | Android-specific configuration. Example: {"priority": "HIGH"}. || apns | object | No | APNS (iOS)-specific configuration. Example: {"headers": {"apns-priority": "10"}}. |Response ExamplesSuccess:{
    "failure_count": 0,
    "success_count": 1
}

Partial Failure:{
    "failure_count": 1,
    "success_count": 1,
    "failed_tokens": [
        {
            "error": "FCM error 400: The registration token is not a valid FCM registration token",
            "token": "AN_INVALID_TOKEN"
        }
    ]
}
```

## 📄 License  
This project is licensed under the MIT License. See the LICENSE file for details.

