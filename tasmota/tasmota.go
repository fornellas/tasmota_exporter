package tasmota

import (
	"github.com/fornellas/tasmota_exporter/prometheus"
)

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
	RealPowerW     float64 `json:"Power"`         // 225
	ApparentPowerW float64 `json:"ApparentPower"` // 283
	ReactivePowerW float64 `json:"ReactivePower"` // 172
	Factor         float64 `json:"Factor"`        // 0.79
	Volts          float64 `json:"Voltage"`       // 233
	Amps           float64 `json:"Current"`       // 1.21
}

func (e Energy) TimeSeriesGroup(prefix string) *prometheus.TimeSeriesGroup {
	tsg := prometheus.NewTimeSeriesGroup()

	ts, err := prometheus.NewTimeSeries(prefix+"_energy_kwh_total", prometheus.Labels{})
	if err != nil {
		panic(err)
	}
	ts.Set(e.TotalKWh)
	tsg.Add(ts)

	ts, err = prometheus.NewTimeSeries(prefix+"_energy_power_watts", prometheus.Labels{"type": "real"})
	if err != nil {
		panic(err)
	}
	ts.Set(e.RealPowerW)
	tsg.Add(ts)

	ts, err = prometheus.NewTimeSeries(prefix+"_energy_power_watts", prometheus.Labels{"type": "apparent"})
	if err != nil {
		panic(err)
	}
	ts.Set(e.ApparentPowerW)
	tsg.Add(ts)

	ts, err = prometheus.NewTimeSeries(prefix+"_energy_power_watts", prometheus.Labels{"type": "reactive"})
	if err != nil {
		panic(err)
	}
	ts.Set(e.ReactivePowerW)
	tsg.Add(ts)

	ts, err = prometheus.NewTimeSeries(prefix+"_energy_factor", prometheus.Labels{})
	if err != nil {
		panic(err)
	}
	ts.Set(e.Factor)
	tsg.Add(ts)

	ts, err = prometheus.NewTimeSeries(prefix+"_energy_voltage_volts", prometheus.Labels{})
	if err != nil {
		panic(err)
	}
	ts.Set(e.Volts)
	tsg.Add(ts)

	ts, err = prometheus.NewTimeSeries(prefix+"_energy_current_amps", prometheus.Labels{})
	if err != nil {
		panic(err)
	}
	ts.Set(e.Amps)
	tsg.Add(ts)

	return tsg
}

type Sensor struct {
	// FIXME convert from local time to UTC
	// Time   Time // "2023-11-19T14:26:14"
	Energy Energy
}

func (s Sensor) TimeSeriesGroup(prefix string) *prometheus.TimeSeriesGroup {
	tsg := prometheus.NewTimeSeriesGroup()
	tsg.Extend(s.Energy.TimeSeriesGroup("sensor"))
	return tsg
}

type Status struct {
	Sensor Sensor `json:"StatusSNS"`
}

func (s Status) TimeSeriesGroup() *prometheus.TimeSeriesGroup {
	tsg := prometheus.NewTimeSeriesGroup()
	tsg.Extend(s.Sensor.TimeSeriesGroup("tasmota"))
	return tsg
}
