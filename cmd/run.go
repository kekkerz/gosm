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

var Name string
var Command string
var Tags string
var InstanceId string

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg aws.Config
		var err error
		if Profile != "" {
			cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(Profile))
		} else {
			cfg, err = config.LoadDefaultConfig(context.TODO())
		}

		if err != nil {
			log.Fatal(err)
			err = nil
		}

		reservations := ec2.GetInstanceMetaData(cfg, Name, Tags, InstanceId)
		var targets []string
		for _, reservation := range reservations {
			for _, instance := range reservation.Instances {
				targets = append(targets, aws.ToString(instance.InstanceId))
			}
		}

		ssm.SendCommand(cfg, targets, Command)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of EC2 instance")
	// Potentially convert `tags` to use a custom flag that only accepts a JSON blob
	runCmd.Flags().StringVarP(&Tags, "tags", "t", "", "List of tags to match against. E.g. `'[{\"Name\": \"tag:Name\", \"Values\": [\"instance1\", \"instance2\"]}]'`")
	runCmd.Flags().StringVarP(&InstanceId, "instance-id", "i", "", "Target Instance ID")
	runCmd.Flags().StringVarP(&Command, "command", "c", "", "Command to send to instance")
	runCmd.MarkFlagsMutuallyExclusive("name", "tags", "instance-id")
	runCmd.MarkFlagsOneRequired("name", "tags", "instance-id")
	runCmd.MarkFlagRequired("command")
}
