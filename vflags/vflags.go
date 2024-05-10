package vflags

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/superwhys/venkit/discover"
	"github.com/superwhys/venkit/internal/shared"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/snail"
)

var (
	v                 = viper.New()
	requiredFlags     []string
	nestedKey         = map[string]interface{}{}
	defaultConfigFile string
	debug             BoolGetter
	config            StringGetter
)

type VflagOption struct {
	autoParseConfig bool
}

type VflagOptionFunc func(*VflagOption)

func Viper() *viper.Viper {
	return v
}

func OverrideDefaultConfigFile(configFile string) {
	defaultConfigFile = configFile
}

func BindPFlag(key string, flag *pflag.Flag) {
	if err := v.BindPFlag(key, flag); err != nil {
		lg.Fatalf("BindPFlag key: %v, err: %v", key, err)
	}
}

func initVFlags() {
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("./tmp/config/")

	declareDefaultFlags()
	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		lg.Fatal("BindPFlags error: %v", err)
	}
}

func ProhibitConsul() VflagOptionFunc {
	return func(_ *VflagOption) {
		shared.UseConsul = func() bool {
			return false
		}
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

	initVFlags()
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
}

func declareDefaultFlags() {
	config = StringP("config", "f", defaultConfigFile, "Specify config file. Support json, yaml")
	debug = Bool("debug", false, "Whether to enable debug mode")
	shared.ServiceName = StringP("service", "s", os.Getenv("VENKIT-SERVICE"), "Set the service name")
	if shared.UseConsul == nil || shared.UseConsul() {
		shared.ConsulAddr = String("consulAddr", fmt.Sprintf("%v:8500", discover.HostAddress), "Set the conusl addr")
		shared.UseConsul = Bool("useConsul", true, "Whether to use the consul service center")
	}
}

func readConfig(opt *VflagOption) {
	if opt.autoParseConfig && config() != "" {
		v.SetConfigFile(config())
		if err := v.ReadInConfig(); err != nil {
			lg.Errorf("Read on local file: %v, error: %v", config(), err)
		} else {
			lg.Infof("Read config from local file: %v!", config())
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
