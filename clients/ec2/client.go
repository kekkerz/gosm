package ec2

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

/*type Tags []struct {
	Name string `json:"Name"`
	Values []string `json:"Values"`
}*/

func GetInstanceMetaData(cfg aws.Config, name string, tags string, instanceId string) (instance []types.Instance) {
	input := &ec2.DescribeInstancesInput{}
	if name != "" {
		input.Filters = []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{name},
			},
		}
	} else if tags != "" {
		//jsonTags := &Tags{}
		filter := &[]types.Filter{}
		err := json.Unmarshal([]byte(tags), &filter)
		if err != nil {
			log.Println(tags)
			log.Fatal("Error unmarshaling JSON: ", err)
		}
		input.Filters = *filter
	} else if instanceId != "" {
		input.InstanceIds = []string{instanceId}
	}

	client := ec2.NewFromConfig(cfg)
	/*output, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{name},
			},
		},
	})*/
	output, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}

	return output.Reservations[0].Instances
}
