package domain

import (
	"crypto/tls"
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/version"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PersistenceProvider defines enum for persistence provider.
type PersistenceProvider string

const (
	// RedisPersistenceProvider uses redis
	RedisPersistenceProvider PersistenceProvider = "REDIS"

	// DynamoDBPersistenceProvider uses DynamoDB
	DynamoDBPersistenceProvider PersistenceProvider = "DYNAMODB"

	// MemoryPersistenceProvider uses in-memory
	MemoryPersistenceProvider PersistenceProvider = "MEMORY"
)

// AuthServiceProvider defines enum for auth service implementation
type AuthServiceProvider string

const (
	// DatabaseAuthServiceProvider uses database based on PersistenceProvider
	DatabaseAuthServiceProvider AuthServiceProvider = "DATABASE"

	// GrpcAuthServiceProvider uses gRPC client based on PersistenceProvider
	GrpcAuthServiceProvider AuthServiceProvider = "GRPC"

	// HttpAuthServiceProvider uses HTTP client based on PersistenceProvider
	HttpAuthServiceProvider AuthServiceProvider = "HTTP"
)

// DynamoDBConfig config
type DynamoDBConfig struct {
	AutoCreateTables    bool   `yaml:"auto_create_tables" mapstructure:"auto_create_tables"`
	TenantPartitionName string `yaml:"tenant_partition_name" mapstructure:"tenant_partition_name"`
	IDName              string `yaml:"id_name" mapstructure:"id_name"`
	ReadCapacityUnits   int64  `yaml:"read_capacity_units" mapstructure:"read_capacity_units"`
	WriteCapacityUnits  int64  `yaml:"write_capacity_units" mapstructure:"write_capacity_units"`
	AWSRegion           string `yaml:"aws_region" mapstructure:"aws_region"`
	Endpoint            string `yaml:"endpoint" mapstructure:"endpoint" env:"DDB_ENDPOINT"`
}

// RedisConfig redis config
type RedisConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	PoolSize int    `yaml:"pool_size" mapstructure:"pool_size"`
}

// Config -- Default Config
type Config struct {
	Redis                      RedisConfig         `yaml:"redis" env:"REDIS"`
	DynamoDB                   DynamoDBConfig      `yaml:"ddb" env:"DYNAMODB"`
	GrpcSasl                   bool                `yaml:"grpc_sasl"`
	GrpcListenPort             string              `yaml:"grpc_listen_port" env:"GRPC_PORT"`
	HttpListenPort             string              `yaml:"http_listen_port" env:"HTTP_PORT"`
	ResourceInstanceExpiration time.Duration       `yaml:"resource_instance_expiration"`
	HttpClientTimeout          time.Duration       `yaml:"http_client_timeout"`
	Debug                      bool                `yaml:"debug"`
	Dir                        string              `yaml:"dir" env:"CONFIG_DIR"`
	PersistenceProvider        PersistenceProvider `yaml:"persistence_provider" env:"PERSISTENCE_PROVIDER"`
	AuthServiceProvider        AuthServiceProvider `yaml:"auth_service_provider" env:"AUTH_SERVICE_PROVIDER"`
	MaxCacheSize               int                 `yaml:"max_cache_size"`
	CacheExpirationMillis      int                 `yaml:"cache_expiration_millis"`
	MaxGroupRoleLevels         int                 `yaml:"max_group_role_levels"`
	ProxyURL                   string              `yaml:"proxy_url"`
	Version                    *version.Info       `yaml:"-"`
}

// NewConfig -- initializes the Default Configuration
func NewConfig(configFile string) (*Config, error) {
	viper.SetDefault("debug", "false")
	viper.SetDefault("dir", "")
	viper.SetDefault("grpc_listen_port", "127.0.0.1:7777")
	viper.SetDefault("http_listen_port", "127.0.0.1:7778")
	viper.SetDefault("max_limit", "10000")
	viper.SetDefault("http_client_timeout", "5s")
	viper.SetDefault("persistence_provider", "REDIS")
	viper.SetDefault("auth_service_provider", "DATABASE")
	viper.SetDefault("ddb.endpoint", "")
	viper.SetDefault("url", "")
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("plexauthz-config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("../config")
		viper.AddConfigPath("../../config")
	}
	var err error
	if err = viper.ReadInConfig(); err == nil {
		logrus.Infof("using config file: %s", viper.ConfigFileUsed())
	} else {
		logrus.Debugf("could not read config %s", err)
	}
	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		logrus.Fatalf("unable to decode into struct, %v", err)
		return nil, err
	}
	confDir := os.Getenv("CONFIG_DIR")
	if confDir != "" {
		config.Dir = confDir
	}
	if err = config.Validate(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) TLSClient() (tlsC TLSConfig, err error) {
	if !c.GrpcSasl {
		return
	}
	var certFile string
	var keyFile string
	var caFile string
	certFile, err = c.ClientCertFile()
	if err != nil {
		return
	}
	keyFile, err = c.ClientKeyFile()
	if err != nil {
		return
	}
	caFile, err = c.CAFile()
	if err != nil {
		return
	}
	tlsC = TLSConfig{
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}
	return
}

