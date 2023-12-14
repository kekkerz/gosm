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

		instances := ec2.GetInstanceMetaData(cfg, Name, Tags, InstanceId)
		resp := ssm.SendCommand(cfg, instances[0], Command)
		log.Print(aws.ToString(resp))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of EC2 instance")
	runCmd.Flags().StringVarP(&Tags, "tags", "t", "", "List of tags to match against")
	runCmd.Flags().StringVarP(&InstanceId, "instance-id", "i", "", "Target Instance ID")
	runCmd.Flags().StringVarP(&Command, "command", "c", "", "Command to send to instance")
	runCmd.MarkFlagsMutuallyExclusive("name", "tags", "instance-id")
	runCmd.MarkFlagsOneRequired("name", "tags", "instance-id")
	runCmd.MarkFlagRequired("command")
}
