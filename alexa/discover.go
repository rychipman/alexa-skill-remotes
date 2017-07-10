package alexa

type DiscoverRequest struct {
	Header  Header
	Payload DiscoverRequestPayload
}

type DiscoverRequestPayload struct {
	AccessToken string
}

type DiscoverResponse struct {
	Header  Header                  `json:"header"`
	Payload DiscoverResponsePayload `json:"payload"`
}

type DiscoverResponsePayload struct {
	Appliances []Appliance `json:"discoveredAppliances"`
}
