package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"sync"
)

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

func (c *AWSResourceChecker) getDynamicDBTables(tablesOutput *dynamodb.ListTablesOutput, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, table := range tablesOutput.TableNames {
		c.detailedMessages.Store(region, fmt.Sprintf("DynamoDB Table in region %s: %s\n", region, table))
	}
}
