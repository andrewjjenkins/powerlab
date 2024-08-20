package hpilo4

import (
	"fmt"

	"github.com/andrewjjenkins/powerlab/pkg/model"
)

type TemperatureSensorResponse struct {
	Label          string `json:"label"`
	XPosition      int    `json:"xposition"`
	YPosition      int    `json:"yposition"`
	Location       string `json:"location"`
	Status         string `json:"status"`
	CurrentReading int    `json:"currentreading"`
	Caution        int    `json:"caution"`
	Critical       int    `json:"critical"`
	TempUnit       string `json:"temp_unit"`
}

type TemperatureSensorsResponse struct {
	HostPowerState string                      `json:"hostpwr_state"`
	Temperatures   []TemperatureSensorResponse `json:"temperature"`
}

func (api *Api) getTemperaturesInternal() (*TemperatureSensorsResponse, error) {
	var temps TemperatureSensorsResponse
	err := api.GetJson("/json/health_temperature", &temps)
	if err != nil {
		return nil, err
	}
	return &temps, nil
}

type PowerSupplyResponse struct {
	Unhealthy        int    `json:"unhealthy"`
	Enabled          int    `json:"enabled"`
	Mismatch         int    `json:"mismatch"`
	PsBay            int    `json:"ps_bay"`
	PsPresent        string `json:"ps_present"`
	PsCondition      string `json:"ps_condition"`
	PsErrorCode      string `json:"ps_error_code"`
	PsIpduCapable    string `json:"ps_ipdu_capable"`
	PsHotplugCapable string `json:"ps_hotplug_capable"`
	PsModel          string `json:"ps_model"`
	PsSpare          string `json:"ps_spare"`
	PsSerialNum      string `json:"ps_serial_num"`
	PsMaxCapWatts    int    `json:"ps_max_cap_watts"`
	PsFwVer          string `json:"ps_fw_ver"`
	PsInputVolts     int    `json:"ps_input_volts"`
	PsOutputWatts    int    `json:"ps_output_watts"`
	Average          int    `json:"avg"`
	Max              int    `json:"max"`
	Supply           bool   `json:"supply"`
	Bbu              bool   `json:"bbu"`
	Charge           int    `json:"charge"`
	Age              int    `json:"age"`
	BatteryHealth    int    `json:"battery_health"`
}

type PowerSuppliesResponse struct {
	PresentPowerReading int                   `json:"present_power_reading"`
	Supplies            []PowerSupplyResponse `json:"supplies"`
}

func (api *Api) getPowerInternal() (*PowerSuppliesResponse, error) {
	var powers PowerSuppliesResponse
	err := api.GetJson("/json/power_supplies", &powers)
	if err != nil {
		return nil, err
	}
	return &powers, nil
}

func (api *Api) GetSensors() (*model.ServerSensorReadings, error) {
	temps, err := api.getTemperaturesInternal()
	if err != nil {
		return nil, err
	}

	tempsByName := make(map[string]TemperatureSensorResponse)
	for _, t := range temps.Temperatures {
		if _, ok := tempsByName[t.Label]; ok {
			return nil, fmt.Errorf("duplicate temp sensor %s", t.Label)
		}
		tempsByName[t.Label] = t
	}

	for _, l := range []string{"02-CPU 1", "12-PS 2 Inlet"} {
		if _, ok := tempsByName[l]; !ok {
			return nil, fmt.Errorf("no temp sensor %s", l)
		}
	}

	powerSupplies, err := api.getPowerInternal()
	if err != nil {
		return nil, err
	}

	return &model.ServerSensorReadings{
		CpuTemp:     float64(tempsByName["02-CPU 1"].CurrentReading),
		ChassisTemp: float64(tempsByName["12-PS 2 Inlet"].CurrentReading),
		FanSpeed:    0.0,
		PowerWatts:  float64(powerSupplies.PresentPowerReading),
	}, nil
}

type rawResponse struct {
	Temperatures  *TemperatureSensorsResponse `json:"temperatures"`
	PowerSupplies *PowerSuppliesResponse      `json:"power_supplies"`
}

func (api *Api) GetSensorsRaw() (interface{}, error) {
	temperatureResp, err := api.getTemperaturesInternal()
	if err != nil {
		return nil, err
	}
	powerSuppliesResponse, err := api.getPowerInternal()
	if err != nil {
		return nil, err
	}

	return &rawResponse{
		Temperatures:  temperatureResp,
		PowerSupplies: powerSuppliesResponse,
	}, nil
}
