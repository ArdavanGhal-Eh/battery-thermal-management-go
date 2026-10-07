package pack

import (
	"testing"
)

func TestBatteryPack_Initialization(t *testing.T) {
	series := 10
	parallel := 2
	p := NewBatteryPack(series, parallel, 25.0)

	if len(p.Cells) != 20 {
		t.Fatalf("Expected 20 cells in 10s2p pack, got %d", len(p.Cells))
	}

	status := p.StepPack(0.0, 0.1)
	if status.NumCells != 20 {
		t.Fatalf("Expected status.NumCells = 20, got %d", status.NumCells)
	}

	// 10 series NMC cells ~ 37V-42V
	if status.PackVoltage < 35.0 || status.PackVoltage > 43.0 {
		t.Fatalf("Expected pack voltage between 35V and 43V, got %.2fV", status.PackVoltage)
	}
}

func TestBatteryPack_StepDischarge(t *testing.T) {
	p := NewBatteryPack(4, 1, 25.0)
	status1 := p.StepPack(10.0, 1.0)
	status2 := p.StepPack(10.0, 10.0)

	if status2.AverageSOC >= status1.AverageSOC {
		t.Fatalf("Expected AverageSOC to drop after discharge, got %.1f vs %.1f",
			status2.AverageSOC, status1.AverageSOC)
	}
}

func TestBatteryPack_ThermalRunawayWarning(t *testing.T) {
	p := NewBatteryPack(4, 1, 25.0)

	// Manually inject extreme heat into one cell
	p.Cells[0].State.Temperature = 65.0

	status := p.StepPack(5.0, 0.1)
	if !status.ThermalRunaway {
		t.Fatalf("Expected ThermalRunaway warning to be TRUE when cell temp exceeds 58C")
	}
}
