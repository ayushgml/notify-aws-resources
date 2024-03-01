package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"sync"
)

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
