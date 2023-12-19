package ssm

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/session-manager-plugin/src/datachannel"
	ssmLog "github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
	"github.com/google/uuid"
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
	}

	for _, instance := range targets {
		log.Printf("\033[32m%s: \n\033[0m%s\n", instance, <-ch)
	}
}

func Connect(cfg aws.Config, target string) {
	client := ssm.NewFromConfig(cfg)
	resp, err := client.StartSession(context.TODO(), &ssm.StartSessionInput{
		Target: aws.String(target),
	})

	if err != nil {
		log.Fatal(err)
	}

	var ssmSession = session.Session{
		SessionId: *resp.SessionId,
		StreamUrl: *resp.StreamUrl,
		TokenValue: *resp.TokenValue,
		ClientId: uuid.New().String(),
		DataChannel: &datachannel.DataChannel{},
		TargetId: target,
		Endpoint: fmt.Sprintf("ssm.%s.amazonaws.com", cfg.Region),
	}

	err = ssmSession.Execute(ssmLog.Logger(true, target))

	if err != nil {
		log.Fatal(err)
	}
}
