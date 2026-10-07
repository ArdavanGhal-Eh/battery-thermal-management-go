package balancing

import (
	"battery-thermal-management-go/internal/cell"
	"testing"
)

func TestActiveBalancer_EqualCellsNoBalancing(t *testing.T) {
	b := NewActiveBalancer()
	cells := []*cell.BatteryCell{
		cell.NewBatteryCell(1, 0.80, 25.0),
		cell.NewBatteryCell(2, 0.805, 25.0),
	}

	shuttled, balanced := b.BalancePack(cells, 1.0)
	if balanced || shuttled > 0.0 {
		t.Fatalf("Expected no balancing when delta SOC < threshold (0.015)")
	}
}

func TestActiveBalancer_ImbalancedCellsChargeTransfer(t *testing.T) {
	b := NewActiveBalancer()
	cellHigh := cell.NewBatteryCell(1, 0.90, 25.0)
	cellLow := cell.NewBatteryCell(2, 0.70, 25.0)
	cells := []*cell.BatteryCell{cellHigh, cellLow}

	initialHighSOC := cellHigh.State.SOC
	initialLowSOC := cellLow.State.SOC

	shuttled, balanced := b.BalancePack(cells, 3600.0) // 1 hour step

	if !balanced {
		t.Fatalf("Expected balancing to trigger on delta SOC = 0.20")
	}

	if shuttled <= 0.0 {
		t.Fatalf("Expected positive charge shuttling, got %.4f Ah", shuttled)
	}

	if cellHigh.State.SOC >= initialHighSOC {
		t.Fatalf("Expected high cell SOC to decrease, got %.4f vs %.4f", cellHigh.State.SOC, initialHighSOC)
	}

	if cellLow.State.SOC <= initialLowSOC {
		t.Fatalf("Expected low cell SOC to increase, got %.4f vs %.4f", cellLow.State.SOC, initialLowSOC)
	}
}
