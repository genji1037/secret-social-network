package config

import log "github.com/sirupsen/logrus"

// Serve is server config.
type Serve struct {
	Host      string       `yaml:"host"`      // 主机地址
	Port      int          `yaml:"port"`      // 端口号
	Level     Level        `yaml:"level"`     // 日志级别
	MySQL     MySQL        `yaml:"mysql"`     // MySQL数据库
	DGraph    DGraph       `yaml:"dgraph"`    // dgraph数据库
	Open      OpenPlatform `yaml:"open"`      // 开放平台
	Consensus Consensus    `yaml:"consensus"` // 共识模块
}

// Level is log level.
type Level string

// Value get log level value.
func (level Level) Value() log.Level {
	lvl, err := log.ParseLevel(string(level))
	if err != nil {
		return log.InfoLevel
	}
	return lvl
}

// MySQL config.
type MySQL struct {
	Host         string `yaml:"host"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	AlwayMigrate bool   `yaml:"alway_migrate"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// DGraph config.
type DGraph struct {
	Addr string `yaml:"addr"`
}

// OpenPlatform config.
type OpenPlatform struct {
	BaseURL string  `yaml:"base_url"`
	AppKeys AppKeys `yaml:"app_keys"`
}

// AppKeys contain app id and corresponding secret key.
type AppKeys []AppKey

// GetByAppID get secret key by corresponding app id.
func (a AppKeys) GetByAppID(appID string) string {
	for i := range a {
		if a[i].AppID == appID {
			return a[i].SecretKey
		}
	}
	return ""
}

// AppKey config.
type AppKey struct {
	AppID     string `yaml:"app_id"`
	SecretKey string `yaml:"secret_key"`
}

// Consensus config.
type Consensus struct {
	Token         string `yaml:"token"`          // 建立连接使用的币种
	PaymentRemark string `yaml:"payment_remark"` // 支付备注
}
