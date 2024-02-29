package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var messages []string
var instanceArray []string

var instanceWaitGroup sync.WaitGroup
var rdsWaitGroup sync.WaitGroup

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
				checkRDSInstances(regionalCfg, regionName)
			}()

			resourceWG.Wait()
		}(regionName)
	}

	wg.Wait()

	for _, message := range messages {
		fmt.Println(message)
	}

	for _, instance := range instanceArray {
		fmt.Println(instance)
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
		instanceWaitGroup.Add(1)
		getEC2Instances(instancesOutput, region)
	}

	instanceWaitGroup.Wait()
}

func getEC2Instances(instancesOutput *ec2.DescribeInstancesOutput, region string) {
	defer instanceWaitGroup.Done()

	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instanceId := *instance.InstanceId
			instanceType := instance.InstanceType
			launchTime := *instance.LaunchTime
			launchTimeFormatted := launchTime.Format("2006-01-02 15:04:05")
			availabilityZone := *instance.Placement.AvailabilityZone
			platform := instance.Platform
			state := instance.State.Name
			instanceArray = append(instanceArray, fmt.Sprintf("EC2 Instance ID: %s, Instance Type: %s, "+
				"Launch Time: %s, Availability Zone: %s, Platform: %s, State: %s, Region: %s\n",
				instanceId, instanceType, launchTimeFormatted, availabilityZone, platform, state, region))
		}
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

func checkRDSInstances(cfg aws.Config, region string) {
	RDSClient := rds.NewFromConfig(cfg)
	instancesOutput, err := RDSClient.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		messages = append(messages, fmt.Sprintf("Unable to retrieve RDS instances in region %s: %v", region, err))
		return
	}

	instanceCount := len(instancesOutput.DBInstances)
	if instanceCount > 0 {
		messages = append(messages, fmt.Sprintf("RDS instances found in region %s. Count=%d", region, instanceCount))
		rdsWaitGroup.Add(1)
		getRDSInstances(instancesOutput, region)
	}

	rdsWaitGroup.Wait()
}

func getRDSInstances(instancesOutput *rds.DescribeDBInstancesOutput, region string) {
	defer rdsWaitGroup.Done()

	for _, instance := range instancesOutput.DBInstances {
		instanceId := *instance.DBInstanceIdentifier
		instanceClass := *instance.DBInstanceClass
		creationTime := *instance.InstanceCreateTime
		creationTimeFormatted := creationTime.Format("2006-01-02 15:04:05")
		availabilityZone := *instance.AvailabilityZone
		engine := *instance.Engine
		instanceArray = append(instanceArray, fmt.Sprintf("RDS Instance ID: %s, Instance Class: %s, "+
			"Creation Time: %s, Availability Zone: %s, Engine: %s, region: %s\n",
			instanceId, instanceClass, creationTimeFormatted, availabilityZone, engine, region))
	}
}
