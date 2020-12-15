package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
)

func Register() {
	_, err := MakeRequest("POST", "register", buildMetadata())
	if err != nil {
		log.Error(err)
	}
	return
}

func PushDataToController(data structs.ProcessedItemsWrapper) (*io.ReadCloser, int, error) {
	log.Info(fmt.Sprintf("Processing %d item from %s", len(data.Items), os.Getenv("MODULE_SVC_NAME")))
	resp, err := MakeRequest("POST", "queue", data)
	if err != nil {
		log.Error(err)
		return nil, http.StatusInternalServerError, err
	}
	return &resp.Body, resp.StatusCode, nil
}

func MakeRequest(httpMethod, internalEndpoint string, data interface{}) (*http.Response, error) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)

	registerUrl := fmt.Sprintf("http://%s:%s/internal/%s", os.Getenv("CONTROLLER_SVC_NAME"), os.Getenv("CONTROLLER_PORT"), internalEndpoint)

	req, err := http.NewRequest(httpMethod, registerUrl, buf)

	if err != nil {
		return nil, err
	}

	token := os.Getenv("INTERNAL_TOKEN")

	req.Header.Set("x-internal-token", token)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func buildMetadata() structs.ModuleMetadata {
	postMethod := structs.HttpMethod{Method: "POST"}
	getMethod := structs.HttpMethod{Method: "GET"}
	optionsMethod := structs.HttpMethod{Method: "OPTIONS"}

	defaultEndpoint := structs.ModuleEndpoint{
		Secure:      true,
		Endpoint:    "/run",
		HttpMethods: []structs.HttpMethod{optionsMethod, postMethod},
	}

	healthEndpoint := structs.ModuleEndpoint{
		Secure:      true,
		Endpoint:    "/health",
		HttpMethods: []structs.HttpMethod{optionsMethod, getMethod},
	}

	testEndpoint := structs.ModuleEndpoint{
		Secure:      true,
		Endpoint:    "/config",
		HttpMethods: []structs.HttpMethod{optionsMethod, getMethod, postMethod},
	}

	desc := "Ingests intelligence from existing blocklists by importing files over http/https or by uploading files manually."

	localPort := os.Getenv("LOCAL_PORT")
	moduleSvcName := os.Getenv("MODULE_SVC_NAME")

	return structs.ModuleMetadata{
		ModuleServiceName: moduleSvcName,
		ModuleDisplayName: "Generic Data Importer",
		ModuleDescription: desc,
		ModuleType:        "ingress",
		InboundRoute:      "/fpimp",
		InternalIP:        GetLocalIP(),
		InternalPort:      localPort,
		Configured:        true,
		Configurable:      true,
		IconURL:           os.Getenv("ICON_URL"),
		InternalEndpoints: []structs.ModuleEndpoint{defaultEndpoint, healthEndpoint, testEndpoint},
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error(err)
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
