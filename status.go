package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/m0a/easyjson"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// FacebookPlatformPlugin mackerel plugin
type FacebookPlatformApiStatusPlugin struct {
	Prefix string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (u FacebookPlatformApiStatusPlugin) MetricKeyPrefix() string {
	if u.Prefix == "" {
		u.Prefix = "status"
	}
	return u.Prefix
}

// GraphDefinition interface for mackerelplugin
func (u FacebookPlatformApiStatusPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(u.Prefix)
	return map[string](mp.Graphs){
		"": mp.Graphs{
			Label: labelPrefix,
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "health", Label: "Health"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (u FacebookPlatformApiStatusPlugin) FetchMetrics() (map[string]interface{}, error) {
	health, err := FacebookPlatformApiStatus()
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch facebook platform metrics: %s", err)
	}
	return map[string]interface{}{"health": health}, nil
}

func FacebookPlatformApiStatus() (float64, error) {
	resp, err := http.Get("https://www.facebook.com/platform/api-status/")
	if err != nil {
		return 0, fmt.Errorf("Failed to connect facebook.com: %s", err)
	}
	defer resp.Body.Close()

	jsonData, err := easyjson.NewEasyJson(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Invalid responses: %s", err)
	}

	health := jsonData.K("current", "health").AsFloat64Panic()
	return health, nil
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "Facebook Platform Api Status", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u := FacebookPlatformApiStatusPlugin{
		Prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(u)
	helper.Tempfile = *optTempfile
	helper.Run()
}