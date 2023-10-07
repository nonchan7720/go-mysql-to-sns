package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplate(t *testing.T) {
	mp := map[string]interface{}{
		"int":   1,
		"float": 1.11,
	}
	tpl := "{{int}}-{{float}}"
	v := templateExecute(tpl, mp)
	require.Equal(t, v, "1-1.11")
}
