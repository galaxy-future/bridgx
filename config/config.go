package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var GlobalConfig *Config

func MustInit() {
	GlobalConfig = loadConfig()
}

func loadConfig() *Config {
	filepath := composeConfigFileName("conf/config.yml", os.Getenv("SpecifiedConfig"))
	log.Printf("config filepath:%s", filepath)

	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var config Config
	if err = yaml.Unmarshal(f, &config); err != nil {
		panic(err)
	}
	return &config
}

func composeConfigFileName(basePath string, suffix string) string {
	var filepath = basePath

	if suffix != "" {
		filepath = strings.Join([]string{filepath, suffix}, ".")
	}

	return filepath
}

type Config struct {
	DebugMode         bool           `yaml:"DebugMode"`
	NeedPublishConfig bool           `yaml:"NeedPublishConfig"`
	ServerPort        int            `yaml:"ServerPort"`
	CostCfg           CostConfig     `yaml:"CostConfig"`
	WriteDB           DBConfig       `yaml:"WriteDB"`
	ReadDB            DBConfig       `yaml:"ReadDB"`
	EtcdConfig        *EtcdConfig    `yaml:"EtcdConfig"`
	JwtToken          JwtTokenConfig `yaml:"JwtToken"`
}

type JwtTokenConfig struct {
	JwtTokenSignKey        string `yaml:"JwtTokenSignKey"`
	JwtTokenCreatedExpires int64  `yaml:"JwtTokenCreatedExpires"`
	JwtTokenRefreshExpires int64  `yaml:"JwtTokenRefreshExpires"`
	BindContextKeyName     string `yaml:"BindContextKeyName"`
}

type DBConfig struct {
	Name         string `yaml:"Name"`
	Host         string `yaml:"Host"`
	Port         string `yaml:"Port"`
	User         string `yaml:"User"`
	Password     string `yaml:"Password"`
	Timeout      string `yaml:"Timeout"`
	ReadTimeout  string `yaml:"ReadTimeout"`
	WriteTimeout string `yaml:"WriteTimeout"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
}

type EtcdConfig struct {
	Endpoints   []string      `yaml:"Endpoints"`
	DailTimeout time.Duration `yaml:"DailTimeout"`
}

type CostConfig struct {
	QueryOrderIntvalSec          int `yaml:"QueryOrderIntvalSec"`
	QueryAlibabaCloudOrderPerMin int `yaml:"QueryAlibabaCloudOrderPerMin"`
}
