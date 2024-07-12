package vflags

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/internal/shared"
)

func watchCnosulConfigChange(path string) {
	plan, err := watch.Parse(map[string]interface{}{"type": "key", "key": path})
	lg.PanicError(err)

	first := true
	var currentVal []byte

	plan.Handler = func(u uint64, data interface{}) {
		kvPair, ok := data.(*api.KVPair)
		if !ok {
			lg.Warnc(lg.Ctx, "Failed to watch remote config.")
			return
		}

		if first {
			first = false
			currentVal = kvPair.Value
			return
		}

		if string(kvPair.Value) == string(currentVal) {
			lg.Debugc(lg.Ctx, "Remote config not change.")
			return
		}
		lg.Debugc(lg.Ctx, "Remote config change")

		if err := v.ReadConfig(bytes.NewBuffer(kvPair.Value)); err != nil {
			lg.Errorc(lg.Ctx, "viper read config error: %v", err)
		}

		if killWhileChange() {
			killToRestartServer()
			return
		}
		structConfReload()
	}

	plan.Run(shared.GetConsulAddress())
}

func readConsulConfig() string {
	// use remote config
	serviceName := getServiceNameWithoutTag()
	serviceTag := getServiceTag()
	path := fmt.Sprintf("/etc/configs/%v/%v.yaml", serviceName, serviceTag)
	v.AddRemoteProvider("consul", shared.GetConsulAddress(), path)
	v.SetConfigType("yaml")
	if err := v.ReadRemoteConfig(); err != nil {
		lg.Errorc(lg.Ctx, "Fetch remote config -> %v:%v, error: %v", serviceName, serviceTag, err)
		lg.Fatal("Failed to read remote config.")
	}
	return path
}

// killToRestartServer will kill the server first.
// If the service runs with docker and is set to start automatically,
// it can implement configuration updates and refresh the service
func killToRestartServer() {
	lg.Infoc(lg.Ctx, "Remote config changed. Shutting down...")
	kill()
}
