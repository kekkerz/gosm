package cmd

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"gosm/clients/ec2"
	"gosm/clients/ssm"
	"log"
)

var ConnectName string
var ConnectInstanceId string

var connectCmd = &cobra.Command{
	Use: "connect",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg aws.Config
		var err error
		var target string
		if Profile != "" {
			cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(Profile))
		} else {
			cfg, err = config.LoadDefaultConfig(context.TODO())
		}

		if err != nil {
			log.Fatal(err)
			err = nil
		}

		if ConnectName != "" {
			reservations := ec2.GetInstanceMetaData(cfg, ConnectName, "", ConnectInstanceId)
			if len(reservations) > 1 || len(reservations[0].Instances) > 1 {
				log.Fatal("More than one instance found.")
			}
			target = aws.ToString(reservations[0].Instances[0].InstanceId)
		} else {
			target = ConnectInstanceId
		}
		ssm.Connect(cfg, target)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringVarP(&ConnectName, "name", "n", "", "Name of EC2 instance")
	connectCmd.Flags().StringVarP(&ConnectInstanceId, "instance-id", "i", "", "Target Instance ID")
	connectCmd.MarkFlagsMutuallyExclusive("name", "instance-id")
	connectCmd.MarkFlagsOneRequired("name", "instance-id")
}
