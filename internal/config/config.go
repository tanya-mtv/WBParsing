package config

import (
	"flag"
	"os"
	"parsingWB/internal/logger"
	"strconv"

	"github.com/spf13/viper"
)

const (
	CONFIG_TYPE    = "CONFIG_TYPE"
	CONFIG_PATH    = "CONFIG_PATH"
	yaml           = "yaml"
	cfgPath        = "config.yaml"
	APIToken       = "WB_API_TOKEN"
	repInterval    = "reportInterval"
	wb_list_url    = "WB_LIST_URL"
	wb_catalog_url = "WB_CATALOG_URL"
	limit          = "LIMIT"
)

var configPath string
var configType string

type Config struct {
	Logger         *logger.Config `mapstructure:"logger"`
	HromeDriver    string         `mapstructure:"hromedriver"`
	HromePort      int            `mapstructure:"hromeport"`
	Token          string
	ReportInterval int
	WbListUrl      string `mapstructure:"wb_list_url"`
	WbCatalogUrl   string `mapstructure:"wb_catalog_url"`
	Limit          int64
	MSSQL          *ConfigMSSQL `mapstructure:"mssql"`
}

func init() {
	flag.StringVar(&configType, "config-type", "", "Format of configuration file type. Supported formats are: JSON, TOML, YAML, HCL, envfile and Java properties config files")
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
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

	if configType == "" {
		configTypeFromEnv := os.Getenv(CONFIG_TYPE)
		if configTypeFromEnv != "" {
			configType = configTypeFromEnv
		} else {
			configType = yaml
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

	wb_list_url := os.Getenv(wb_list_url)
	if cfg.WbListUrl == "" {
		cfg.WbListUrl = wb_list_url
	}

	wb_catalog_url := os.Getenv(wb_catalog_url)
	if cfg.WbCatalogUrl == "" {
		cfg.WbCatalogUrl = wb_catalog_url
	}

	return cfg, nil
}
