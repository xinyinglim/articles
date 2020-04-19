// The Cloud Function that will update Google Cloud IoT Core device configuration
// To deploy the function, run the following command:
//
// gcloud functions deploy UpdateWeather --runtime go111 --trigger-http --allow-unauthenticated

package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	cloudiot "google.golang.org/api/cloudiot/v1"
)

type FanConfig struct {
	On    bool `json:"on"`
	Speed int  `json:"speed"`
}

func UpdateWeather(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	service, err := cloudiot.NewService(ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	deviceService := cloudiot.NewProjectsLocationsRegistriesDevicesService(service)
	sunnyConfigData := FanConfig{
		On:    true,
		Speed: 20,
	}
	bytes, err := json.Marshal(sunnyConfigData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	encodedString := b64.StdEncoding.EncodeToString(bytes)
	configRequest := cloudiot.ModifyCloudToDeviceConfigRequest{
		BinaryData: encodedString,
	}

	projectID := "YOUR-GCP-PROJECT-ID"
	location := "REGISTRY-LOCATION"
	registryID := "REGISTRY-ID"
	deviceID := "DEVICE-ID"

	devicePath := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, location, registryID, deviceID)
	call := deviceService.ModifyCloudToDeviceConfig(devicePath, &configRequest)
	call.Context(ctx)

	_, err = call.Do()
	if err != nil {
		fmt.Printf("%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success!"))
}
