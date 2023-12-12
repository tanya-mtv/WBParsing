package config

import (
	"flag"
	"os"
	"parsingWB/internal/logger"
	"strconv"

	"github.com/spf13/viper"
)

const (
	CONFIG_TYPE = "CONFIG_TYPE"
	CONFIG_PATH = "CONFIG_PATH"
	yaml        = "yaml"
	cfgPath     = "config.yaml"
	APIToken    = "APIToken"
	repInterval = "reportInterval"
	WBURLList   = "WBURLList"
	WBURLParce  = "WBURLParce"
	limit       = "LIMIT"
)

var configPath string
var configType string

type Config struct {
	Logger         *logger.Config `mapstructure:"logger"`
	Token          string
	ReportInterval int
	WBURLList      string `mapstructure:"wbURLList"`
	WBURLParce     string `mapstructure:"wbURLParce"`
	Limit          int64
	MSSQL          *ConfigMSSQL `mapstructure:"mssql"`
}

func init() {
	flag.StringVar(&configType, "config-type", "yaml", "Format of configuration file type. Supported formats is: yaml")
	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(CONFIG_PATH)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			configPath = cfgPath
		}
	}
	cfg := &Config{}

	configType := os.Getenv(CONFIG_TYPE)

	if configType == "" {
		configType = yaml
	}
	viper.SetConfigType(configType)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	token := os.Getenv(APIToken)
	if token != "" {
		cfg.Token = token
	}

	reportInterval, err := strconv.Atoi(os.Getenv(repInterval))
	if err != nil {
		reportInterval = 30
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = reportInterval
	}

	intlimit, err := strconv.ParseInt(os.Getenv(limit), 10, 64)
	if err != nil {
		intlimit = int64(1000)
	}

	if cfg.Limit == 0 {
		cfg.Limit = intlimit
	}

	WBURLList := os.Getenv(WBURLList)
	if cfg.WBURLList == "" {
		cfg.WBURLList = WBURLList
	}

	WBURLParce := os.Getenv(WBURLParce)
	if cfg.WBURLParce == "" {
		cfg.WBURLParce = WBURLParce
	}

	return cfg, nil
}
