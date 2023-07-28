package cfg

import (
	log "github.com/ndmsystems/golog"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var (
	Debug bool
)

// эта обертка нужна, чтобы логировать отсутствие параметра
var config *viper.Viper

func Init(env string) {
	config = viper.New()

	config.AddConfigPath("./config") // path to folder

	config.SetConfigName(env) // (without extension)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	// viper сам конвертирует on/off yes/no в true/false
	Debug = config.GetBool("debug")
}

func GetString(key string) (value string) {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	value = config.GetString(key)
	if value == "" {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	return
}

func GetBool(key string) (result bool) {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	value := config.GetString(key)
	if value == "" {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatal("config: value is not bool: " + key)
	}
	return result
}

func GetInt(key string) int {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	// все как строку берем
	strValue := config.GetString(key)
	if strValue == "" {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	value64, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		log.Errorw("config: invalid type of int", "key", key, "value", strValue)
		return 0
	}
	return int(value64)
}

func GetStringSlice(key string) []string {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	slice := config.GetStringSlice(key)
	if slice == nil || len(slice) == 0 {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	return slice
}

func GetTime(key string) time.Duration {
	exists := config.IsSet(key)
	if !exists {
		log.Fatal("config: value is not set key: " + key)
	}
	strValue := GetString(key)
	if strValue == "" {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	value, err := time.ParseDuration(strValue)
	if err != nil {
		log.Errorw("config: invalid type of time", "key", key, "value", strValue)
		return 0
	}
	return value
}
