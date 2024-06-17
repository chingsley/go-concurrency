package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	// start 10 'producer' goroutines
	// 		 Each producer should generate and send 10 integers to the consumer using ch.
	n := 10
	for i := 0; i < n; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer func(){
				wg.Done()
				if i == 0 { // first producer is responsible for closing the ch
					wg.Wait()
					close(ch)
				}
			}()
			for j := 0; j < 10; j++ {
				ch <- j
			}
			fmt.Printf("producer %d done\n", i)
		}()
	}

	// start 1 consumer goroutine
	//		 The consumer should receive all of the integers from ch and print them out.

	go func() {
		for data := range ch {
			fmt.Println("data received:", data)
		}
		fmt.Println("consumer done")
	}()

	// Use wait-groups to ensure that every goroutine returns before the main() func stops
	wg.Wait()
	fmt.Println("main func done")
}
