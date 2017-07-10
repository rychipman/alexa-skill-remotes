package alexa

import "fmt"

type ControlRequest struct {
	Header  Header
	Payload ControlRequestPayload
}

type ControlRequestPayload struct {
	AccessToken string
	Appliance   Appliance
}

type ControlResponse struct {
	Header  Header            `json:"header"`
	Payload map[string]string `json:"payload"`
}

func NewControlResponse(name string) ControlResponse {
	return ControlResponse{
		Header: Header{
			PayloadVersion: "2",
			Namespace:      "Alexa.ConnectedHome.Control",
			Name:           name,
		},
		Payload: make(map[string]string),
	}
}

func NewControlResponseFromRequest(req ControlRequest) (ControlResponse, error) {
	name, err := getConfirmationName(req.Header.Name)
	res := NewControlResponse(name)
	res.Header.MessageId = req.Header.MessageId
	return res, err
}

func getConfirmationName(requestName string) (str string, err error) {
	switch requestName {
	case "TurnOnRequest":
		str = "TurnOnConfirmation"
	case "TurnOffRequest":
		str = "TurnOffConfirmation"
	default:
		err = fmt.Errorf("Could not find corresponding confirmaion name for requestName '%s'", requestName)
	}
	return
}
