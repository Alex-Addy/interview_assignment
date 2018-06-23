package interview_assignment

import (
	"testing"
)

func TestHash(t *testing.T) {
	input := "angryMonkey"
	expected := `ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==`
	output := HashAndEncode(input)

	if output != expected {
		t.Errorf("Expected: %s\nReceived: %s", expected, output)
	}
}
