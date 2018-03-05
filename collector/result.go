package collector

type Result struct {
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Location string `json:"location,omitempty"`
	Source   string `json:"source"`
	LiveTime int    `json:"live_time,omitempty"`
}
