package conf

import (
	"github.com/spf13/viper"
)

var (
	// Addr is websocket server address
	Addr string
	// Pub stores config related to publisher
	Pub struct {
		Turns   int
		Workers int
	}
	// Lis stores config related to listener
	Lis struct {
		Workers int
	}
)

// Load uses viper to load config from .env.yml in the root dir
func Load() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.SetConfigName(".env")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	// fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	Addr = viper.GetString("addr")
	Pub.Turns = viper.GetInt("pub.turns")
	Pub.Workers = viper.GetInt("pub.workers")
	Lis.Workers = viper.GetInt("lis.workers")
	return nil
}
