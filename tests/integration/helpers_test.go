//go:build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"
)

func normalizeTimeDurations(in string) string {
	duration := regexp.MustCompile(`\(\d+(\.\d+)?(ns|us|µs|ms|s|m|h)\)`)
	return duration.ReplaceAllString(in, "(<duration>)")
}

func normalizeLogTime(in string) string {
	duration := regexp.MustCompile(`time="[^"]*"`)
	return duration.ReplaceAllString(in, `time="<time>"`)
}

func CloneRepo(t *testing.T, url string, branch string) (string, func()) {
	t.Helper()

	repoDir, err := ioutil.TempDir("", strings.ReplaceAll(t.Name(), "/", "-"))
	require.NoError(t, err)

	cleanup := func() {
		err = os.RemoveAll(repoDir)
		require.NoError(t, err)
	}

	_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
	})
	require.NoError(t, err)

	return repoDir, cleanup
}

type Executor struct {
	envs       map[string]string
	timeout    time.Duration
	binaryPath string
}

func Exec() *Executor {
	return &Executor{
		envs: map[string]string{},
	}
}

// WithEnv adds given env. Overrides if previously existed
func (s *Executor) WithEnv(key string, value string) *Executor {
	s.envs[key] = value
	return s
}

type ExecuteOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

func (s *Executor) AwaitResultAtMost(timeout time.Duration) *ExecuteOutput {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, s.binaryPath)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	for k, v := range s.envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	err := cmd.Run()
	return &ExecuteOutput{
		ExitCode: cmd.ProcessState.ExitCode(),
		Err:      err,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}
}

func (s *Executor) WithTimeout(timeout time.Duration) *Executor {
	s.timeout = timeout
	return s
}

func (s *Executor) Binary(binaryPath string) *Executor {
	s.binaryPath = binaryPath
	return s
}
