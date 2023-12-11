/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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
var Profile string

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

		instances := ec2.GetInstanceMetaData(cfg, Name)
		resp := ssm.SendCommand(cfg, instances[0], Command)
		log.Print(aws.ToString(resp))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of EC2 instance")
	runCmd.Flags().StringVarP(&Command, "command", "c", "", "Command to send to instance")
	runCmd.Flags().StringVarP(&Profile, "profile", "p", "", "AWS profile")
	runCmd.MarkFlagRequired("name")
	runCmd.MarkFlagRequired("command")
}
