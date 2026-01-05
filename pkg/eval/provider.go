package eval

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// LLMProvider defines the interface for generating text from a prompt.
type LLMProvider interface {
	Generate(prompt string) (string, error)
}

// CLIProvider implements LLMProvider by executing a shell command.
type CLIProvider struct {
	Command string
	Args    []string
}

func (c *CLIProvider) Generate(prompt string) (string, error) {
	cmd := exec.Command(c.Command, c.Args...)
	cmd.Stdin = strings.NewReader(prompt)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command failed: %v | stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}
