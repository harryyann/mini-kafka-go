package config

import "sync"

type KafkaConfig struct {
	BrokerId                  string `yaml:"brokerId"`
	LogDir                    string `yaml:"logDir"`
	Listeners                 string `yaml:"listeners"`
	MaxConnectionsPerListener int    `yaml:"maxConnectionsPerListener,omitempty"`
	MaxMessageBytes           int    `yaml:"maxMessageBytes,omitempty"`
	clientTimeoutSeconds      int    `yaml:"clientTimeoutSeconds,omitempty"`
}

var config *KafkaConfig
var once sync.Once

func NewKafkaConfig() *KafkaConfig {
	once.Do(func() {
		config = &KafkaConfig{}
	})
	return config
}
