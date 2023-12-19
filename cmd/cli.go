package cmd

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/kekkerz/gosm/clients/ec2"
	"github.com/kekkerz/gosm/clients/ssm"
	"log"
	"os"
)

var (
	Profile    string
	Name       string
	Command    string
	Tags       string
	InstanceId string
	cfg        aws.Config
	err        error
)

var rootCmd = &cobra.Command{
	Use: "gosm",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if Profile != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(Profile))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		log.Fatal(err)
	}

	var runCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
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

	var connectCmd = &cobra.Command{
		Use: "connect",
		Run: func(cmd *cobra.Command, args []string) {
			var target string
			if Name != "" {
				reservations := ec2.GetInstanceMetaData(cfg, Name, "", InstanceId)
				if len(reservations) > 1 || len(reservations[0].Instances) > 1 {
					log.Fatal("More than one instance found.")
				}
				target = aws.ToString(reservations[0].Instances[0].InstanceId)
			} else {
				target = InstanceId
			}
			ssm.Connect(cfg, target)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "AWS profile")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(connectCmd)

	runCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of EC2 instance")
	// Potentially convert `tags` to use a custom flag that only accepts a JSON blob
	runCmd.Flags().StringVarP(&Tags, "tags", "t", "", "List of tags to match against. E.g. `'[{\"Name\": \"tag:Name\", \"Values\": [\"instance1\", \"instance2\"]}]'`")
	runCmd.Flags().StringVarP(&InstanceId, "instance-id", "i", "", "Target Instance ID")
	runCmd.Flags().StringVarP(&Command, "command", "c", "", "Command to send to instance")
	runCmd.MarkFlagsMutuallyExclusive("name", "tags", "instance-id")
	runCmd.MarkFlagsOneRequired("name", "tags", "instance-id")
	runCmd.MarkFlagRequired("command")

	connectCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of EC2 instance")
	connectCmd.Flags().StringVarP(&InstanceId, "instance-id", "i", "", "Target Instance ID")
	connectCmd.MarkFlagsMutuallyExclusive("name", "instance-id")
	connectCmd.MarkFlagsOneRequired("name", "instance-id")
}
