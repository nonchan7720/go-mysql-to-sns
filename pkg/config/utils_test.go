package config

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFastTemplate(t *testing.T) {
	mp := map[string]interface{}{
		"int":   1,
		"float": 1.11,
		"key":   "a6ea368b-3ceb-4729-820b-22e85cace255",
		"nested": map[string]interface{}{
			"a": "a",
			"b": map[string]interface{}{
				"c": "d",
			},
		},
	}
	tpl := `{{int}}-{{float}}-{{replaceAll key "-" ""}}-{{nested.a}}-{{replaceAll nested.b.c "d" "e"}}`
	v := fastTemplate(tpl, mp)
	require.Equal(t, v, `1-1.11-a6ea368b3ceb4729820b22e85cace255-a-e`)
}

func TestGoTemplate(t *testing.T) {
	mp := map[string]interface{}{
		"int":   1,
		"float": 1.11,
		"key":   "a6ea368b-3ceb-4729-820b-22e85cace255",
		"nested": map[string]interface{}{
			"a": "a",
			"b": map[string]interface{}{
				"c": "d",
			},
		},
	}
	tpl := `{{.int}}-{{.float}}-{{replaceAll .key "-" ""}}-{{.nested.a}}-{{replaceAll .nested.b.c "d" "e"}}`
	v := goTemplate(tpl, mp)
	require.Equal(t, v, "1-1.11-a6ea368b3ceb4729820b22e85cace255-a-e")
}

func TestUUIDFunction(t *testing.T) {
	tpl := `{{newUUID}}`
	for _, fn := range []func(string, map[string]interface{}) string{goTemplate, fastTemplate} {
		actual := fn(tpl, nil)
		_, err := uuid.Parse(actual)
		require.NoError(t, err)
	}
}
