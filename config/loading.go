package config

import (
	"fmt"
	"os"

	"github.com/Netflix/go-env"
)

func LoadConfig() (config Config) {
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		fmt.Println(getExampleString())
		fmt.Println(err)
		os.Exit(1)
	}
	return
}
