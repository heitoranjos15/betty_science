package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RiotAPIKey string `mapstructure:"RIOT_API_KEY"`
	MongoURI   string `mapstructure:"MONGO_URI"`
	MongoDB    string `mapstructure:"MONGO_DB_NAME"`
}

func LoadConfig() *Config {
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &Config{
		RiotAPIKey: viper.GetString("RIOT_API_KEY"),
		MongoURI:   viper.GetString("MONGO_URI"),
		MongoDB:    viper.GetString("MONGO_DB_NAME"),
	}
}
