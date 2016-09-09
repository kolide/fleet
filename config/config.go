package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix = "KOLIDE"
)

type AppConfig struct {
	Mysql struct {
		Address  string
		Username string
		Password string
		Database string
	}
	Server struct {
		Address string
		Cert    string
		Key     string
	}
	Auth struct {
		JwtKey      string
		BcryptCost  int
		SaltKeySize int
	}
}

// LoadConfig will load the config variables into a fully initialized AppConfig struct
func (man ConfigManager) LoadConfig() AppConfig {
	var config AppConfig

	// MySQL
	config.Mysql.Address = man.GetConfigString("mysql.address")
	config.Mysql.Username = man.GetConfigString("mysql.username")
	config.Mysql.Password = man.GetConfigString("mysql.password")
	config.Mysql.Database = man.GetConfigString("mysql.database")

	// Server
	config.Server.Address = man.GetConfigString("server.address")
	config.Server.Cert = man.GetConfigString("server.cert")
	config.Server.Key = man.GetConfigString("server.key")

	// Auth
	config.Auth.JwtKey = man.GetConfigString("auth.jwt_key")
	config.Auth.BcryptCost = man.GetConfigInt("auth.bcrypt_cost")
	config.Auth.SaltKeySize = man.GetConfigInt("auth.salt_key_size")

	return config
}

func (man ConfigManager) AttachConfigs() {
	// MySQL
	man.AddConfigString("mysql.address", "localhost:3306")
	man.AddConfigString("mysql.username", "kolide")
	man.AddConfigString("mysql.password", "kolide")
	man.AddConfigString("mysql.database", "kolide")

	// Server
	man.AddConfigString("server.address", "localhost:8080")
	man.AddConfigString("server.cert", "./tools/osquery/kolide.crt")
	man.AddConfigString("server.key", "./tools/osquery/kolide.key")

	// Auth
	man.AddConfigString("auth.jwt_key", "CHANGEME")
	man.AddConfigInt("auth.bcrypt_cost", 12)
	man.AddConfigInt("auth.salt_key_size", 24)
}

var (
	// File may or may not contain the path to the config file
	File string
)

func init() {
	cobra.OnInitialize(initConfig)
}

// Due to a deficiency in viper (https://github.com/spf13/viper/issues/71), one
// can not set the default values of nested config elements. For example, if the
// "mysql" section of the config allows a user to define "username", "password",
// and "database", but the only wants to override the default for "username".
// they should be able to create a config which looks like:
//
//   mysql:
//     username: foobar
//
// In viper, that would nullify the default values of all other config keys in
// the mysql section ("mysql.*"). To get around this, instead of using the
// provided API for setting default values, after we've read the config and env,
// we manually check to see if the value has been set and, if it hasn't, we set
// it manually.
func setDefaultConfigValue(key string, value interface{}) {
	if viper.Get(key) == nil {
		viper.Set(key, value)
	}
}

func recurseConfig(config AppConfig) {
	refType := reflect.TypeOf(&config).Elem()
	numFields := refType.NumField()
	for i := 0; i < numFields; i++ {
		field := refType.Field(i)
		recurseConfigValue(func(leaf reflect.StructField, prefix string) { reflect.ValueOf(leaf).SetString(prefix) }, field, "")
	}
}

/*
func recurseConfig2(config AppConfig) {
	refType := reflect.TypeOf(&config).Elem()
	refVal := reflect.ValueOf(&config).Elem()
	numFields := refType.NumField()
	for i := 0; i < numFields; i++ {
		field := refType.Field(i)
		recurseConfigValue2(func(leaf reflect.StructField, prefix string) { reflect.ValueOf(leaf).SetString(prefix) }, field, "")
	}
}
*/

type structLeafFunc func(leaf reflect.StructField, prefix string)

/*
func recurseConfigValue2(fun structLeafFunc, field reflect.StructField, val reflect.Value, prefix string) {
	switch root.Type().Kind() {
	case reflect.Struct:
		fmt.Println("Got struct")
		tag := root.Tag.Get("config")
		numFields := root.Type.NumField()
		for i := 0; i < numFields; i++ {
			field := root.Type.Field(i)
			recurseConfigValue(fun, field, prefix+tag+".")
		}

	}
}
*/

