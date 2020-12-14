package idgen

import (
	"fmt"
	"testing"
)

func TestTraceIdGenerator(t *testing.T) {
	idGenerator := NewTraceIdGenerator(4)

	for i := 0; i < 100000; i++ {
		id, err := idGenerator.GenerateId("192.168.1.2", "9001")
		fmt.Println(i, string(id), err)
	}

	for i := 0; i < 100000; i++ {
		id, err := DefaultTraceIdGenerator.GenerateId("192.168.1.2", "9001")
		fmt.Println(i, string(id), err)
	}
}
