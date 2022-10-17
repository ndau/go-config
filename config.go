package config

import (
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/fsnotify/fsnotify"
	logger "github.com/ndau/go-logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Interface defines getters for any config data type
type Config interface {
	BindPFlag(key string, flag *pflag.Flag) error
	BindPFlags(flags *pflag.FlagSet) error
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetIntSlice(key string) []int
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	IsSet(key string) bool
	AllSettings() map[string]interface{}
}

// Returns a type compliant with the Config interface
func New(configFiles ...string) (*viper.Viper, error) {
	conf := viper.New()
	replacer := strings.NewReplacer(".", "_")
	conf.SetEnvKeyReplacer(replacer)
	conf.SetEnvPrefix("NDAU")
	conf.AutomaticEnv()

	logLevel := os.Getenv("NDAU_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	log, err := logger.New("config", logLevel)
	if err != nil {
		return nil, err
	}

	watchConfig := false

	// Extract information from params and env
	switch len(configFiles) {
	case 0:
		// Use sensible defaults for config name and path
		configName := os.Getenv("NDAU_CONFIG_NAME")
		configPath := os.Getenv("NDAU_CONFIG_PATH")

		log.Infof("Load config from: %s/%s", configPath, configName)
		if configName == "" || configPath == "" {
			break
		}

		conf.SetConfigName(configName)
		conf.AddConfigPath(configPath)
		err := backoff.Retry(func() error {
			var err error
			err = conf.ReadInConfig()
			if err != nil {
				log.Info(err)
			}
			return err
		}, backoff.NewConstantBackOff(2*time.Second))
		if err != nil {
			return nil, err
		}

		watchConfig = true
	default:
		for i, path := range configFiles {
			conf.SetConfigFile(path)
			err := backoff.Retry(func() error {
				var err error
				if i == 0 {
					err = conf.ReadInConfig()
				} else {
					err = conf.MergeInConfig()
				}
				if err != nil {
					log.Info(err)
				}
				return err
			}, backoff.NewConstantBackOff(2*time.Second))
			if err != nil {
				return nil, err
			}
		}
		watchConfig = true
	}

	if watchConfig {
		conf.OnConfigChange(func(e fsnotify.Event) {
			log.Infof("Reloaded config: %s\n", e.Name)
		})
		conf.WatchConfig()
	}

	return conf, nil
}
