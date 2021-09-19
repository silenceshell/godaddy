package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetExternalIP(t *testing.T) {
	t.Run("get external ip", func(t *testing.T) {
		addr, err := getExternalIP()

		fmt.Println(addr)
		assert.Equal(t, err, nil)
	})
}
