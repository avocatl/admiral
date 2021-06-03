package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	Text string `json:"message"`
}

func TestDisplay_JsonDisplayable_Raw(t *testing.T) {
	tst := test{
		Text: "Hello from admiral!",
	}

	dps := Json(tst, false)

	assert.False(t, dps.Filterable())
	assert.True(t, dps.NoHeaders())
	assert.Empty(t, dps.ColMap())
	assert.Equal(t, dps.Cols(), []string{""})
	assert.Equal(t, []map[string]interface{}{
		{
			"": "{\"message\":\"Hello from admiral!\"}",
		},
	}, dps.KV())
}

func TestDisplay_JsonDisplayable_Pretty(t *testing.T) {
	tst := test{
		Text: "Hello from admiral!",
	}

	dps := Json(tst, true)

	assert.False(t, dps.Filterable())
	assert.True(t, dps.NoHeaders())
	assert.Empty(t, dps.ColMap())
	assert.Equal(t, dps.Cols(), []string{""})
	assert.Equal(t, []map[string]interface{}{
		{
			"": "{\n    \"message\": \"Hello from admiral!\"\n}",
		},
	}, dps.KV())
}

func TestDisplay_JsonDisplayable_JsonError(t *testing.T) {
	x := map[string]interface{}{
		"foo": make(chan int),
	}

	dps := Json(x, true)

	assert.False(t, dps.Filterable())
	assert.True(t, dps.NoHeaders())
	assert.Empty(t, dps.ColMap())
	assert.Equal(t, dps.Cols(), []string{""})
	assert.Equal(t, []map[string]interface{}{
		{
			"": "json: unsupported type: chan int",
		},
	}, dps.KV())
}

func TestDisplay_TextDisplayable_Raw(t *testing.T) {
	dps := Text("", "Hello!")

	assert.False(t, dps.Filterable())
	assert.True(t, dps.NoHeaders())
	assert.Empty(t, dps.ColMap())
	assert.Equal(t, dps.Cols(), []string{""})
	assert.Equal(t, []map[string]interface{}{
		{
			"": "Hello!",
		},
	}, dps.KV())
}

func TestDisplay_TextDisplayable_Divider(t *testing.T) {
	dps := Text("*", "Hello!")

	assert.False(t, dps.Filterable())
	assert.True(t, dps.NoHeaders())
	assert.Empty(t, dps.ColMap())
	assert.Equal(t, dps.Cols(), []string{""})
	assert.Equal(t, []map[string]interface{}{
		{
			"": "\n**************************************************\nHello!\n**************************************************\n",
		},
	}, dps.KV())
}
