package main

import (
	"os"
	"strconv"
	"sync"
)

func main() {
	var replicas, _ = strconv.Atoi(os.Args[1])
	var wg = sync.WaitGroup{}
	wg.Add(replicas)
	for i := 0; i < replicas; i++ {
		go func() {
			var proc *os.Process
			defer proc.Kill()
			for {
				proc, _ = os.StartProcess(os.Args[2], os.Args[2:], &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
				if state, _ := proc.Wait(); state.Success() {
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
