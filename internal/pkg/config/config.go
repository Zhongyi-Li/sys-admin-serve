package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Logger LoggerConfig `mapstructure:"logger"`
	Casbin CasbinConfig `mapstructure:"casbin"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	ParseTime    bool   `mapstructure:"parse_time"`
	Loc          string `mapstructure:"loc"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type CasbinConfig struct {
	ModelPath  string `mapstructure:"model_path"`
	PolicyPath string `mapstructure:"policy_path"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.setDefaults()

	return &cfg, nil
}

func (c *Config) setDefaults() {
	if c.App.Name == "" {
		c.App.Name = "sys-admin-serve"
	}

	if c.App.Env == "" {
		c.App.Env = "local"
	}

	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}

	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}

	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}

	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}

	if c.Logger.Filename == "" {
		c.Logger.Filename = "logs/app.log"
	}

	if c.Logger.MaxSize == 0 {
		c.Logger.MaxSize = 100
	}

	if c.Logger.MaxBackups == 0 {
		c.Logger.MaxBackups = 7
	}

	if c.Logger.MaxAge == 0 {
		c.Logger.MaxAge = 30
	}
}

func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c MySQLConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.ParseTime,
		url.QueryEscape(c.Loc),
	)
}

func (c MySQLConfig) AdminDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?charset=%s&parseTime=%t&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Charset,
		c.ParseTime,
		url.QueryEscape(c.Loc),
	)
}

func (c RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
