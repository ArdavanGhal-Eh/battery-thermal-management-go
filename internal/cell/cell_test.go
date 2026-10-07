package cell

import (
	"testing"
)

func TestBatteryCell_OCV(t *testing.T) {
	c := NewBatteryCell(1, 1.0, 25.0)

	// At 100% SOC (1.0), OCV should be around 4.2V
	ocv100 := c.OCV(1.0)
	if ocv100 < 4.10 || ocv100 > 4.30 {
		t.Fatalf("Expected OCV at 100%% SOC to be ~4.2V, got %.3fV", ocv100)
	}

	// At 0% SOC, OCV should be around 3.1V
	ocv0 := c.OCV(0.0)
	if ocv0 < 3.0 || ocv0 > 3.20 {
		t.Fatalf("Expected OCV at 0%% SOC to be ~3.1V, got %.3fV", ocv0)
	}
}

func TestBatteryCell_DischargeHeating(t *testing.T) {
	c := NewBatteryCell(1, 0.90, 25.0)
	initTemp := c.State.Temperature
	initSOC := c.State.SOC

	// Discharge at 20A for 60 seconds (high C-rate discharge)
	current := 20.0
	dt := 1.0
	for i := 0; i < 60; i++ {
		c.Step(current, dt)
	}

	if c.State.SOC >= initSOC {
		t.Fatalf("Expected SOC to decrease during discharge, got %.3f", c.State.SOC)
	}

	if c.State.Temperature <= initTemp {
		t.Fatalf("Expected cell temperature to rise due to Joule heating, initial %.2fC, final %.2fC",
			initTemp, c.State.Temperature)
	}

	if c.State.Q_gen <= 0.0 {
		t.Fatalf("Expected heat generation Q_gen > 0 during high-current discharge")
	}
}

func TestBatteryCell_ChargeRecovery(t *testing.T) {
	c := NewBatteryCell(1, 0.50, 25.0)

	// Charge at -5A for 30 seconds
	current := -5.0
	c.Step(current, 30.0)

	if c.State.SOC <= 0.50 {
		t.Fatalf("Expected SOC to increase during charging, got %.4f", c.State.SOC)
	}
}
