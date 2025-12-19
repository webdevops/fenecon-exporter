package config

import (
	"encoding/json"
	"time"
)

type (
	Opts struct {
		// logger
		Logger struct {
			Level  string `long:"log.level"    env:"LOG_LEVEL"   description:"Log level" choice:"trace" choice:"debug" choice:"info" choice:"warning" choice:"error" default:"info"`                          // nolint:staticcheck // multiple choices are ok
			Format string `long:"log.format"   env:"LOG_FORMAT"  description:"Log format" choice:"logfmt" choice:"json" default:"logfmt"`                                                                     // nolint:staticcheck // multiple choices are ok
			Source string `long:"log.source"   env:"LOG_SOURCE"  description:"Show source for every log message (useful for debugging and bug reports)" choice:"" choice:"short" choice:"file" choice:"full"` // nolint:staticcheck // multiple choices are ok
			Color  string `long:"log.color"    env:"LOG_COLOR"   description:"Enable color for logs" choice:"" choice:"auto" choice:"yes" choice:"no"`                                                        // nolint:staticcheck // multiple choices are ok
			Time   bool   `long:"log.time"     env:"LOG_TIME"    description:"Show log time"`
		}

		Fenecon struct {
			Request struct {
				Timeout          time.Duration `long:"fenecon.request.timeout"       env:"FENECON_REQUEST_TIMEOUT"       description:"Request timeout"              default:"10s"`
				Parallel         int           `long:"fenecon.request.parallel"      env:"FENECON_REQUEST_PARALLEL"      description:"Number of parallel requests"  default:"1"`
				RetryCount       int           `long:"fenecon.request.retries"       env:"FENECON_REQUEST_RETRIES"       description:"Request retries"              default:"2"`
				RetryWaitTime    time.Duration `long:"fenecon.request.waittime"      env:"FENECON_REQUEST_WAITTIME"      description:"Request retries"              default:"2s"`
				RetryMaxWaitTime time.Duration `long:"fenecon.request.maxwaittime"   env:"FENECON_REQUEST_MAXWAITTIME"   description:"Request retries"              default:"5s"`
			}

			Auth struct {
				Username string `long:"fenecon.auth.username"  env:"FENECON_AUTH_USERNAME"  description:"Username for fenecon login"`
				Password string `long:"fenecon.auth.password"  env:"FENECON_AUTH_PASSWORD"  description:"Password for fenecon login" default:"user" json:"-"`
			}
		}

		// general options
		Server struct {
			// general options
			Bind         string        `long:"server.bind"              env:"SERVER_BIND"           description:"Server address"        default:":8080"`
			ReadTimeout  time.Duration `long:"server.timeout.read"      env:"SERVER_TIMEOUT_READ"   description:"Server read timeout"   default:"5s"`
			WriteTimeout time.Duration `long:"server.timeout.write"     env:"SERVER_TIMEOUT_WRITE"  description:"Server write timeout"  default:"60s"`
		}
	}
)

func (o *Opts) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}
