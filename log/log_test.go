package log

import (
	"testing"
)

func TestCreateLogrusInstance(t *testing.T) {
	log := CreateLogrusInstance("./api.log", 5, 2, DebugLevel, false)
	log.Debug("hello")
}
