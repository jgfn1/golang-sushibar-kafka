package main

import (
	"fmt"
	"sync"
	"log"
	"os/exec"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	lenOfByteArray := []int{256, 512, 1024, 2048, 5096}

	var wg sync.WaitGroup
	for _, lenOfByte := range lenOfByteArray {
		for numOfPublishers := 0; numOfPublishers <= 40; numOfPublishers++{
			invokePublisherCommand := fmt.Sprintf("go run publisher/publisher.go %d", lenOfByte)
			wg.Add(1)
			go callPublisher(invokePublisherCommand, &wg)
		}
		wg.Wait()
	}
}

func callPublisher(invokePublisherCommand string, wg *sync.WaitGroup){
	defer wg.Done()
	bash := exec.Command("bash", "-c", invokePublisherCommand)
	_, err := bash.Output()
	if err != nil {
		failOnError(err, "Call Publisher Error")
	}
}
