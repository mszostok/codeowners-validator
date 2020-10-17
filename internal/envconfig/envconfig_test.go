package envconfig_test

import (
	"os"
	"testing"

	"github.com/mszostok/codeowners-validator/internal/envconfig"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

type testConfig struct {
	Key1 string
}

func TestInit(t *testing.T) {
	t.Run("Should read env variable without prefix", func(t *testing.T) {
		// given
		var cfg testConfig

		require.NoError(t, os.Setenv("KEY1", "test-value"))

		// when
		err := envconfig.Init(&cfg)

		// then
		require.NoError(t, err)
		assert.Equal(t, "test-value", cfg.Key1)
	})

	t.Run("Should read env variable with prefix", func(t *testing.T) {
		// given
		var cfg testConfig

		require.NoError(t, os.Setenv("ENVS_PREFIX", "TEST_PREFIX"))
		require.NoError(t, os.Setenv("TEST_PREFIX_KEY1", "test-value"))

		// when
		err := envconfig.Init(&cfg)

		// then
		require.NoError(t, err)
		assert.Equal(t, "test-value", cfg.Key1)
	})
}
