package main

import (
	"battery-thermal-management-go/internal/pack"
	"battery-thermal-management-go/internal/telemetry"
	"fmt"
	"math"
	"net/http"
	"time"
)

func main() {
	fmt.Println("============================================================")
	fmt.Println("⚡ Real-Time EV Battery Pack Electro-Thermal Simulator (Go)")
	fmt.Println("   Author: Ardavan Ghal-Eh | Sharif University of Tech")
	fmt.Println("============================================================")

	// 16 Series x 4 Parallel = 64 cells (~60V, 20Ah pack for light EV / AMR)
	series := 16
	parallel := 4
	ambientCoolantC := 25.0

	bp := pack.NewBatteryPack(series, parallel, ambientCoolantC)
	ts := telemetry.NewTelemetryServer(bp)

	// Start background HTTP telemetry listener
	http.HandleFunc("/api/battery/status", ts.StatusHandler)
	go func() {
		http.ListenAndServe(":8088", nil)
	}()

	fmt.Printf("Pack Config: %dS%dP (%d cells total)\n", series, parallel, series*parallel)
	fmt.Printf("Nominal Capacity: %.1f Ah | Coolant Temp: %.1f °C\n", float64(parallel)*5.0, ambientCoolantC)
	fmt.Println("Simulation Drive Profile: Aggressive Highway Cycle (3C Discharge Pulses)")
	fmt.Println("Telemetry Stream: http://localhost:8088/api/battery/status")
	fmt.Println("------------------------------------------------------------")

	dt := 0.2 // 200ms simulation timestep
	totalSteps := 100

	startTime := time.Now()

	for step := 1; step <= totalSteps; step++ {
		t := float64(step) * dt

		// Dynamic drive profile current: Base 25A + 30A acceleration pulse every 10s
		pulse := 0.0
		if math.Mod(t, 10.0) < 4.0 {
			pulse = 35.0
		}
		packCurrent := 20.0 + pulse + 5.0*math.Sin(0.5*t)

		status := bp.StepPack(packCurrent, dt)
		ts.UpdateStatus(status)

		if step%20 == 0 || step == totalSteps {
			runawayAlert := "Normal"
			if status.ThermalRunaway {
				runawayAlert = "⚠️ CRITICAL TEMP ALERT"
			}
			fmt.Printf("Time: %5.1fs | V_pack: %5.1fV | Current: %5.1fA | Power: %5.2fkW | SOC: %4.1f%% | T_avg: %5.2f°C | T_max: %5.2f°C | ΔT: %4.2f°C [%s]\n",
				t, status.PackVoltage, status.PackCurrent, status.TotalPowerKW, status.AverageSOC, status.AverageTemp, status.MaxTemp, status.TempGradient, runawayAlert)
		}
	}

	elapsed := time.Since(startTime)
	cellUpdates := float64(totalSteps*series*parallel) / elapsed.Seconds()
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("⏱️ Completed %d timesteps in %v\n", totalSteps, elapsed)
	fmt.Printf("🚀 Concurrent Throughput: %.0f cell state updates/sec\n", cellUpdates)
	fmt.Println("============================================================")
}
