package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RiotAPIKey string `mapstructure:"RIOT_API_KEY"`
	RiotSleep  int    `mapstructure:"RIOT_SLEEP"`
	MongoURI   string `mapstructure:"MONGO_URI"`
	MongoDB    string `mapstructure:"MONGO_DB_NAME"`
	Workers    int    `mapstructure:"WORKERS"`
}

func LoadConfig() *Config {
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &Config{
		RiotAPIKey: viper.GetString("RIOT_API_KEY"),
		RiotSleep:  viper.GetInt("RIOT_SLEEP"),
		MongoURI:   viper.GetString("MONGO_URI"),
		MongoDB:    viper.GetString("MONGO_DB_NAME"),
		Workers:    viper.GetInt("WORKERS"),
	}
}
