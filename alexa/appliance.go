package alexa

type Appliance struct {
	Id          string   `json:"applianceId"`
	Name        string   `json:"friendlyName"`
	Description string   `json:"friendlyDescription"`
	IsReachable bool     `json:"isReachable"`
	Actions     []string `json:"actions"`

	ManufacturerName string            `json:"manufacturerName"`
	ModelName        string            `json:"modelName"`
	Version          string            `json:"version"`
	Details          map[string]string `json:"additionalApplianceDetails"`
}
