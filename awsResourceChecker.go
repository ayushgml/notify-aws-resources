package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"sync"
)

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
