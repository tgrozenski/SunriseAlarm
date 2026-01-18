package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	fcmEndpoint = "https://fcm.googleapis.com/v1/projects/sunrise-alarm/messages:send"
	maxRetries  = 3
	initialWait = 1 * time.Second
)

type Notifier interface {
	SendNotification(ctx context.Context, token string, deviceID string) error
}

type FCMNotifier struct {
	apiKey string
	client *http.Client
}

func NewFCMNotifier(apiKey string) *FCMNotifier {
	return &FCMNotifier{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (n *FCMNotifier) SendNotification(ctx context.Context, token, deviceID string) error {
	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"token": token,
			"data": map[string]string{
				"type":     "SUNRISE_ALARM",
				"deviceId": deviceID,
			},
			"android": map[string]string{
				"priority": "high",
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal FCM payload: %w", err)
	}

	var lastErr error
	wait := initialWait
	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, "POST", fcmEndpoint, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+n.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := n.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("FCM request failed: %w", err)
			time.Sleep(wait)
			wait *= 2
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		respBody, _ := io.ReadAll(resp.Body)
		lastErr = fmt.Errorf("FCM error %d: %s", resp.StatusCode, respBody)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			break
		}
		time.Sleep(wait)
		wait *= 2
	}
	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

type SecretManager interface {
	GetSecret(ctx context.Context, secretName string) (string, error)
}

type AWSSecretManager struct {
	client *secretsmanager.Client
}

func NewAWSSecretManager(client *secretsmanager.Client) *AWSSecretManager {
	return &AWSSecretManager{client: client}
}

func (s *AWSSecretManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	result, err := s.client.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret %s: %w", secretName, err)
	}
	if result.SecretString == nil {
		return "", fmt.Errorf("secret %s has no string value", secretName)
	}
	return *result.SecretString, nil
}
