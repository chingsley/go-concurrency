package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan int) // a channel of integers. Could be a channel of other data types including custom data type
	n := 10 // no. of workers

	var wg sync.WaitGroup
	for i := 0; i < n; i++ { // i goes up to the no. of workers (n)
		wg.Add(1)
		i := i
		go func () {
			for data := range ch { // keep reading from ch until ch is closed
				LongRunningRPC(data) // data pulled from the channel is passed as argument to the fn. NOTE: ch data and fn argument are of same data type
			}
			wg.Done()
			fmt.Printf("worker #%d stopped\n", i)
		}()
		// gaurantee: at some point in future, LongRunningRPC will run in parallel
	}

	// in the main thread
	for msg := 0; msg < 100; msg++ { // msg is the data we're processng. It must be of same data type as ch
		ch <- msg // channel sends and receives block
	}

	close(ch) // signal to receivers that no more data is coming
	wg.Wait()
	fmt.Println("ending main func")
}

func LongRunningRPC(data int) {
	fmt.Printf("sending %d to server\n", data)
	time.Sleep(150 * time.Millisecond)
}