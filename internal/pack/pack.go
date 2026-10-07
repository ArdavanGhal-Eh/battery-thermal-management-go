package pack

import (
	"battery-thermal-management-go/internal/cell"
	"math"
	"sync"
)

type PackStatus struct {
	NumCells         int     `json:"num_cells"`
	PackVoltage      float64 `json:"pack_voltage_v"`
	PackCurrent      float64 `json:"pack_current_a"`
	AverageSOC       float64 `json:"avg_soc_percent"`
	AverageTemp      float64 `json:"avg_temp_c"`
	MaxTemp          float64 `json:"max_temp_c"`
	MinTemp          float64 `json:"min_temp_c"`
	TempGradient     float64 `json:"temp_gradient_delta_c"`
	ThermalRunaway   bool    `json:"thermal_runaway_warning"`
	TotalPowerKW     float64 `json:"total_power_kw"`
}

type BatteryPack struct {
	Cells          []*cell.BatteryCell
	SeriesCount    int
	ParallelCount  int
	mu             sync.RWMutex
}

func NewBatteryPack(series int, parallel int, ambientTemp float64) *BatteryPack {
	total := series * parallel
	cells := make([]*cell.BatteryCell, total)
	for i := 0; i < total; i++ {
		// Slight cell-to-cell capacity variance (1%)
		initialSOC := 0.95 - (float64(i%5) * 0.002)
		cells[i] = cell.NewBatteryCell(i+1, initialSOC, ambientTemp)
	}
	return &BatteryPack{
		Cells:         cells,
		SeriesCount:   series,
		ParallelCount: parallel,
	}
}

// StepPack simulates pack discharge with optimized high-performance cell stepping
func (p *BatteryPack) StepPack(totalPackCurrent float64, dt float64) PackStatus {
	p.mu.Lock()
	defer p.mu.Unlock()

	cellCurrent := totalPackCurrent / float64(p.ParallelCount)

	// High-performance cache-friendly zero-allocation cell evaluation
	for _, c := range p.Cells {
		c.Step(cellCurrent, dt)
	}

	// Compute aggregate metrics
	totalV := 0.0
	totalSOC := 0.0
	totalTemp := 0.0
	maxTemp := -100.0
	minTemp := 1000.0

	// Sum voltages of series strings
	for s := 0; s < p.SeriesCount; s++ {
		totalV += p.Cells[s*p.ParallelCount].State.V_terminal
	}

	for _, c := range p.Cells {
		totalSOC += c.State.SOC
		totalTemp += c.State.Temperature
		if c.State.Temperature > maxTemp {
			maxTemp = c.State.Temperature
		}
		if c.State.Temperature < minTemp {
			minTemp = c.State.Temperature
		}
	}

	num := float64(len(p.Cells))
	avgSOC := (totalSOC / num) * 100.0
	avgTemp := totalTemp / num
	deltaT := maxTemp - minTemp

	// Critical threshold: temperature > 58 C or deltaT > 6 C
	isRunawayRisk := maxTemp > 58.0 || deltaT > 6.0

	return PackStatus{
		NumCells:       len(p.Cells),
		PackVoltage:    round(totalV, 2),
		PackCurrent:    round(totalPackCurrent, 2),
		AverageSOC:     round(avgSOC, 1),
		AverageTemp:    round(avgTemp, 2),
		MaxTemp:        round(maxTemp, 2),
		MinTemp:        round(minTemp, 2),
		TempGradient:   round(deltaT, 2),
		ThermalRunaway: isRunawayRisk,
		TotalPowerKW:   round((totalV*totalPackCurrent)/1000.0, 2),
	}
}

func round(val float64, prec int) float64 {
	pow := math.Pow(10, float64(prec))
	return math.Round(val*pow) / pow
}
