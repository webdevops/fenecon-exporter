package config

import (
	"encoding/json"
	"time"
)

type (
	Opts struct {
		Logger struct {
			Debug       bool `long:"log.debug"    env:"LOG_DEBUG"  description:"debug mode"`
			Development bool `long:"log.devel"    env:"LOG_DEVEL"  description:"development mode"`
			Json        bool `long:"log.json"     env:"LOG_JSON"   description:"Switch log output to json format"`
		}

		Fenecon struct {
			Request struct {
				Timeout time.Duration `long:"fenecon.request.timeout"  env:"FENECON_REQUEST_TIMEOUT"  description:"Request timeout" default:"5s"`
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
