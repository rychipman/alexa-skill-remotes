package alexa

type Header struct {
	MessageId      string `json:"messageId,omitempty"`
	PayloadVersion string `json:"payloadVersion"`
	Namespace      string `json:"namespace"`
	Name           string `json:"name"`
}
