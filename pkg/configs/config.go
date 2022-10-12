package configs

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GrpcPort       string               `yaml:"grpc_port"`
	HttpPort       string               `yaml:"http_port"`
	MongoDb        MongoDB              `yaml:"mongo_db"`
	RedisCache     Redis                `yaml:"redis_cache"`
	KeyJwt         string               `yaml:"key_jwt"`
	GrpcClientConn GrpcClientConnConfig `yaml:"grpc_conn"`
	Exp            time.Duration        `yaml:"exp"`
	TotpSecret     string               `yaml:"totp_secret"`
	TimeoutRedis   time.Duration        `yaml:"time_out_redis"`
	TimeRequestId  time.Duration        `yaml:"time_request_id"`
	TimeEmailOtp   time.Duration        `yaml:"time_email_otp"`
}

type MongoDB struct {
	Url        string `yaml:"url"`
	DbName     string `yaml:"db_name"`
	Collection string `yaml:"collection"`
}

type Redis struct {
	Address string `yaml:"address"`
	Url     string `yaml:"url"`
}

type GrpcClientConnConfig struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	AccessToken string        `mapstructure:"access_token"`
}

func Get(this interface{}, key string) interface{} {
	return this.(map[interface{}]interface{})[key]
}
func String(payload interface{}) string {
	var load string
	if pay, oh := payload.(string); oh {
		load = pay
	} else {
		load = ""
	}
	return load
}
func LoadConfig() (*Config, error) {
	filename, err := filepath.Abs("../pkg/configs/config.yaml")
	if err != nil {
		return nil, errors.WithMessage(err, "filepath.Abs")
	}
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.WithMessage(err, "ioutil.ReadFile")
	}
	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, errors.WithMessage(err, "yaml.Unmarshal")
	}
	return &config, nil
}
