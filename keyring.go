package keyring

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashibuto/nimble"
)

type Config struct {
	Backend   string
	NoDisplay bool
}

type Keyring struct {
	config *Config
	group  string
}

// NewKeyring returns a new Keyring object
func NewKeyring(group string, config *Config) *Keyring {
	return &Keyring{
		config: config,
		group:  group,
	}
}

// mkEnv prepares the environment for the keyring executable with user configurables
func (kr *Keyring) mkEnv() []string {
	env := []string{}
	if kr.config.Backend != "" {
		env = append(env, fmt.Sprintf("PYTHON_KEYRING_BACKEND=%s", kr.config.Backend))
	}

	if kr.config.NoDisplay {
		env = append(env, "DISPLAY=:0.0")
	}

	keys := nimble.NewSet[string]()
	for _, v := range env {
		parts := strings.SplitN(v, "=", 2)
		keys.Add(parts[0])
	}
	filtered := nimble.Filter(func(i int, v string) bool {
		parts := strings.SplitN(v, "=", 2)
		return !keys.Has(parts[0])
	}, os.Environ()...)
	filtered = append(filtered, env...)
	return filtered
}

// Set sets the key to value in the current group
func (kr *Keyring) Set(key string, value string) error {
	cmd := exec.Command("keyring", "set", kr.group, key)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s\n", value))
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Env = kr.mkEnv()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", cmd.Stderr)
	}

	return nil
}

// Get retrieves the value from the key in the current group
func (kr *Keyring) Get(key string) (string, error) {
	cmd := exec.Command("keyring", "get", kr.group, key)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Env = kr.mkEnv()

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s", cmd.Stderr)
	}

	return strings.TrimRight(outb.String(), "\n"), nil
}

// Del removes the key from the current group
func (kr *Keyring) Del(key string) error {
	var outb, errb bytes.Buffer
	cmd := exec.Command("keyring", "del", kr.group, key)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Env = kr.mkEnv()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", cmd.Stderr)
	}

	return nil
}
