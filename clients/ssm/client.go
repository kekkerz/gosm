package ssm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"log"
	"time"
)

func SendCommand(cfg aws.Config, targets []string, command string) {
	client := ssm.NewFromConfig(cfg)
	parameters := make(map[string][]string)
	parameters["commands"] = []string{command}
	resp, err := client.SendCommand(context.TODO(), &ssm.SendCommandInput{
		InstanceIds:  targets,
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters:   parameters,
	})
	if err != nil {
		log.Fatal(err)
	}

	commandId := resp.Command.CommandId
	waiter := ssm.NewCommandExecutedWaiter(client)
	maxWaitTime := 5 * time.Minute
	ch := make(chan string)
	for _, instance := range targets {
		go func(instance string) {
			params := &ssm.GetCommandInvocationInput{
				CommandId:  commandId,
				InstanceId: aws.String(instance),
			}
			waitResp, err := waiter.WaitForOutput(context.TODO(), params, maxWaitTime)
			if err != nil {
				log.Fatal(err)
			}
			ch <- aws.ToString(waitResp.StandardOutputContent)
		}(instance)
		log.Print(<-ch)
	}
}
