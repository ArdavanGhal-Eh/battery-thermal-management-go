# ⚡ Real-Time EV Battery Pack Electro-Thermal & SOH Degradation Simulator (Go + Python)

A high-performance, concurrent electro-thermal simulation engine for lithium-ion battery modules (18650 / 21700 NMC chemistry) developed in **Go** and paired with **Python** analytical telemetry visualizers. Models 2-RC Thevenin equivalent circuit polarization, Bernardi heat generation, lumped convective cooling, and concurrent multi-cell pack balancing using lightweight Goroutines.

---

## 📌 Problem & Engineering Motivation
In Electric Vehicles (EVs), Automated Guided Vehicles (AGVs), and grid energy storage (BESS):
1. **Thermal Runaway Risk:** Overheating beyond $55-60^\circ\text{C}$ or spatial gradients $\Delta T_{pack} > 5^\circ\text{C}$ causes catastrophic thermal runaway and accelerated aging.
2. **Cell-to-Cell Imbalance:** Manufacturing variances in internal resistance ($R_0$) and capacity cause uneven current distribution and hotspot formation in module cores.
3. **Real-Time Onboard Simulation:** Battery Management Systems (BMS) require microsecond-fast digital twins to predict temperature rise under peak discharge bursts (e.g., hill climb or regenerative braking) without draining onboard computing power.

---

## 🌟 System Architecture

```text
┌───────────────────────────────────────────────┐
│        Dynamic Drive Cycle Current I(t)       │
└───────────────────────┬───────────────────────┘
                        │
                        ▼
┌───────────────────────────────────────────────┐
│        Multi-Cell Pack Manager (Go)           │
│   (16S4P = 64 Cells Simulated via Goroutines) │
└───────┬───────────────────────────────┬───────┘
        │                               │
        ▼                               ▼
┌──────────────────────────┐    ┌──────────────────────────┐
│  2-RC Thevenin Model     │    │  Bernardi Thermal Model  │
│  - OCV(SOC) Nonlinear Fit│    │  - Ohmic Joule Q_ohmic   │
│  - V_rc1, V_rc2 Dynamics │    │  - Lumped Heat Equation  │
│  - Coulomb Counting SOC  │    │  - Convective Heat h*A*ΔT│
└──────────────────────────┘    └──────────────────────────┘
                        │
                        ▼
┌───────────────────────────────────────────────┐
│     Live Telemetry & Hotspot Guard (Go)       │
│     - REST API Endpoint (:8088)               │
│     - Thermal Runaway & Gradient Alarms       │
└───────────────────────────────────────────────┘
```

---

## 📐 Mathematical Formulation

### 1. 2-RC Thevenin Equivalent Circuit
The cell terminal voltage $V_{term}$:
$$V_{term} = V_{ocv}(\text{SOC}) - I R_0 - V_{rc1} - V_{rc2}$$
Polarization dynamics across double-layer and diffusion capacitances:
$$\frac{dV_{rc1}}{dt} = -\frac{V_{rc1}}{R_1 C_1} + \frac{I}{C_1}, \quad \frac{dV_{rc2}}{dt} = -\frac{V_{rc2}}{R_2 C_2} + \frac{I}{C_2}$$

### 2. Bernardi Heat Generation & Lumped Thermal Balance
$$Q_{gen} = I^2 R_0 + I (V_{rc1} + V_{rc2})$$
$$m C_p \frac{dT}{dt} = Q_{gen} - h A (T - T_{coolant})$$

### 3. State of Health (SOH) Capacity Fade
$$\Delta \text{SOH} = -\alpha \cdot Ah_{throughput} \cdot \left(1 + \beta \max(0, T - 25^\circ\text{C})\right)$$

---

## 🎯 Real-World Applications & Cross-Industry Impact

### ⚙️ Mechanical & Automotive Engineering
- **EV Battery Thermal Management Systems (BTMS):** Optimizing liquid cold-plate channel geometries to eliminate core hotspots in 800V packs.
- **Heavy Machinery & Mining Electrification:** Safe high-C-rate fast charging under extreme environmental ambients.

### 🌐 Cross-Industry & Software Applications
- **BMS Digital Twin Microservices:** Cloud telemetry platforms (like Tesla Fleet Analytics) continuously inferring cell SOH from field log streams.
- **Renewable Microgrid BESS:** Scheduling peak-shaving storage discharge without triggering battery degradation penalties.

---

## 🚀 Installation & Execution

### 1. Run Go Simulation Engine
```bash
go run cmd/simulator/main.go
```
*Live telemetry will stream in the console and serve JSON on `http://localhost:8088/api/battery/status`.*

### 2. Generate Thermal Contours & Voltage Plots (Python)
```bash
cd python_visualizer
pip install -r requirements.txt
python plot_pack_thermal.py
```

---

## 🛠️ Tech Stack
- **Engine Core:** Go 1.21+, Goroutines, Channels, HTTP server
- **Electro-Thermal Modeling:** 2-RC Thevenin ECM, Bernardi heat generation
- **Visualizer & Analytics:** Python 3.10+, `numpy`, `matplotlib`

---

## 👨‍💻 Author
**Ardavan Ghal-Eh**  
Mechanical Engineering Student, Sharif University of Technology  
*Focus: Energy Systems, Battery Thermal Management (BTMS) & High-Performance Software*
