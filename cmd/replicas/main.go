package main

import (
	"os"
	"os/exec"
	"strconv"
	"sync"
)

func main() {
	var replicas, _ = strconv.Atoi(os.Args[1])
	var wg = sync.WaitGroup{}
	wg.Add(replicas)
	for i := 0; i < replicas; i++ {
		go func() {
			for {
				if err := exec.Command(os.Args[2], os.Args[2:]...).Run(); err == nil {
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
