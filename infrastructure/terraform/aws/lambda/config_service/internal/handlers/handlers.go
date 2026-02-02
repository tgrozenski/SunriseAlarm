package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"myproject/internal/dynamodb"
	"myproject/internal/models"
	"myproject/internal/notification"
	"myproject/internal/sunrise"
	"myproject/internal/validation"
)

type ConfigHandler struct {
	store dynamodb.ConfigStore
}

func NewConfigHandler(store dynamodb.ConfigStore) *ConfigHandler {
	return &ConfigHandler{store: store}
}

type AlarmHandler struct {
	store         dynamodb.ConfigStore
	notifier      notification.Notifier
	secretManager notification.SecretManager
}

func NewAlarmHandler(store dynamodb.ConfigStore, notifier notification.Notifier, secretManager notification.SecretManager) *AlarmHandler {
	return &AlarmHandler{store: store, notifier: notifier, secretManager: secretManager}
}

func (h *AlarmHandler) CheckAlarm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now().UTC()
	dateBucket := now.Format("2006-01-02")
	startTime := now.Add(-1 * time.Minute).Format(time.RFC3339)
	endTime := now.Add(1 * time.Minute).Format(time.RFC3339)

	configs, err := h.store.QueryByAlarmTime(ctx, dateBucket, startTime, endTime)
	if err != nil {
		writeError(w, "failed to query alarms", http.StatusInternalServerError)
		return
	}

	if len(configs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"fired": 0})
		return
	}

	var notifier notification.Notifier
	if h.notifier != nil {
		notifier = h.notifier
	} else {
		secret, err := h.secretManager.GetSecret(ctx, "LAMBDA_SERVICE_KEY")
		if err != nil {
			writeError(w, "failed to retrieve service key", http.StatusInternalServerError)
			return
		}
		notifier = notification.NewFCMNotifier(secret)
	}
	fired := 0
	for _, config := range configs {
		if err := notifier.SendNotification(ctx, config.FCMToken, config.DeviceID); err != nil {
			// Notification failed, but we still need to reschedule the alarm
			// Log error but continue to reschedule
			log.Printf("failed to send notification for device %s: %v", config.DeviceID, err)
		} else {
			fired++
		}
		// Always reschedule after attempting notification
		nextAlarmTime, alarmDateBucket, err := sunrise.ComputeAlarmFields(&config)
		if err != nil {
			log.Printf("failed to compute alarm fields for device %s: %v", config.DeviceID, err)
			continue
		}
		config.NextAlarmTime = nextAlarmTime
		config.AlarmDateBucket = alarmDateBucket
		if err := h.store.PutConfig(ctx, &config); err != nil {
			log.Printf("failed to update config for device %s: %v", config.DeviceID, err)
			continue
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"fired": fired})
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["deviceId"]
	if deviceID == "" {
		writeError(w, "deviceId is required", http.StatusBadRequest)
		return
	}

	config, err := h.store.GetConfig(r.Context(), deviceID)
	if err != nil {
		if err == dynamodb.ErrConfigNotFound {
			writeError(w, "config not found", http.StatusNotFound)
			return
		}
		writeError(w, "failed to retrieve config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(config)
}

func (h *ConfigHandler) PutConfig(w http.ResponseWriter, r *http.Request) {
	var config models.UserConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validation.ValidateUserConfig(&config); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle manual override
	if !config.Enabled {
		config.NextAlarmTime = ""
		config.AlarmDateBucket = "DISABLED"
	} else if config.NextAlarmTime == "" {
		// Calculate from sunrise (no manual override provided)
		nextAlarmTime, alarmDateBucket, err := sunrise.ComputeAlarmFields(&config)
		if err != nil {
			log.Printf("failed to compute alarm fields: %v", err)
			writeError(w, "failed to compute next alarm time", http.StatusInternalServerError)
			return
		}
		config.NextAlarmTime = nextAlarmTime
		config.AlarmDateBucket = alarmDateBucket
	}

	// If NextAlarmTime was already set, keep it as-is (manual override for testing)
	if err := h.store.PutConfig(r.Context(), &config); err != nil {
		log.Printf("failed to save config: %v", err)
		writeError(w, "failed to save config", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
