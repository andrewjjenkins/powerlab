package megarac

type SensorResult struct {
	Id            int     `json:"id"`
	SensorNumber  int     `json:"sensor_number"`
	Name          string  `json:"name"`
	OwnerId       int     `json:"owner_id"`
	OwnerLun      int     `json:"owner_lun"`
	RawReading    float64 `json:"raw_reading"`
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

func (api *Api) GetSensors() (*SensorsResult, error) {
	s := SensorsResult{}
	err := api.GetJson("/api/sensors", &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (api *Api) GetSensorsRaw() (interface{}, error) {
	return api.GetSensors()
}
