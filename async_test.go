package async

import (
	"sync"
	"testing"
)

var results []fibNum
var mu sync.Mutex
var wg sync.WaitGroup

// The n-th fibonacci number.
type fibNum struct {
	n   int
	num int
}

// A task that will calculate the n-th fibonacci number.
type calcFibNum struct {
	n int
}

// Calculate the n-th fibonacci number.
func fib(i int) int {
	if i <= 0 {
		return 0
	}
	if i == 1 {
		return 1
	}
	return fib(i-1) + fib(i-2)
}

// Process calculates a fibonacci number and appends to results.
func (t *calcFibNum) Process() error {
	num := fib(t.n)
	result := fibNum{t.n, num}
	mu.Lock()
	results = append(results, result)
	mu.Unlock()
	wg.Done()
	return nil
}

// TestAsync queues a range of fibonacci number tasks, waits, then
// tests that the correct number of values were calculated.
func TestAsync(t *testing.T) {

	d, e := NewDispatcher(AutoSize)
	d.LogErrors(e)
	d.Start()
	defer d.Stop()

	var sent int
	for n := 37; n >= 25; n-- { // Reasonably fast range that still requires some time.
		wg.Add(1)
		t.Logf("queue task for fib num %d", n)
		TaskQueue <- &calcFibNum{n}
		sent++
	}

	wg.Wait()

	if len(results) != sent {
		t.Fatalf("result does not contain expected number of values: %d", sent)
	}
	for _, r := range results {
		t.Logf("got fib num %d = %d", r.n, r.num)
	}
}