func recurseConfigValue(fun structLeafFunc, root reflect.StructField, prefix string) {

	switch root.Type.Kind() {
	case reflect.Struct:
		fmt.Println("Got struct")
		tag := root.Tag.Get("config")
		numFields := root.Type.NumField()
		for i := 0; i < numFields; i++ {
			field := root.Type.Field(i)
			recurseConfigValue(fun, field, prefix+tag+".")
		}
	default:
		fmt.Println("Got other: ", root.Type.Kind(), prefix+root.Tag.Get("config"))
		fun(root, prefix)
	}
}

type ConfigReader struct {
	EnvKeyReplacer *strings.Replacer
}

func NewConfigReader() ConfigReader {
	return ConfigReader{
		EnvKeyReplacer: strings.NewReplacer(".", "_"),
	}
}

func envVarNameFromConfigKey(key string) string {
	return envPrefix + "_" + strings.ToUpper(strings.Replace(key, ".", "_", -1))
}

type ConfigManager struct {
	Command  *cobra.Command
	defaults map[string]interface{}
}

// AddConfigString Adds a string config to the config options
func (man ConfigManager) AddConfigString(key string, defVal string) {
	man.Command.PersistentFlags().String(key, defVal, "Env: "+envVarNameFromConfigKey(key))
	viper.BindPFlag(key, man.Command.PersistentFlags().Lookup(key))
	viper.BindEnv(key, envVarNameFromConfigKey(key))
	// Look up default from defaults map
	// if defaults.
}

// GetConfigString Retrieves a string from the loaded config
func (man ConfigManager) GetConfigString(key string) string {
	if viper.Get(key) == nil {
		// Flag did not appear in config, try to use default
		flag := man.Command.PersistentFlags().Lookup(key)
		if flag == nil {
			panic("Tried to look up default value for nonexistent config option " + key)
		}
		return flag.DefValue
	}
	return viper.GetString(key)
}

// AddConfigInt Adds an int config to the config options
func (man ConfigManager) AddConfigInt(key string, defVal int) {
	man.Command.PersistentFlags().Int(key, defVal, "Env: "+envVarNameFromConfigKey(key))
	viper.BindPFlag(key, man.Command.PersistentFlags().Lookup(key))
	viper.BindEnv(key, envVarNameFromConfigKey(key))
}

// GetConfigString Retrieves a string from the loaded config
func (man ConfigManager) GetConfigInt(key string) int {
	if viper.Get(key) == nil {
		// Flag did not appear in config, try to use default
		flag := man.Command.PersistentFlags().Lookup(key)
		if flag == nil {
			panic("Tried to look up default value for nonexistent config option " + key)
		}
		return 0 // TODO return actual defaults value
	}
	return viper.GetInt(key)
}

func initConfig() {
	if File != "" {
		viper.SetConfigFile(File)
	} else {
		viper.SetConfigName("kolide")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath("./tools/app")
		viper.AddConfigPath("/etc/kolide")
	}

	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Error reading config file: %s", viper.ConfigFileUsed())
	}

	logrus.Info("Using config file: ", viper.ConfigFileUsed())

	setDefaultConfigValue("mysql.address", "foo:3306")
	setDefaultConfigValue("mysql.username", "kolide")
	setDefaultConfigValue("mysql.password", "kolide")
	setDefaultConfigValue("mysql.database", "kolide")

	setDefaultConfigValue("server.address", "0.0.0.0:8080")

	setDefaultConfigValue("app.web_address", "0.0.0.0:8080")

	setDefaultConfigValue("auth.jwt_key", "CHANGEME")
	setDefaultConfigValue("auth.bcrypt_cost", 12)
	setDefaultConfigValue("auth.salt_key_size", 24)

	setDefaultConfigValue("smtp.token_key_size", 24)
	setDefaultConfigValue("smtp.address", "localhost:1025")
	setDefaultConfigValue("smtp.pool_connections", 4)

	setDefaultConfigValue("session.key_size", 64)
	setDefaultConfigValue("session.expiration_seconds", 60*60*24*90)
	setDefaultConfigValue("session.cookie_name", "KolideSession")

	setDefaultConfigValue("osquery.node_key_size", 24)
	setDefaultConfigValue("osquery.status_log_file", "/tmp/osquery_status")
	setDefaultConfigValue("osquery.result_log_file", "/tmp/osquery_result")
	setDefaultConfigValue("osquery.label_up_interval", 1*time.Minute)

	setDefaultConfigValue("logging.debug", false)
	setDefaultConfigValue("logging.disable_banner", false)

	if viper.GetBool("logging.debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	if viper.GetBool("logs.json") {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
