package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Component struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

type Info struct {
	Components []Component `json:"components"`
}

type ServerInfo struct {
	Info Info `json:"info"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	hyperionUrl := os.Getenv("HYPERION_JSON_RPC_URL")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		isEnabled := isHyperionLedDeviceEnabled(hyperionUrl)
		enableDisable(hyperionUrl, isEnabled)
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "Hyperion enabled: " + strconv.FormatBool(!isEnabled)})
	})

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./templates/favicon.ico")
	})

	router.Run()
}

func isHyperionLedDeviceEnabled(hyperionUrl string) bool {
	jsonBody := []byte(`{"command":"serverinfo"}`)
	res := sendRequest(hyperionUrl, jsonBody)
	resBody, _ := io.ReadAll(res.Body)

	var serverInfo ServerInfo
	json.Unmarshal(resBody, &serverInfo)

	for _, component := range serverInfo.Info.Components {
		if component.Name == "LEDDEVICE" {
			return component.Enabled
		}
	}

	return true
}

func enableDisable(hyperionUrl string, isEnabled bool) {
	for _, component := range []string{"V4L", "LEDDEVICE"} {
		body := []byte(fmt.Sprintf(`{ "command":"componentstate", "componentstate":{ "component":"%s","state": %t } }`, component, !isEnabled))
		sendRequest(hyperionUrl, body)
	}
}

func sendRequest(hyperionUrl string, body []byte) http.Response {
	reader := bytes.NewReader(body)
	request, _ := http.NewRequest(http.MethodPost, hyperionUrl+"/json-rpc", reader)
	res, _ := http.DefaultClient.Do(request)

	return *res
}
