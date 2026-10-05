package telemetry

import (
	"battery-thermal-management-go/internal/pack"
	"encoding/json"
	"net/http"
	"sync"
)

type TelemetryServer struct {
	Pack   *pack.BatteryPack
	Status pack.PackStatus
	mu     sync.RWMutex
}

func NewTelemetryServer(bp *pack.BatteryPack) *TelemetryServer {
	return &TelemetryServer{Pack: bp}
}

func (s *TelemetryServer) UpdateStatus(status pack.PackStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = status
}

func (s *TelemetryServer) StatusHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.Status)
}
