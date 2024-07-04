package vqueue

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	sliceQueue *SliceQueue[string]
)

func TestSliceQueueEnqueue(t *testing.T) {
	for {
		val, err := sliceQueue.Dequeue()
		if err != nil {
			if errors.Is(err, QueueEmptyError) {
				time.Sleep(2)
				continue
			}
			assert.Nil(t, err)
		}

		fmt.Println(val)
	}
}
