package ants

import (
	"log"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

var wg sync.WaitGroup

const (
	numCPU = 3

	poolSize    = 480
	sizePerPool = 120
	lbs         = ants.RoundRobin

	taskCount    = 10000
	taskDuration = 250 * time.Millisecond
)

var task = func(i any) {
	defer wg.Done()
	time.Sleep(taskDuration)
}

func init() {
	runtime.GOMAXPROCS(numCPU)
	log.Printf("Running with %d CPUs\n", numCPU)
}

func TestNewPoolWithFunc(t *testing.T) {
	start := time.Now()
	defer func() {
		t.Logf("time spent duration: %v", time.Since(start))
	}()

	p, err := NewPoolWithFunc(sizePerPool, task)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		err = p.Invoke(i)
		if err != nil {
			t.Fatalf("failed to invoke: %v", err)
		}
	}

	wg.Wait()
	p.Release()
}

func TestNewMultiPoolWithFunc(t *testing.T) {
	start := time.Now()
	defer func() {
		t.Logf("test spent duration: %v", time.Since(start))
	}()

	mp, err := NewMultiPoolWithFunc(poolSize, sizePerPool, task, lbs)
	if err != nil {
		t.Fatalf("failed to create multipool: %v", err)
	}

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		err = mp.Invoke(i)
		if err != nil {
			t.Fatalf("failed to invoke: %v", err)
		}
	}

	wg.Wait()
	err = mp.ReleaseTimeout(5 * time.Second)
	if err != nil {
		t.Fatalf("failed to release multipool: %v", err)
		return
	}
}
