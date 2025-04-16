package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/mohammadne/fesghel/internal/entities"
)

func Load(config any) (err error) {
	switch entities.EnvironmentValue {
	case entities.EnvironmentLocal:
		if err := setDefaults(); err != nil {
			return fmt.Errorf("error set default config: %v", err)
		}

		defer func() {
			fmt.Println("================ Loaded Configuration ================")
			object, _ := json.MarshalIndent(config, "", "  ")
			fmt.Println(string(object))
			fmt.Println("======================================================")
		}()
	}

	prefix := strings.ToUpper(strings.ReplaceAll(entities.System, "-", "_"))
	if err = envconfig.Process(prefix, config); err != nil {
		return fmt.Errorf("error processing config via envconfig: %v", err)
	}

	return nil
}

const seperator = "_"

//go:embed defaults.env
var defaults string

func setDefaults() error {
	lines := strings.Split(defaults, "\n")
	for _, line := range lines {
		splits := strings.Split(line, "=")
		if len(splits) < 2 {
			continue
		}

		key := strings.ReplaceAll(splits[0], seperator+seperator, seperator)
		err := os.Setenv(key, splits[1])
		if err != nil {
			return fmt.Errorf("error set environment %s: %v", key, err)
		}
	}

	return nil
}
