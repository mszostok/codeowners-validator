package envconfig

import (
	"os"

	"github.com/vrischmann/envconfig"
)

// Init the given config. Supports also envs prefix if set.
func Init(conf interface{}) error {
	envPrefix := os.Getenv("ENVS_PREFIX")
	return envconfig.InitWithPrefix(conf, envPrefix)
}
