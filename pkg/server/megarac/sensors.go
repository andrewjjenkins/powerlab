package megarac

import (
	"fmt"
	"log/slog"

	"github.com/andrewjjenkins/powerlab/pkg/model"
)

type SensorResult struct {
	Id            int     `json:"id"`
	SensorNumber  int     `json:"sensor_number"`
	Name          string  `json:"name"`
	OwnerId       int     `json:"owner_id"`
	OwnerLun      int     `json:"owner_lun"`
	RawReading    float64 `json:"raw_reading"`
	Reading       float64 `json:"reading"`
	Type          string  `json:"type"`
	TypeNumber    int     `json:"type_number"`
	SensorState   int     `json:"sensor_state"`
	DiscreteState int     `json:"discrete_state"`
	Accessible    int     `json:"accessible"`
	SettableFlag  int     `json:"settable_flag"`
	//LowerNonRecoverableThreshold  string  `json:"lower_non_recoverable_threshold"`
	//LowerCriticalThreshold        float64 `json:"lower_critical_threshold"`
	//LowerNonCriticalThreshold     float64 `json:"lower_non_critical_threshold"`
	//HigherNonCriticalThreshold    float64 `json:"higher_non_critical_threshold"`
	//HigherCriticalThreshold       float64 `json:"higher_critical_threshold"`
	//HigherNonRecoverableThreshold string  `json:"higher_non_recoverable_threshold"`
	Unit string `json:"unit"`
}

type SensorsResult []SensorResult

func (api *Api) getSensorsInternal() (*SensorsResult, error) {
	u := "/api/sensors"
	cached, ok := api.cache.Responses.Get(u)
	if ok {
		res := cached.(SensorsResult)
		if res == nil {
			return nil, fmt.Errorf("cached result invalid type")
		}
		slog.Debug("Cache hit", "path", u)
		return &res, nil
	}

	s := SensorsResult{}
	err := api.GetJson(u, &s)
	if err != nil {
		return nil, err
	}
	api.cache.Responses.Add(u, s)
	return &s, nil
}

func (api *Api) GetSensorsRaw() (interface{}, error) {
	return api.getSensorsInternal()
}

func (api *Api) GetSensors() (*model.ServerSensorReadings, error) {
	sr, err := api.getSensorsInternal()
	if err != nil {
		return nil, err
	}

	sensorsByName := make(map[string]SensorResult)
	for _, s := range *sr {
		if _, ok := sensorsByName[s.Name]; ok {
			return nil, fmt.Errorf("duplicate sensor %s", s.Name)
		}
		sensorsByName[s.Name] = s
	}

	for _, k := range []string{"CPU0_TEMP", "MB_TEMP1", "BPB_FAN_1A", "SYS_POWER"} {
		if _, ok := sensorsByName[k]; !ok {
			return nil, fmt.Errorf("no sensor %s", k)
		}
	}

	readings := model.ServerSensorReadings{
		CpuTemp:     sensorsByName["CPU0_TEMP"].Reading,
		ChassisTemp: sensorsByName["MB_TEMP1"].Reading,
		FanSpeed:    sensorsByName["BPB_FAN_1A"].Reading,
		PowerWatts:  sensorsByName["SYS_POWER"].Reading,
	}
	return &readings, nil
}
