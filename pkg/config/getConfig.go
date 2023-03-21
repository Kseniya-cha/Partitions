package config

import (
	"flag"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// GetConfig инициализирует и заполняет структуру конфигурационного файла
func GetConfig() (*Config, error) {
	var cfg Config

	// Чтение пути до конфигурационного файла
	configPath, configName, configType := readConfigPath()

	var v = viper.New()

	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)

	err := readParametersFromConfig(v, &cfg)
	if err != nil {
		return &cfg, err
	}

	// Проверка наличия параметров в командной строке
	readFlags(&cfg)
	return &cfg, nil
}

func readParametersFromConfig(v *viper.Viper, cfg *Config) error {
	// Попытка чтения конфига
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	// Попытка заполнение структуры Config полученными данными
	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}

// checkConfigPath проверяет, есть ли среди флагов путь до конфигурационного файла
func readConfigPath() (string, string, string) {
	var absoluteConfigPath, configPath, configName, configType string

	args := os.Args
	for _, arg := range args {
		if strings.Split(arg, "=")[0][1:] == "configPath" {
			absoluteConfigPath = strings.Split(arg, "=")[1]
			configPath, configName, configType = getParamsConf(absoluteConfigPath)
			break
		}
	}
	// Если путь не был указан, выставляется по умолчанию ./
	if absoluteConfigPath == "" {
		configPath = "./"
		configName = "config"
		configType = "yaml"
	}
	return configPath, configName, configType
}

func getParamsConf(absoluteConfigPath string) (string, string, string) {

	pathSplit := strings.Split(absoluteConfigPath, "/")
	configNameType := strings.Split(pathSplit[len(pathSplit)-1], ".")
	configPath := strings.Join(pathSplit[:len(pathSplit)-1], "/") + "/"

	return configPath, configNameType[0], configNameType[1]
}

// readFlags реализует возможность передачи параметров
// конфигурационного файла при запуске из командной строки
func readFlags(cfg *Config) {
	var stub string

	flag.StringVar(&cfg.LogLevel, "logLevel", cfg.LogLevel, "The level of logging parameter")
	flag.BoolVar(&cfg.LogFileEnable, "logFileEnable", cfg.LogFileEnable, "The statement whether to log to a file")
	flag.BoolVar(&cfg.LogStdoutEnable, "logStdoutEnable", cfg.LogStdoutEnable, "The statement whether to log to console")
	flag.StringVar(&cfg.LogFile, "logpath", cfg.LogFile, "The path to file of logging out")
	flag.IntVar(&cfg.MaxSize, "maxSize", cfg.MaxSize, "The path to file of logging out")
	flag.IntVar(&cfg.MaxAge, "maxAge", cfg.MaxAge, "The path to file of logging out")
	flag.IntVar(&cfg.MaxBackups, "maxBackups", cfg.MaxBackups, "The path to file of logging out")
	flag.BoolVar(&cfg.RewriteLog, "rewriteLog", cfg.RewriteLog, "Is rewrite log file")

	flag.IntVar(&cfg.Port, "port", cfg.Port, "The host parameter")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "The host parameter")
	flag.StringVar(&cfg.DbName, "dbName", cfg.DbName, "The db_name parameter")
	flag.StringVar(&cfg.TableName, "tableName", cfg.TableName, "The name of table")
	flag.StringVar(&cfg.User, "user", cfg.User, "The user parameter")
	flag.StringVar(&cfg.Password, "password", cfg.Password, "The password parameter")
	flag.StringVar(&cfg.Driver, "driver", cfg.Driver, "The driver parameter")

	flag.StringVar(&stub, "configPath", `./`, "The path to file of configuration")

	flag.Parse()
}