func (c *Config) TLSRootClient() (tlsC TLSConfig, err error) {
	if !c.GrpcSasl {
		return
	}
	var certFile string
	var keyFile string
	var caFile string
	certFile, err = c.ClientRootCertFile()
	if err != nil {
		return
	}
	keyFile, err = c.ClientRootKeyFile()
	if err != nil {
		return
	}
	caFile, err = c.CAFile()
	if err != nil {
		return
	}
	tlsC = TLSConfig{
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}
	return
}

func (c *Config) TLSNobodyClient() (tlsC TLSConfig, err error) {
	if !c.GrpcSasl {
		return
	}
	var certFile string
	var keyFile string
	var caFile string
	certFile, err = c.ClientNobodyCertFile()
	if err != nil {
		return
	}
	keyFile, err = c.ClientNobodyKeyFile()
	if err != nil {
		return
	}
	caFile, err = c.CAFile()
	if err != nil {
		return
	}
	tlsC = TLSConfig{
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}
	return
}

func (c *Config) SetupTLSClient() (tlsConfig *tls.Config, err error) {
	tlsC, err := c.TLSClient()
	if err != nil {
		return nil, err
	}
	return tlsC.SetupTLS()
}

func (c *Config) SetupTLSServer(addr string) (tlsConfig *tls.Config, err error) {
	certFile, err := c.ServerCertFile()
	if err != nil {
		return nil, err
	}
	keyFile, err := c.ServerKeyFile()
	if err != nil {
		return nil, err
	}
	caFile, err := c.CAFile()
	if err != nil {
		return nil, err
	}
	tlsC := TLSConfig{
		CertFile:      certFile,
		KeyFile:       keyFile,
		CAFile:        caFile,
		ServerAddress: addr,
		Server:        true,
	}
	return tlsC.SetupTLS()
}

func (c *Config) CAFile() (string, error) {
	return configFile(c.Dir, "ca.pem")
}

func (c *Config) ServerCertFile() (string, error) {
	return configFile(c.Dir, "server.pem")
}

func (c *Config) ServerKeyFile() (string, error) {
	return configFile(c.Dir, "server-key.pem")
}

func (c *Config) ClientCertFile() (string, error) {
	return configFile(c.Dir, "client.pem")
}

func (c *Config) ClientKeyFile() (string, error) {
	return configFile(c.Dir, "client-key.pem")
}

func (c *Config) ClientRootCertFile() (string, error) {
	return configFile(c.Dir, "root-client.pem")
}

func (c *Config) ClientRootKeyFile() (string, error) {
	return configFile(c.Dir, "root-client-key.pem")
}

func (c *Config) ClientNobodyCertFile() (string, error) {
	return configFile(c.Dir, "nobody-client.pem")
}

func (c *Config) ClientNobodyKeyFile() (string, error) {
	return configFile(c.Dir, "nobody-client-key.pem")
}

func (c *Config) ACLModelFile() (string, error) {
	return configFile(c.Dir, "model.conf")
}

func (c *Config) ACLPolicyFile() (string, error) {
	return configFile(c.Dir, "policy.csv")
}

func configFile(dir string, filename string) (string, error) {
	f := filepath.Join(dir, filename)
	if st, err := os.Stat(f); err == nil {
		if st.IsDir() || !st.Mode().IsRegular() {
			return "", fmt.Errorf("%s is directory", f)
		}
		return f, nil
	} else {
		cwd, _ := os.Getwd()
		return "", fmt.Errorf("failed to find '%s' config file [cwd %s] in %s due to %s", f, cwd, dir, err)
	}
}

// Validate ensures config is correct
func (c *Config) Validate() error {
	if err := c.Redis.Validate(); err != nil {
		return err
	}
	if c.ResourceInstanceExpiration.Seconds() <= 0 {
		c.ResourceInstanceExpiration = 15 * time.Minute
	}
	if c.MaxCacheSize <= 0 {
		c.MaxCacheSize = 10000
	}
	if c.CacheExpirationMillis <= 0 {
		c.CacheExpirationMillis = 15000
	}
	if c.MaxGroupRoleLevels <= 0 {
		c.MaxGroupRoleLevels = 5
	}

	if c.PersistenceProvider == "" {
		c.PersistenceProvider = "REDIS"
	}
	if c.GrpcListenPort == "" {
		c.GrpcListenPort = "127.0.0.1:7777"
	}
	if c.HttpListenPort == "" {
		c.HttpListenPort = ":7778"
	}
	if c.HttpClientTimeout.Seconds() < 0 {
		c.HttpClientTimeout = time.Second * 5
	}
	return nil
}

// Validate - validates
func (c *RedisConfig) Validate() error {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 6379
	}
	return nil
}

// Validate - validates
func (c *DynamoDBConfig) Validate() error {
	if c.TenantPartitionName == "" {
		c.TenantPartitionName = "tenant"
	}
	if c.IDName == "" {
		c.IDName = "id"
	}
	if c.ReadCapacityUnits <= 0 {
		c.ReadCapacityUnits = 5
	}
	if c.WriteCapacityUnits <= 0 {
		c.WriteCapacityUnits = 5
	}
	if c.AWSRegion == "" {
		if os.Getenv("AWS_DEFAULT_REGION") != "" {
			c.AWSRegion = os.Getenv("AWS_DEFAULT_REGION")
		} else {
			c.AWSRegion = "us-west-1"
		}
	}
	if c.Endpoint == "" && os.Getenv("DDB_ENDPOINT") != "" {
		c.Endpoint = os.Getenv("DDB_ENDPOINT")
	}

	return nil
}
