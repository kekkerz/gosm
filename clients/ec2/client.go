package ec2

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

type FilterOptions struct {
	Name       string
	Tags       string
	InstanceId string
}

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetInstanceMetaData(client EC2DescribeInstancesAPI, filters FilterOptions) (instance []types.Reservation) {
	input := &ec2.DescribeInstancesInput{}
	if filters.Name != "" {
		input.Filters = []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{filters.Name},
			},
		}
	} else if filters.Tags != "" {
		filter := &[]types.Filter{}
		err := json.Unmarshal([]byte(filters.Tags), &filter)
		if err != nil {
			log.Fatal("Error unmarshaling JSON: ", err)
		}
		input.Filters = *filter
	} else if filters.InstanceId != "" {
		input.InstanceIds = []string{filters.InstanceId}
	}

	output, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}

	return output.Reservations
}
