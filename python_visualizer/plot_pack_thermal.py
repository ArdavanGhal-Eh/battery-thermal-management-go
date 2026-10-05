"""
Battery Pack Thermal Gradient & Degradation Visualizer
Part of Battery Pack Electro-Thermal Suite in Go & Python
Author: Ardavan Ghal-Eh | Sharif University of Technology
"""

import numpy as np
import matplotlib.pyplot as plt


def simulate_and_plot_battery_pack():
    time_s = np.linspace(0, 600, 300)
    
    # Pack average & max temperature curves
    t_ambient = 25.0
    t_avg = t_ambient + 18.0 * (1.0 - np.exp(-time_s / 180.0)) + 1.5 * np.sin(0.1 * time_s)
    t_max = t_avg + 3.2 * (1.0 - np.exp(-time_s / 120.0))
    t_min = t_avg - 2.1 * (1.0 - np.exp(-time_s / 120.0))

    # SOC curve
    soc = 95.0 - (time_s / 600.0) * 45.0

    # 16S4P 2D Thermal spatial map at t=600s
    rows, cols = 4, 16
    temp_grid = np.zeros((rows, cols))
    for r in range(rows):
        for c in range(cols):
            # Center cells heat up more due to restricted convective cooling
            dist_from_center = np.hypot(r - 1.5, c - 7.5) / np.hypot(1.5, 7.5)
            temp_grid[r, c] = t_max[-1] - 4.5 * dist_from_center + np.random.uniform(-0.3, 0.3)

    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(13, 8), gridspec_kw={'height_ratios': [1.2, 1.0]})

    # 1. Spatial Cell Temperature Map (16S4P Pack)
    cmap = plt.colormaps['plasma']
    im = ax1.imshow(temp_grid, cmap=cmap, aspect='auto', interpolation='nearest')
    cbar = fig.colorbar(im, ax=ax1, label='Cell Temperature (°C)')
    ax1.set_title("16S4P Battery Module Spatial Thermal Distribution (Hotspot Detection at Core)", fontsize=12, fontweight='bold')
    ax1.set_xlabel("Series Cell Index (1 to 16)")
    ax1.set_ylabel("Parallel Bank (1 to 4)")
    ax1.set_xticks(range(16))
    ax1.set_xticklabels([f"S{i+1}" for i in range(16)])
    ax1.set_yticks(range(4))
    ax1.set_yticklabels([f"P{i+1}" for i in range(4)])

    # Annotate cell temperatures
    for r in range(rows):
        for c in range(cols):
            ax1.text(c, r, f"{temp_grid[r, c]:.1f}", ha='center', va='center', color='white' if temp_grid[r, c] < 42 else 'black', fontsize=7)

    # 2. Time-series Temperature & SOC
    ax2.plot(time_s, t_max, 'r-', label='Maximum Cell Temp ($T_{max}$)', linewidth=2.0)
    ax2.plot(time_s, t_avg, 'orange', label='Average Pack Temp ($T_{avg}$)', linewidth=1.8)
    ax2.plot(time_s, t_min, 'g--', label='Minimum Cell Temp ($T_{min}$)', linewidth=1.5)
    ax2.axhline(55.0, color='darkred', linestyle=':', label='Thermal Throttling Limit (55°C)', linewidth=1.5)

    ax2_twin = ax2.twinx()
    ax2_twin.plot(time_s, soc, 'b-', label='Pack SOC (%)', linewidth=1.8)
    ax2_twin.set_ylabel('State of Charge - SOC (%)', color='blue')

    ax2.set_title("Thermal Rise & SOC Depletion under Aggressive Drive Cycle", fontsize=11, fontweight='bold')
    ax2.set_xlabel("Time (seconds)")
    ax2.set_ylabel("Temperature (°C)")
    ax2.grid(True, linestyle=':', alpha=0.6)
    ax2.legend(loc='upper left')

    plt.tight_layout()
    output_png = "/working_dir/c_db360fcc6464ba55/daily_projects_day6/battery-thermal-management-go/python_visualizer/battery_thermal_profile.png"
    plt.savefig(output_png, dpi=200)
    print("✅ Battery Thermal Visualization generated:", output_png)


if __name__ == "__main__":
    simulate_and_plot_battery_pack()
