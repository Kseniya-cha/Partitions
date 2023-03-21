package config

// Config - структура конфига
type Config struct {
	Logger   `mapstructure:"logger"`
	Database `mapstructure:"database"`
}

// Logger содержит параметры логгера
type Logger struct {
	LogLevel        string `mapstructure:"logLevel"`
	LogFileEnable   bool   `mapstructure:"logFileEnable"`
	LogStdoutEnable bool   `mapstructure:"logStdoutEnable"`
	LogFile         string `mapstructure:"logpath"`
	MaxSize         int    `mapstructure:"maxSize"`
	MaxAge          int    `mapstructure:"maxAge"`
	MaxBackups      int    `mapstructure:"maxBackups"`
	RewriteLog      bool   `mapstructure:"rewriteLog"`
}

// Database содержит параметры базы данных
type Database struct {
	Port      int    `mapstructure:"port"`
	Host      string `mapstructure:"host"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DbName    string `mapstructure:"dbName"`
	TableName string `mapstructure:"tableName"`
	Driver    string `mapstructure:"driver"`
}
