package vflags

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/discover"
	"github.com/superwhys/venkit/v2/internal/shared"
	"github.com/superwhys/venkit/v2/snail"
	venkitUtils "github.com/superwhys/venkit/v2/utils"
)

var (
	v                 = viper.New()
	requiredFlags     []string
	nestedKey         = map[string]interface{}{}
	defaultConfigFile = "config.yaml"
	debug             BoolGetter
	config            StringGetter
	useRemoteConfig   BoolGetter
	watchConfig       BoolGetter
	// killWhileChange will kill this service while config change
	// If used in conjunction with the restart configuration of docker,
	// the service can be restarted immediately upon configuration change
	killWhileChange BoolGetter
)

type VflagOption struct {
	autoParseConfig bool
	useConsul       bool
}

type VflagOptionFunc func(*VflagOption)

func Viper() *viper.Viper {
	return v
}

func declareConsulFlags() {
	shared.UseConsul = Bool("useConsul", true, "Whether to use the consul service center.")
	shared.ConsulAddr = String("consulAddr", fmt.Sprintf("%v:8500", discover.HostAddress), "Set the conusl addr.")
	useRemoteConfig = Bool("useRemoteConfig", false, "Set true to use remote config.")
}

func declareDefaultFlags(o *VflagOption) {
	config = StringP("config", "f", defaultConfigFile, "Specify config file. Support json, yaml.")
	debug = Bool("debug", false, "Whether to enable debug mode.")
	shared.ServiceName = StringP("service", "s", os.Getenv("VENKIT_SERVICE"), "Set the service name.")
	watchConfig = Bool("watchConfig", false, "Set true to watch config.")
	killWhileChange = Bool("killWhenChange", false, `It will kill this service while config change. 
If used in conjunction with the restart configuration of docker,
the service can be restarted immediately upon configuration change.`)
	if o.useConsul {
		declareConsulFlags()
	}
}

func OverrideDefaultConfigFile(configFile string) {
	defaultConfigFile = configFile
}

func BindPFlag(key string, flag *pflag.Flag) {
	if err := v.BindPFlag(key, flag); err != nil {
		lg.Fatalf("BindPFlag key: %v, err: %v", key, err)
	}
}

func initVFlags(o *VflagOption) {
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("./tmp/config/")
	v.SetConfigFile("config.yaml")

	declareDefaultFlags(o)
	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		lg.Fatal("BindPFlags error: %v", err)
	}
}

func EnableConsul() VflagOptionFunc {
	return func(vo *VflagOption) {
		vo.useConsul = true
	}
}

// Deprecated: ProhibitConsul Disable consul
func ProhibitConsul() VflagOptionFunc {
	return func(vo *VflagOption) {
		vo.useConsul = false
	}
}

func ProhibitAutoParseConfig() VflagOptionFunc {
	return func(vo *VflagOption) {
		vo.autoParseConfig = false
	}
}

func ConfigFile() ([]byte, error) {
	if config() == "" {
		return nil, errors.New("no config file specify")
	}
	b, err := os.ReadFile(config())
	if err != nil {
		return nil, errors.Wrap(err, "readConf")
	}
	return b, nil
}

func Parse(opts ...VflagOptionFunc) {
	o := &VflagOption{
		autoParseConfig: true,
	}

	for _, opt := range opts {
		opt(o)
	}

	initVFlags(o)
	pflag.Parse()

	injectNestedKey()
	readConfig(o)
	checkFlagKey()
	optionInit()
	snail.Init()
}

func optionInit() {
	if debug() {
		lg.EnableDebug()
	}

	if shared.GetIsUseConsul() {
		discover.SetConsulFinderToDefault()
	}

	if watchConfig() {
		setStructConfWatch()
	}
}

func getServiceNameWithoutTag() string {
	s := GetServiceName()
	segs := strings.SplitN(s, ":", 2)
	if len(segs) < 2 {
		return s
	}
	return segs[0]
}

func getServiceTag() string {
	s := GetServiceName()
	segs := strings.SplitN(s, ":", 2)
	if len(segs) < 2 {
		return ""
	}
	return segs[1]
}

func readConfig(opt *VflagOption) {
	if opt.useConsul && useRemoteConfig() {
		path := readConsulConfig()
		if watchConfig() {
			go watchCnosulConfigChange(path)
		}
		lg.Infoc(lg.Ctx, "Read consul config success. Config=%v", path)
	} else if opt.autoParseConfig && config() != "" && venkitUtils.FileExists(config()) {
		// use local config
		v.SetConfigFile(config())
		if err := v.ReadInConfig(); err != nil {
			lg.Errorc(lg.Ctx, "Read on local file: %v, error: %v", config(), err)
		} else {
			lg.Infoc(lg.Ctx, "Read local config success. Config=%v", config())
		}
	}
}

func checkFlagKey() {
	for _, rk := range requiredFlags {
		if isZero(v.Get(rk)) {
			lg.Fatalf("Missing key: %v", rk)
		}
	}
}

func isZero(i interface{}) bool {
	switch i.(type) {
	case bool:
		// It's trivial to check a bool, since it makes the flag no sense(always true).
		return !i.(bool)
	case string:
		return i.(string) == ""
	case time.Duration:
		return i.(time.Duration) == 0
	case float64:
		return i.(float64) == 0
	case int:
		return i.(int) == 0
	case []string:
		return len(i.([]string)) == 0
	case []interface{}:
		return len(i.([]interface{})) == 0
	default:
		return true
	}
}

func GetServiceName() string {
	return shared.GetServiceName()
}
