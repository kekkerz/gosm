package ssm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	et "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"log"
	"time"
	//"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// Command should be updated to accept type Array
func SendCommand(cfg aws.Config, instance et.Instance, command string) (stdout *string) {
	if instance.Platform == "windows" {
		log.Fatal("Windows not supported.")
	}

	client := ssm.NewFromConfig(cfg)
	parameters := make(map[string][]string)
	parameters["commands"] = []string{command}
	resp, err := client.SendCommand(context.TODO(), &ssm.SendCommandInput{
		InstanceIds:  []string{aws.ToString(instance.InstanceId)},
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters:   parameters,
	})
	if err != nil {
		log.Fatal(err)
	}

	commandId := resp.Command.CommandId
	waiter := ssm.NewCommandExecutedWaiter(client)
	params := &ssm.GetCommandInvocationInput{
		CommandId:  commandId,
		InstanceId: instance.InstanceId,
	}
	maxWaitTime := 5 * time.Minute

	waitResp, err := waiter.WaitForOutput(context.TODO(), params, maxWaitTime)
	if err != nil {
		log.Fatal(err)
	}

	return waitResp.StandardOutputContent
}
