package ec2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

func GetInstanceMetaData(cfg aws.Config, name string) (instance []types.Instance) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{name},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return output.Reservations[0].Instances
}
