//go:build (all || most || spanner) && !no_spanner

package internal

// Code generated by gen.go. DO NOT EDIT.

import (
	_ "github.com/xo/usql/drivers/spanner" // Google Spanner driver
)