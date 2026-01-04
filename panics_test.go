package panics

import (
	"fmt"
	"testing"
)

func TestPanics(t *testing.T) {
	err := fmt.Errorf("err")
	OnError(err, "test message")
}
