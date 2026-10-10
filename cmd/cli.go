package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	e "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/kekkerz/gosm/clients/ec2"
	"github.com/kekkerz/gosm/clients/ssm"
	"github.com/spf13/cobra"
)

// Tags implementation of pflags.Value for custom validation of flag
//
//	type Value interface {
//		String() string
//		Set(string) error
//		Type() string
//	}
type Tags []types.Filter

// default value
func (t *Tags) String() string {
	return "[]"
}

// Set unmarshals input into Tags type
// Returns error if wrong format
func (t *Tags) Set(v string) error {
	if err := json.Unmarshal([]byte(fmt.Sprintf("[{\"Name\": \"tag:Name\", \"Values\": %s}]", v)), &t); err != nil {
		return errors.New("input does not match expected format")
	}

	return nil
}

// Type returns type name for var
func (t *Tags) Type() string {
	return "Tags"
}

var (
	Profile    string
	Name       string
	Command    string
	tagsVal    Tags
	InstanceId string
	cfg        aws.Config
	err        error
	ec2Client  *e.Client
)

var rootCmd = &cobra.Command{
	Use: "gosm",
}

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		filters := ec2.FilterOptions{
			Name:       Name,
			Tags:       tagsVal,
			InstanceId: InstanceId,
		}
		reservations := ec2.GetInstanceMetaData(ec2Client, filters)
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
		filters := ec2.FilterOptions{
			Name:       Name,
			Tags:       tagsVal,
			InstanceId: InstanceId,
		}
		if Name != "" {
			reservations := ec2.GetInstanceMetaData(ec2Client, filters)
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "AWS profile")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(connectCmd)

	runCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of EC2 instance")
	runCmd.Flags().VarP(&tagsVal, "tags", "t", "List of tags to match against. E.g. `'[\"instance1\", \"instance2\"]'`")
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

func initConfig() {
	if Profile != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(Profile))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		log.Fatal(err)
	}

	ec2Client = e.NewFromConfig(cfg)
}
