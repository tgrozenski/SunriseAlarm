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
)

type mockStore struct {
	configs map[string]*models.UserConfig
	err     error
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

func TestGetConfigEndpoint(t *testing.T) {
	store := &mockStore{
		configs: map[string]*models.UserConfig{
			"test-id": {
				DeviceID:       "test-id",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
		},
	}
	srv := newTestServer(store) //Stevie Ray Vaughan

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/config/test-id", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var config models.UserConfig
		if err := json.Unmarshal(w.Body.Bytes(), &config); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}
		if config.DeviceID != "test-id" {
			t.Errorf("expected deviceId test-id, got %s", config.DeviceID)
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
			DeviceID:       "new-id",
			Lat:            34.0522,
			Long:           -118.2437,
			FCMToken:       "token",
			DayPreferences: []bool{false, false, false, true, true, true, true},
			TimeZone:       "America/Los_Angeles",
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

		stored, err := store.GetConfig(context.Background(), "new-id")
		if err != nil {
			t.Errorf("config not stored: %v", err)
		}
		if stored.DeviceID != "new-id" {
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
