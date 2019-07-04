package utils

import (
	"fmt"
	"strings"
)

// DefaultPrefix ...
const DefaultPrefix = ""

// DefaultSep ...
const DefaultSep = ":"

// KeyNameHelper ...
type KeyNameHelper struct {
	Prefix string
	Sep    string
}

// InitKeyNameHelper ...
func InitKeyNameHelper() *KeyNameHelper {
	return &KeyNameHelper{
		Prefix: DefaultPrefix,
		Sep:    DefaultSep,
	}
}

// GetPrefix ...
func (k *KeyNameHelper) GetPrefix() string {
	return k.Prefix
}

// SetPrefix ...
func (k *KeyNameHelper) SetPrefix(p string) {
	k.Prefix = p
}

// GetSep ...
func (k *KeyNameHelper) GetSep() string {
	return k.Sep
}

// SetSep ...
func (k *KeyNameHelper) SetSep(s string) {
	k.Sep = s
}

// CreateKeyName ...
func (k *KeyNameHelper) CreateKeyName(vals []string) string {
	key := k.CreateFieldName(vals)
	if k.Prefix != "" {
		return fmt.Sprintf("%s%s%s", k.Prefix, k.Sep, key)
	}
	return key
}

// CreateFieldName ...
func (k *KeyNameHelper) CreateFieldName(vals []string) string {
	return strings.Join(vals[:], k.Sep)
}
