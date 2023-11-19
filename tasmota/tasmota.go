package tasmota

// type Time time.Time

// func (t *Time) UnmarshalJSON(b []byte) error {
// 	gt, err := time.Parse("2006-01-02T15:04:05", string(b[1:len(b)-1]))
// 	*t = Time(gt)
// 	return err
// }

type Energy struct {
	// FIXME convert from local time to UTC
	// TotalStartTime Time    `json:"TotalStartTime"` // "2023-11-16T20:37:21"
	TotalKWh float64 `json:"Total"` // 3.396
	// YesterdayKWh   float64 `json:"Yesterday"`     // 0.870
	// TodayKWh       float64 `json:"Today"`         // 0.645
	PowerW         float64 `json:"Power"`         // 225
	ApparentPowerW float64 `json:"ApparentPower"` // 283
	ReactivePowerW float64 `json:"ReactivePower"` // 172
	Factor         float64 `json:"Factor"`        // 0.79
	Volts          float64 `json:"Voltage"`       // 233
	Amps           float64 `json:"Current"`       // 1.21
}
type Sensor struct {
	// FIXME convert from local time to UTC
	// Time   Time // "2023-11-19T14:26:14"
	Energy Energy
}

type Status struct {
	Sensor Sensor `json:"StatusSNS"`
}
