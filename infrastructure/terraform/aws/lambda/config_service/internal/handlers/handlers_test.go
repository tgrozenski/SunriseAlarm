package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	internaldynamodb "myproject/internal/dynamodb"
	"myproject/internal/models"
	"myproject/internal/notification"
)

type mockNotifier struct {
	calls []struct {
		token    string
		deviceID string
	}
	err error
}

func (m *mockNotifier) SendNotification(ctx context.Context, token, deviceID string) error {
	m.calls = append(m.calls, struct {
		token    string
		deviceID string
	}{token, deviceID})
	return m.err
}

type mockSecretManager struct {
	secret string
	err    error
}

func (m *mockSecretManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.secret, nil
}

type mockStore struct {
	configs      map[string]*models.UserConfig
	err          error
	queryResults []models.UserConfig
	queryErr     error
}

func (m *mockStore) GetConfig(ctx context.Context, deviceID string) (*models.UserConfig, error) {
	if m.err != nil {
		return nil, m.err
	}
	config, ok := m.configs[deviceID]
	if !ok {
		return nil, internaldynamodb.ErrConfigNotFound
	}
	return config, nil
}

func (m *mockStore) PutConfig(ctx context.Context, config *models.UserConfig) error {
	if m.err != nil {
		return m.err
	}
	m.configs[config.DeviceID] = config
	return nil
}

func (m *mockStore) QueryByAlarmTime(ctx context.Context, dateBucket string, startTime, endTime string) ([]models.UserConfig, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return m.queryResults, nil
}

func newTestServer(store internaldynamodb.ConfigStore) http.Handler {
	router := mux.NewRouter()
	handler := NewConfigHandler(store)
	router.HandleFunc("/config/{deviceId}", handler.GetConfig).Methods("GET")
	router.HandleFunc("/config", handler.PutConfig).Methods("POST")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	return router
}

func newAlarmTestServer(store internaldynamodb.ConfigStore, notifier notification.Notifier, secretManager notification.SecretManager) http.Handler {
	router := mux.NewRouter()
	handler := NewAlarmHandler(store, notifier, secretManager)
	router.HandleFunc("/check_alarm", handler.CheckAlarm).Methods("GET")
	return router
}

func TestGetConfigEndpoint(t *testing.T) {
	store := &mockStore{
		configs: map[string]*models.UserConfig{
			"550e8400-e29b-41d4-a716-446655440002": {
				DeviceID:       "550e8400-e29b-41d4-a716-446655440002",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Offset:         0,
				Enabled:        true,
			},
		},
	}
	srv := newTestServer(store) //Stevie Ray Vaughan

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/config/550e8400-e29b-41d4-a716-446655440002", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var config models.UserConfig
		if err := json.Unmarshal(w.Body.Bytes(), &config); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}
		if config.DeviceID != "550e8400-e29b-41d4-a716-446655440002" {
			t.Errorf("expected deviceId 550e8400-e29b-41d4-a716-446655440002, got %s", config.DeviceID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/config/not-found", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to decode error response: %v", err)
		}
		if resp["error"] != "config not found" {
			t.Errorf("expected error 'config not found', got %s", resp["error"])
		}
	})

	t.Run("missing deviceId", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/config/", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestPutConfigEndpoint(t *testing.T) {
	store := &mockStore{
		configs: make(map[string]*models.UserConfig),
	}
	srv := newTestServer(store)

	t.Run("success", func(t *testing.T) {
		config := models.UserConfig{
			DeviceID:       "550e8400-e29b-41d4-a716-446655440001",
			Lat:            34.0522,
			Long:           -118.2437,
			FCMToken:       "token",
			DayPreferences: []bool{false, false, false, true, true, true, true},
			TimeZone:       "America/Los_Angeles",
			Offset:         0,
			Enabled:        true,
		}
		body, _ := json.Marshal(config)
		req := httptest.NewRequest("POST", "/config", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		stored, err := store.GetConfig(context.Background(), "550e8400-e29b-41d4-a716-446655440001")
		if err != nil {
			t.Errorf("config not stored: %v", err)
		}
		if stored.DeviceID != "550e8400-e29b-41d4-a716-446655440001" {
			t.Errorf("stored config mismatch")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/config", bytes.NewReader([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		config := models.UserConfig{
			DeviceID:       "",
			Lat:            34.0522,
			Long:           -118.2437,
			FCMToken:       "token",
			DayPreferences: []bool{false, false, false, true, true, true, true},
			TimeZone:       "America/Los_Angeles",
			Offset:         0,
			Enabled:        true,
		}
		body, _ := json.Marshal(config)
		req := httptest.NewRequest("POST", "/config", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestCheckAlarmEndpoint(t *testing.T) {
	t.Run("no alarms", func(t *testing.T) {
		store := &mockStore{
			queryResults: []models.UserConfig{},
		}
		notifier := &mockNotifier{}
		secretManager := &mockSecretManager{secret: "fake-key"}
		srv := newAlarmTestServer(store, notifier, secretManager)

		req := httptest.NewRequest("GET", "/check_alarm", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp map[string]int
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}
		if resp["fired"] != 0 {
			t.Errorf("expected fired 0, got %d", resp["fired"])
		}
	})

	t.Run("alarms found", func(t *testing.T) {
		config := models.UserConfig{
			DeviceID:        "550e8400-e29b-41d4-a716-446655440003",
			Lat:             34.0522,
			Long:            -118.2437,
			FCMToken:        "token-1",
			DayPreferences:  []bool{true, true, true, true, true, true, true},
			TimeZone:        "America/Los_Angeles",
			Offset:          0,
			Enabled:         true,
			NextAlarmTime:   "2026-01-15T14:32:00Z",
			AlarmDateBucket: "2026-01-15",
		}
		store := &mockStore{
			queryResults: []models.UserConfig{config},
			configs:      make(map[string]*models.UserConfig),
		}
		notifier := &mockNotifier{}
		secretManager := &mockSecretManager{secret: "fake-key"}
		srv := newAlarmTestServer(store, notifier, secretManager)

		req := httptest.NewRequest("GET", "/check_alarm", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp map[string]int
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}
		if resp["fired"] != 1 {
			t.Errorf("expected fired 1, got %d", resp["fired"])
		}
		if len(notifier.calls) != 1 {
			t.Errorf("expected 1 notification call, got %d", len(notifier.calls))
		}
		if notifier.calls[0].token != "token-1" {
			t.Errorf("notification token mismatch")
		}
		// Verify config was updated (store.PutConfig called)
		if _, ok := store.configs["550e8400-e29b-41d4-a716-446655440003"]; !ok {
			t.Errorf("config not updated")
		}
	})
}

func TestHealthEndpoint(t *testing.T) {
	store := &mockStore{configs: make(map[string]*models.UserConfig)}
	srv := newTestServer(store)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "OK" {
		t.Errorf("expected body OK, got %s", w.Body.String())
	}
}
