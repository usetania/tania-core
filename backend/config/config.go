package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Config Configuration

const (
	DB_INMEMORY = "inmemory"
	DB_SQLITE   = "sqlite"
	DB_MYSQL    = "mysql"
)

type Configuration struct {
	AppPort                *string   `mapstructure:"app_port"`
	DemoMode               *bool     `mapstructure:"demo_mode"`
	UploadPathArea         *string   `mapstructure:"upload_path_area"`
	UploadPathCrop         *string   `mapstructure:"upload_path_crop"`
	TaniaPersistenceEngine *string   `mapstructure:"tania_persistence_engine"`
	SqlitePath             *string   `mapstructure:"sqlite_path"`
	MysqlHost              *string   `mapstructure:"mysql_host"`
	MysqlPort              *string   `mapstructure:"mysql_port"`
	MysqlDbname            *string   `mapstructure:"mysql_dbname"`
	MysqlUsername          *string   `mapstructure:"mysql_username"`
	MysqlPassword          *string   `mapstructure:"mysql_password"`
	RedirectURI            []*string `mapstructure:"redirect_uri"`
	ClientID               *string   `mapstructure:"client_id"`
}

/*
InitViperConfig https://github.com/spf13/viper
Viper uses the following precedence order. Each item takes precedence over the item below it:
- explicit call to Set
- flag
- env
- config
- key/value store
- default
*/
func InitViperConfig() {
	v := viper.New()

	v.AutomaticEnv()

	// App Ports
	pflag.String("app_port", "8080", "Tania server port")

	// Demo Mode
	pflag.Bool("demo_mode", true, "Switch for the demo mode. This will bypass auth check and use hardcoded token demo")

	// Persistence Config
	pflag.String("tania_persistence_engine", "sqlite", "Tania persistence engine. Available engine: mysql, sqlite, inmemory")

	// Persistence Config - SQLite
	pflag.String("sqlite_path", "tania.db", "Path of sqlite file db")

	// Persistence Config - MySQL
	pflag.String("mysql_host", "127.0.0.1", "Mysql Host")
	pflag.String("mysql_port", "3306", "Mysql Port")
	pflag.String("mysql_dbname", "tania", "Mysql DBName")
	pflag.String("mysql_username", "root", "Mysql username")
	pflag.String("mysql_password", "root", "Mysql password")

	// Local Upload Path
	pflag.String("upload_path_area", "tania-uploads/area", "Upload path for the Area photo")
	pflag.String("upload_path_crop", "tania-uploads/crop", "Upload path for the Crop photo")

	// Built-In implicit grant OAuth 2
	pflag.StringSlice("redirect_uri", []string{"http://localhost:8080/oauth2_implicit_callback"}, "URI for redirection after authorization server grants access token")
	pflag.String("client_id", "f0ece679-3f53-463e-b624-73e83049d6ac", "OAuth2 Implicit Grant Client ID for frontend")

	pflag.Parse()
	err := v.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}

	c := Configuration{}

	v.SetConfigType("json")

	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	Config = c

	// Start to read the File Config if any //
	v.SetConfigName("conf")
	v.AddConfigPath("./")

	err = v.ReadInConfig()
	if err != nil {
		log.Warnf("No configuration file found. Err %v", err)
	} else {
		log.Info("Using config file at " + v.ConfigFileUsed())
	}

	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	Config = c
}
