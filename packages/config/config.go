package config

import (
	"strconv"
	"time"

	log "github.com/ndmsystems/golog"
	"github.com/spf13/viper"
)

var (
	Debug bool   // режим отладки
	env   string // текущее окружение
)

// эта обертка нужна, чтобы логировать отсутствие параметра
var config *viper.Viper

func Init(e string) {
	config = viper.New()
	env = e
	config.AddConfigPath("./config") // path to folder

	config.SetConfigName(env) // (without extension)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
}

func Get(key string) (value string) {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	value = config.GetString(key)
	if value == "" {
		log.Warningw("config: value is empty, key: "+key, "key", key)
	}
	log.Debug(key, value)
	return
}

func Bool(key string) (result bool) {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	value := config.GetString(key)
	if value == "" {
		log.Warningw("config: value is empty, key: "+key, "key", key)
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatal("config: value is not bool: " + key)
	}
	return result
}

func Int(key string) int {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	// все как строку берем
	strValue := config.GetString(key)
	if strValue == "" {
		log.Warningw("config: value is empty, key: "+key, "key", key)
	}
	value64, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		log.Fatalw("config: invalid type of int", "key", key, "value", strValue)
		return 0
	}
	return int(value64)
}

func List(key string) []string {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	slice := config.GetStringSlice(key)
	if slice == nil || len(slice) == 0 {
		log.Warningw("config: value is empty, key: "+key, "key", key)
	}
	return slice
}

func Map(key string) map[string]interface{} {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	m := config.GetStringMap(key)
	if m == nil || len(m) == 0 {
		log.Warningw("config: value is empty, key: "+key, "key", key)
	}

	return m
}

func Time(key string) time.Duration {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	strValue := Get(key)
	if strValue == "" {
		log.Warningw("config: value is empty, key: "+key, "key", key)
	}
	value, err := time.ParseDuration(strValue)
	if err != nil {
		log.Fatal("config: invalid type of time", "key", key, "value", strValue)
	}
	return value
}

// Info returns info about device
func Info() map[string]interface{} {
	info := Map("info")
	info["ts"] = time.Now().Unix()
	return info
}

// IsProd returns true if env is prod
func IsProd() bool {
	return env == "prod"
}

// IsKeenetic if keenetic device
func IsKeenetic() bool {
	return env == "keenetic"
}
