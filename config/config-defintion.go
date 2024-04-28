package config

import (
	"encoding/json"
	"strings"

	"github.com/Netflix/go-env"
)

type Config struct {
	PortToTest     string `env:"PORT_TO_TEST,required=true"`
	SwitchIp       string `env:"SWITCH_IP,required=true"`
	SwitchUsername string `env:"SWITCH_USERNAME,required=true"`
	SwitchPassword string `env:"SWITCH_PASSWORD,required=true" json:"-"`
	SshPort        int    `env:"SSH_PORT,required=true"`
}

func (c Config) String() string {
	raw, _ := json.MarshalIndent(c, "", "   ")
	return string(raw)
}

func getExampleString() string {
	set, _ := env.Marshal(&Config{
		PortToTest:     "ether3",
		SwitchIp:       "192.168.88.2",
		SwitchUsername: "admin",
		SwitchPassword: "admin",
		SshPort:        22,
	})
	envStrings := env.EnvSetToEnviron(set)
	return "env " + strings.Join(envStrings, " ")
}
