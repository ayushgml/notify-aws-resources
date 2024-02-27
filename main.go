package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var messages []string

func main() {
	startTime := time.Now()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		messages = append(messages, fmt.Sprintf("Unable to load SDK config, %v", err))
	}

	ec2Client := ec2.NewFromConfig(cfg)
	regionOutput, err := ec2Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		messages = append(messages, fmt.Sprintf("Unable to retrieve AWS regions: %v", err))
	}

	var wg sync.WaitGroup

	for _, region := range regionOutput.Regions {
		regionName := *region.RegionName

		wg.Add(1)
		go func(regionName string) {
			defer wg.Done()

			regionalCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(regionName))
			if err != nil {
				messages = append(messages, fmt.Sprintf("Error loading configuration for region %s: %v", regionName, err))
				return
			}

			if regionName == "ap-south-1" {
				checkS3Buckets(regionalCfg)
			}

			var resourceWG sync.WaitGroup
			resourceWG.Add(1)

			go func() {
				defer resourceWG.Done()
				checkEC2Instances(regionalCfg, regionName)
			}()

			resourceWG.Wait()
		}(regionName)
	}

	wg.Wait()

	for _, message := range messages {
		fmt.Println(message)
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time to execute: %s\n", elapsedTime)
}

func checkEC2Instances(cfg aws.Config, region string) {
	ec2Client := ec2.NewFromConfig(cfg)
	instancesOutput, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		messages = append(messages, fmt.Sprintf("Unable to retrieve EC2 instances in region %s: %v", region, err))
		return
	}

	instanceCount := 0
	for _, reservation := range instancesOutput.Reservations {
		instanceCount += len(reservation.Instances)
	}

	if instanceCount > 0 {
		messages = append(messages, fmt.Sprintf("EC2 instances found in region %s. Count=%d", region, instanceCount))
	}
}

func checkS3Buckets(cfg aws.Config) {
	s3Client := s3.NewFromConfig(cfg)
	bucketsOutput, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		messages = append(messages, fmt.Sprintf("Unable to retrieve S3 buckets: %v", err))
		return
	}

	bucketCount := len(bucketsOutput.Buckets)
	if bucketCount > 0 {
		messages = append(messages, fmt.Sprintf("S3 buckets found. Count=%d", bucketCount))
	}
}