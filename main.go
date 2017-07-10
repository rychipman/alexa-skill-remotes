package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/rychipman/alexa-skill-remotes/alexa"
)

func main() {
	http.HandleFunc("/authorization", handleAuthorization)
	http.HandleFunc("/access-token", handleAccessToken)
	http.HandleFunc("/action", handleAction)
	http.HandleFunc("/action/discover", handleDiscover)
	http.HandleFunc("/action/control", handleControl)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAction(w http.ResponseWriter, req *http.Request) {
	log.Printf("received action with unknown namespac")
}

func handleDiscover(w http.ResponseWriter, req *http.Request) {
	action := alexa.DiscoverRequest{}
	err := json.NewDecoder(req.Body).Decode(&action)
	if err != nil {
		log.Printf("got error while decoding: %s", err.Error())
	}
	log.Printf("received discover action: %#v", action)
	response := alexa.DiscoverResponse{
		alexa.Header{
			PayloadVersion: "2",
			Namespace:      "Alexa.ConnectedHome.Discovery",
			Name:           "DiscoverAppliancesResponse",
		},
		alexa.DiscoverResponsePayload{
			[]alexa.Appliance{
				alexa.Appliance{
					Id:          "ryan_air_conditioner",
					Name:        "Air Conditioner",
					Description: "Sarah's Living Room Air Conditioner",
					IsReachable: true,
					Actions:     []string{"turnOn", "turnOff"},

					ManufacturerName: "frigidaire",
					ModelName:        "mine",
					Version:          "1.0",
					Details:          make(map[string]string),
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleControl(w http.ResponseWriter, req *http.Request) {
	action := alexa.ControlRequest{}
	err := json.NewDecoder(req.Body).Decode(&action)
	if err != nil {
		log.Printf("got error while decoding: %s", err.Error())
		return
	}
	log.Printf("received control action '%s' for appliance '%s", action.Header.Name, action.Payload.Appliance.Id)

	device, key, err := controlRequestToLIRC(action)
	if err != nil {
		log.Printf("got error while building lirc command: %s", err.Error())
		return
	}
	sendRemoteKey(device, key)

	response, _ := alexa.NewControlResponseFromRequest(action)
	json.NewEncoder(w).Encode(response)
}

func handleAccessToken(w http.ResponseWriter, req *http.Request) {
	response := make(map[string]string)
	response["access_token"] = "i_am_an_access_token"
	response["refresh_token"] = "i_am_a_refresh_token"
	json.NewEncoder(w).Encode(response)
}

func handleAuthorization(w http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	state := params.Get("state")
	redirect_uri := params.Get("redirect_uri")

	code := "123456789ABC"

	redirect_uri_full := fmt.Sprintf("%s?state=%s&code=%s", redirect_uri, state, code)

	log.Printf("redirecting to uri: '%s'", redirect_uri_full)

	http.Redirect(w, req, redirect_uri_full, 302)
}

func sendRemoteKey(device, key string) error {
	log.Printf("Sending key '%s' to device '%s'", key, device)
	cmd := exec.Command("irsend", "send_once", device, key)
	//cmd := exec.Command("notify-send", fmt.Sprintf("sending key '%s' to device '%s", key, device))
	err := cmd.Run()
	return err
}

func controlRequestToLIRC(req alexa.ControlRequest) (device, key string, err error) {
	switch req.Header.Name {
	case "TurnOnRequest":
		key = "KEY_ON"
	case "TurnOffRequest":
		key = "KEY_OFF"
	default:
		err = fmt.Errorf("could not find lirc key for request type '%s'", req.Header.Name)
	}

	switch req.Payload.Appliance.Id {
	case "ryan_air_conditioner":
		device = "ac"
	default:
		err = fmt.Errorf("could not find lirc key for appliance '%s'", req.Payload.Appliance.Id)
	}

	return
}
