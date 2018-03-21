package result

type Result struct {
	Ip       string  `json:"ip"`
	Port     int     `json:"port"`
	Location string  `json:"location,omitempty"`
	Source   string  `json:"source"`
	Speed    float64 `json:"speed,omitempty"`
}
