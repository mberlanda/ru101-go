package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultKeyNameHelper(t *testing.T) {
	k := InitKeyNameHelper()

	assert.Equal(t, "", k.Prefix)
	assert.Equal(t, ":", k.Sep)

	assert.Equal(t, "", k.CreateFieldName([]string{}))
	assert.Equal(t, "foo", k.CreateFieldName([]string{"foo"}))
	assert.Equal(t, "foo:bar:baz", k.CreateFieldName([]string{"foo", "bar", "baz"}))
	assert.Equal(t, "", k.CreateKeyName([]string{}))
	assert.Equal(t, "foo", k.CreateKeyName([]string{"foo"}))
	assert.Equal(t, "foo:bar:baz", k.CreateKeyName([]string{"foo", "bar", "baz"}))
}

func TestCustomKeyNameHelper(t *testing.T) {
	k := KeyNameHelper{
		Prefix: "pre",
		Sep:    "::",
	}

	assert.Equal(t, "", k.CreateFieldName([]string{}))
	assert.Equal(t, "foo", k.CreateFieldName([]string{"foo"}))
	assert.Equal(t, "foo::bar::baz", k.CreateFieldName([]string{"foo", "bar", "baz"}))
	// Same behaviour as lib provided by redislab
	assert.Equal(t, "pre::", k.CreateKeyName([]string{}))
	assert.Equal(t, "pre::foo", k.CreateKeyName([]string{"foo"}))
	assert.Equal(t, "pre::foo::bar::baz", k.CreateKeyName([]string{"foo", "bar", "baz"}))
}
