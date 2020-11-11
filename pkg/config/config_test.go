package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("check-logger", func(t *testing.T) {
		out := &bytes.Buffer{}
		log.Logger = zerolog.New(out)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		c, err := NewConfig("testdata/good.yaml")
		assert.NoError(t, err)
		enabledRules := make([]string, len(c.Rules))
		for i := range c.Rules {
			enabledRules[i] = fmt.Sprintf("%q", c.Rules[i].Name)
		}

		assert.Equal(t,
			fmt.Sprintf(`{"level":"debug","config":"testdata/good.yaml","message":"loaded config file"}`+"\n"+`{"level":"debug","rules":[%s],"message":"rules enabled"}`, strings.Join(enabledRules, ","))+"\n",
			out.String())
	})

	t.Run("config-good", func(t *testing.T) {
		c, err := NewConfig("testdata/good.yaml")
		assert.NoError(t, err)

		expectedRules := []*rule.Rule{}
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule1",
			Terms:        []string{"rule1"},
			Alternatives: []string{"alt-rule1"},
			Severity:     rule.SevWarn,
		})
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule2",
			Terms:        []string{"rule2", "rule-2"},
			Alternatives: []string{"alt-rule2", "alt-rule-2"},
			Severity:     rule.SevError,
		})
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "whitelist",
			Terms:        []string{"rulewl", "rule-wl"},
			Alternatives: []string{"alt-rulewl", "alt-rule-wl"},
			Severity:     rule.SevError,
		})

		expected := &Config{
			Rules:       expectedRules,
			IgnoreFiles: []string{"README.md", "pkg/rule/default.go", "testdata/good.yaml"},
		}
		expected.ConfigureRules()

		assert.EqualValues(t, expected.Rules, c.Rules)
	})

	t.Run("config-empty-missing", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("")
		assert.NoError(t, err)

		expectedEmpty := &Config{
			Rules:       rule.DefaultRules,
			IgnoreFiles: []string(nil),
		}
		assert.Equal(t, expectedEmpty, c)
	})

	t.Run("config-missing", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("testdata/missing.yaml")
		assert.Error(t, err)
		assert.Nil(t, c)
	})
}

func Test_relative(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	assert.Equal(t, ".woke.yml", relative(filepath.Join(cwd, ".woke.yml")))
	assert.Equal(t, ".woke.yml", relative(".woke.yml"))
	assert.Equal(t, "dir/.woke.yml", relative("dir/.woke.yml"))
}
