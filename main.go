package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (c *AWSResourceChecker) checkAWSResources(ctx context.Context) {
	ec2Client := ec2.NewFromConfig(c.cfg)
	regionOutput, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		fmt.Printf("Unable to retrieve AWS regions: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	for _, region := range regionOutput.Regions {
		wg.Add(1)
		go func(regionName string) {
			defer wg.Done()
			c.checkRegion(ctx, regionName)
		}(*region.RegionName)
	}
	wg.Wait()
}

func (c *AWSResourceChecker) checkRegion(ctx context.Context, region string) {
	regionalCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		c.messages.Store(region, fmt.Sprintf("Error loading configuration for region %s: %v", region, err))
		return
	}

	var wg sync.WaitGroup
	if region == "us-east-1" {
		//Because S3 is a global resource, so it needs to be checked only once. (We will check it in us-ease-1)
		wg.Add(4)
		go c.checkS3Buckets(ctx, regionalCfg, &wg)
	} else {
		wg.Add(3)
	}
	go c.checkEC2Instances(ctx, regionalCfg, region, &wg)
	go c.checkRDSInstances(ctx, regionalCfg, region, &wg)
	go c.checkDynamoDBTables(ctx, regionalCfg, region, &wg)

	wg.Wait()
}

func (c *AWSResourceChecker) checkS3Buckets(ctx context.Context, cfg aws.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	s3Client := s3.NewFromConfig(cfg)
	bucketsOutput, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		c.messages.Store("global", fmt.Sprintf("Unable to retrieve S3 buckets: %v", err))
		return
	}

	bucketCount := len(bucketsOutput.Buckets)
	if bucketCount > 0 {
		c.detailedMessages.Store("global", fmt.Sprintf("S3 buckets found. Count=%d", bucketCount))
	}
}

func (c *AWSResourceChecker) getEC2Instances(instancesOutput *ec2.DescribeInstancesOutput, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instanceId := *instance.InstanceId
			instanceType := instance.InstanceType
			launchTime := *instance.LaunchTime
			launchTimeFormatted := launchTime.Format("2006-01-02 15:04:05")
			availabilityZone := *instance.Placement.AvailabilityZone
			platform := instance.Platform
			state := instance.State.Name
			c.detailedMessages.Store(instanceId, fmt.Sprintf("EC2 Instance ID: %s, Instance Type: %s, "+
				"Launch Time: %s, Availability Zone: %s, Platform: %s, State: %s, Region: %s\n",
				instanceId, instanceType, launchTimeFormatted, availabilityZone, platform, state, region))
		}
	}
}

func (c *AWSResourceChecker) checkEC2Instances(ctx context.Context, cfg aws.Config, region string, wg *sync.WaitGroup) {
	defer wg.Done()

	ec2Client := ec2.NewFromConfig(cfg)
	instancesOutput, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		c.messages.Store(region, fmt.Sprintf("Unable to retrieve EC2 instances in region %s: %v", region, err))
		return
	}

	instanceCount := 0
	for _, reservation := range instancesOutput.Reservations {
		instanceCount += len(reservation.Instances)
	}

	if instanceCount > 0 {
		c.detailedMessages.Store(region, fmt.Sprintf("EC2 instances found in region %s. Count=%d", region, instanceCount))
		wg.Add(1)
		go c.getEC2Instances(instancesOutput, region, wg)
	}
}

func (c *AWSResourceChecker) getRDSInstances(instancesOutput *rds.DescribeDBInstancesOutput, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, instance := range instancesOutput.DBInstances {
		instanceId := *instance.DBInstanceIdentifier
		instanceClass := *instance.DBInstanceClass
		creationTime := *instance.InstanceCreateTime
		creationTimeFormatted := creationTime.Format("2006-01-02 15:04:05")
		availabilityZone := *instance.AvailabilityZone
		engine := *instance.Engine
		c.detailedMessages.Store(instanceId, fmt.Sprintf("RDS Instance ID: %s, Instance Class: %s, "+
			"Creation Time: %s, Availability Zone: %s, Engine: %s, region: %s\n",
			instanceId, instanceClass, creationTimeFormatted, availabilityZone, engine, region))
	}
}

func (c *AWSResourceChecker) checkRDSInstances(ctx context.Context, cfg aws.Config, region string, wg *sync.WaitGroup) {
	defer wg.Done()

	RDSClient := rds.NewFromConfig(cfg)
	instancesOutput, err := RDSClient.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		c.messages.Store(region, fmt.Sprintf("Unable to retrieve RDS instances in region %s: %v", region, err))
		return
	}

	instanceCount := len(instancesOutput.DBInstances)
	if instanceCount > 0 {
		c.detailedMessages.Store(region, fmt.Sprintf("RDS instances found in region %s. Count=%d", region, instanceCount))
		wg.Add(1)
		go c.getRDSInstances(instancesOutput, region, wg)
	}

}

func (c *AWSResourceChecker) getDynamicDBTables(tablesOutput *dynamodb.ListTablesOutput, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, table := range tablesOutput.TableNames {
		c.detailedMessages.Store(region, fmt.Sprintf("DynamoDB Table in region %s: %s\n", region, table))
	}
}

func (c *AWSResourceChecker) checkDynamoDBTables(ctx context.Context, cfg aws.Config, region string, wg *sync.WaitGroup) {
	defer wg.Done()

	DynamoDBClient := dynamodb.NewFromConfig(cfg)
	tablesOutput, err := DynamoDBClient.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		c.messages.Store(region, fmt.Sprintf("Unable to retrieve DynamoDB tables in region %s: %v", region, err))
		return
	}

	tableCount := len(tablesOutput.TableNames)
	if tableCount > 0 {
		c.detailedMessages.Store(region, fmt.Sprintf("DynamoDB tables found in region %s. Count=%d", region, tableCount))
		wg.Add(1)
		go c.getDynamicDBTables(tablesOutput, region, wg)
	}

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
