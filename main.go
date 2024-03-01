package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSResourceChecker struct {
	cfg              aws.Config
	messages         sync.Map
	detailedMessages sync.Map
}

func main() {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("Unable to load SDK config, %v\n", err)
		return
	}

	checker := AWSResourceChecker{cfg: cfg}
	checker.checkAWSResources(ctx)
	checker.printResults()
	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time to execute: %s\n", elapsedTime)
}

func (c *AWSResourceChecker) printResults() {
	c.messages.Range(func(key, value interface{}) bool {
		fmt.Println(value)
		return true
	})

	c.detailedMessages.Range(func(key, value interface{}) bool {
		fmt.Println(value)
		return true
	})
}
