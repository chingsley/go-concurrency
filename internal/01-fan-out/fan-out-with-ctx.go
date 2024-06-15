package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer done()
	ch := make(chan int) // a channel of integers. Could be a channel of other data types including custom data type
	n := 10 // no. of workers

	var wg sync.WaitGroup
	for i := 0; i < n; i++ { // i goes up to the no. of workers (n)
		wg.Add(1)
		i := i
		go func () {
			defer wg.Done()
			for {
				select {

					case data, ok := <-ch:
						if !ok {
							fmt.Printf("worker #%d stopped; channel closed\n", i)
							return
						}
						LongRunningRPC2(ctx, data) // data pulled from the channel is passed as argument to the fn. NOTE: ch data and fn argument are of same data type
					case <-ctx.Done():
						fmt.Printf("worker #%d stopped; deadline expired\n", i)
						return
				}
			}
		}()
		// gaurantee: at some point in future, LongRunningRPC will run in parallel
	}

	// in the main thread
	loop:
	for msg := 0; msg < 100; msg++ { // msg is the data we're processng. It must be of same data type as ch
		select {
		case ch <- msg: // channel sends and receives block
		case <- ctx.Done():
			break loop
		}
	}

	close(ch) // signal to receivers that no more data is coming
	wg.Wait()
	fmt.Println("ending main func")
}

func LongRunningRPC2(ctx context.Context, data int) {
	fmt.Printf("sending %d to server\n", data)
	time.Sleep(150 * time.Millisecond) // TODO: pass context to networking library
}
