package async

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var sum uint64
var wg sync.WaitGroup

type testTask struct {
	value uint64
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Process atomically increments the `sum` variable.
func (tt *testTask) Process() error {
	r := rand.Intn(100) + 100
	time.Sleep(time.Duration(r) * time.Millisecond) // Simulate I/O
	atomic.AddUint64(&sum, tt.value)
	wg.Done()
	return nil
}

func TestAsync(t *testing.T) {

	d := NewDispatcher(AutoSize)
	d.Start()
	defer d.Stop()

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		TaskQueue <- &testTask{value: uint64(i)}
	}

	wg.Wait()

	var expectedSum uint64 = 55 // 1+2+3+4+5+6+7+8+9+10
	if sum != expectedSum {
		t.Errorf("sum does not equal expected value: %d != %d", sum, expectedSum)
	}
}
