package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const BufferSize = 10

var (
	buffer  = make([]byte, BufferSize)
	rwMutex = sync.RWMutex{}
)

func main() {
	var M, N int
	// Example: M = 8, N = 2
	M, N = 8, 2
	startRoutines(M, N)

	// M, N = 8, 8
	// startRoutines(M, N)

	// M, N = 8, 16
	// startRoutines(M, N)

	// M, N = 2, 8
	// startRoutines(M, N)

	select {} // Keep the main goroutine running
}

func startRoutines(M, N int) {
	for i := 0; i < M; i++ {
		go readBuffer(i)
	}
	for i := 0; i < N; i++ {
		go writeBuffer(i)
	}
}

func readBuffer(id int) {
	for {
		rwMutex.RLock() // Acquire the read lock
		fmt.Printf("Reader %d: Reading data: %v\n", id, buffer)
		rwMutex.RUnlock() // Release the read lock

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func writeBuffer(id int) {
	for {
		rwMutex.Lock() // Acquire the write lock
		byteToWrite := byte(rand.Intn(256))
		buffer[rand.Intn(BufferSize)] = byteToWrite
		fmt.Printf("Writer %d: Writing data: %d\n", id, byteToWrite)
		rwMutex.Unlock() // Release the write lock

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}
