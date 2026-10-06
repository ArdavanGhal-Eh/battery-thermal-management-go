package balancing

import (
	"battery-thermal-management-go/internal/cell"
	"math"
)

type ActiveBalancer struct {
	MaxBalancingCurrentA float64 // Maximum DC-DC transfer current (e.g. 2.0A)
	Efficiency           float64 // Converter efficiency ~ 90%
	ThresholdDeltaSOC    float64 // Minimum delta SOC to trigger balancing (e.g. 0.015)
}

func NewActiveBalancer() *ActiveBalancer {
	return &ActiveBalancer{
		MaxBalancingCurrentA: 2.0,
		Efficiency:           0.90,
		ThresholdDeltaSOC:    0.015,
	}
}

// BalancePack transfers charge from the highest-SOC cell to the lowest-SOC cell
func (b *ActiveBalancer) BalancePack(cells []*cell.BatteryCell, dt float64) (shuttledAh float64, balanced bool) {
	if len(cells) < 2 {
		return 0.0, false
	}

	maxIdx := 0
	minIdx := 0
	maxSOC := -1.0
	minSOC := 2.0

	for i, c := range cells {
		if c.State.SOC > maxSOC {
			maxSOC = c.State.SOC
			maxIdx = i
		}
		if c.State.SOC < minSOC {
			minSOC = c.State.SOC
			minIdx = i
		}
	}

	deltaSOC := maxSOC - minSOC
	if deltaSOC < b.ThresholdDeltaSOC {
		return 0.0, false
	}

	// Dynamic current proportional to imbalance
	transferCurrent := math.Min(b.MaxBalancingCurrentA, deltaSOC*10.0)
	transferAh := transferCurrent * (dt / 3600.0)

	// Shuttle charge: discharge highest cell, charge lowest cell with efficiency loss
	cells[maxIdx].State.SOC -= transferAh / cells[maxIdx].CapacityAh
	cells[minIdx].State.SOC += (transferAh * b.Efficiency) / cells[minIdx].CapacityAh

	return transferAh, true
}
