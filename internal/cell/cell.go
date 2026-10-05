package cell

import (
	"math"
)

// CellState represents the electro-thermal state of a single lithium-ion cell
type CellState struct {
	ID          int     `json:"cell_id"`
	SOC         float64 `json:"soc"`          // State of Charge [0.0 - 1.0]
	SOH         float64 `json:"soh"`          // State of Health [0.0 - 1.0]
	V_terminal  float64 `json:"v_terminal"`   // Terminal voltage (V)
	V_rc1       float64 `json:"v_rc1"`        // First RC pair polarization voltage
	V_rc2       float64 `json:"v_rc2"`        // Second RC pair polarization voltage
	Temperature float64 `json:"temp_c"`       // Cell temperature in Celsius
	Q_gen       float64 `json:"q_gen_watts"`  // Total heat generation rate (Watts)
}

// BatteryCell encapsulates parameters for an 18650 / 21700 NMC cell
type BatteryCell struct {
	State       CellState
	CapacityAh  float64 // Nominal capacity (e.g. 5.0 Ah for 21700)
	MassKg      float64 // Cell mass (e.g. 0.070 kg)
	Cp          float64 // Specific heat capacity (J/(kg*K)) ~ 950
	R0          float64 // Ohmic internal resistance (Ohms) ~ 0.020
	R1          float64 // Charge transfer resistance ~ 0.015
	C1          float64 // Electrochemical double-layer capacitance (F) ~ 1500
	R2          float64 // Diffusion resistance ~ 0.010
	C2          float64 // Diffusion capacitance (F) ~ 8000
	H_conv      float64 // Convective cooling coefficient (W/(m^2*K)) ~ 25.0
	Area_m2     float64 // Surface area ~ 0.0055 m^2
	T_coolant   float64 // Coolant temperature (Celsius)
}

func NewBatteryCell(id int, initialSOC float64, tempCoolant float64) *BatteryCell {
	return &BatteryCell{
		State: CellState{
			ID:          id,
			SOC:         initialSOC,
			SOH:         1.0,
			V_terminal:  3.7,
			Temperature: tempCoolant,
		},
		CapacityAh: 5.0,
		MassKg:     0.070,
		Cp:         950.0,
		R0:         0.018,
		R1:         0.012,
		C1:         1800.0,
		R2:         0.010,
		C2:         6000.0,
		H_conv:     30.0,
		Area_m2:    0.0058,
		T_coolant:  tempCoolant,
	}
}

// OCV estimates Open Circuit Voltage from SOC using standard NMC polynomial
func (c *BatteryCell) OCV(soc float64) float64 {
	// Empirical fit for 21700 NMC chemistry: 3.0V at 0% to 4.2V at 100%
	return 3.10 + 1.10*soc - 0.45*math.Pow(soc, 2) + 0.45*math.Pow(soc, 3)
}

// Step advances cell electro-thermal dynamics by dt seconds with current I (Amperes, positive=discharge)
func (c *BatteryCell) Step(currentAmps float64, dt float64) {
	// 1. Coulomb counting for SOC
	deltaSOC := -(currentAmps * (dt / 3600.0)) / (c.CapacityAh * c.State.SOH)
	c.State.SOC = math.Max(0.0, math.Min(1.0, c.State.SOC+deltaSOC))

	// 2. RC polarization voltage updates
	tau1 := c.R1 * c.C1
	tau2 := c.R2 * c.C2
	c.State.V_rc1 = c.State.V_rc1*math.Exp(-dt/tau1) + currentAmps*c.R1*(1.0-math.Exp(-dt/tau1))
	c.State.V_rc2 = c.State.V_rc2*math.Exp(-dt/tau2) + currentAmps*c.R2*(1.0-math.Exp(-dt/tau2))

	// 3. Terminal voltage
	vocv := c.OCV(c.State.SOC)
	c.State.V_terminal = vocv - currentAmps*c.R0 - c.State.V_rc1 - c.State.V_rc2

	// 4. Heat Generation (Bernardi Model: Ohmic Joule heating + Polarization)
	qOhmic := currentAmps * currentAmps * c.R0
	qPolarization := currentAmps * (c.State.V_rc1 + c.State.V_rc2)
	c.State.Q_gen = math.Max(0.0, qOhmic+qPolarization)

	// 5. Lumped Thermal Mass Conservation: m * Cp * dT/dt = Q_gen - h * A * (T - T_coolant)
	qCooling := c.H_conv * c.Area_m2 * (c.State.Temperature - c.T_coolant)
	dTdt := (c.State.Q_gen - qCooling) / (c.MassKg * c.Cp)
	c.State.Temperature += dTdt * dt

	// 6. Empirical SOH degradation per Ah throughput
	ahThroughput := math.Abs(currentAmps) * (dt / 3600.0)
	c.State.SOH -= 2.5e-6 * ahThroughput * (1.0 + 0.05*math.Max(0.0, c.State.Temperature-25.0))
}
