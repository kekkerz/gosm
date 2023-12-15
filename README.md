# gosm
This is a rewrite of [ssm-connect](https://github.com/kekkerz/ssm-connect), which was written in Python. I'm currently in the process of learning Go, so I decided to re-write my original script here to try to get some experience with the language.

## Current functionality
- Send remote commands to instances targeted by their name, instance-id, or tags.
```
[abevil@epsilon ~/git/gosm]$ go run main.go run -n instance-name -c hostname
2023/12/11 15:46:17 i-xxxxxxxx:
ip-172-0-143-106.us-east-2.compute.internal

[abevil@epsilon ~/git/gosm]$ go run main.go run -i i-xxxxxxxx -c hostname
2023/12/14 17:57:52 i-xxxxxxxx: 
ip-172-0-143-106.us-east-2.compute.internal

[abevil@epsilon ~/git/gosm]$ go run main.go run -t '[{"Name": "tag:Name", "Values": ["instance-name", "instance2-name"]}]' -c "hostname"
2023/12/14 17:56:49 i-xxxxxxxx: 
ip-172-0-143-106.us-east-2.compute.internal

2023/12/14 17:56:50 i-xxxxxxxx: 
ip-172-0-132-40.us-east-2.compute.internal
```

## Tag targeting
Tags are passed in using the `-t`/`--tags` flag as a JSON blob. The format should be as follows:

```
[
  {
    "Name": "tag:<name-of-tag>",
    "Values": [
      "<tag_value>",
      "<tag_value>"
    ]
  }
]
```

Additional tags can be supplied by simply adding a new {} block at the top level.

## Usage
```
[abevil@epsilon ~/git/gosm]$ ./gosm run help
Error: required flag(s) "command" not set
Usage:
  gosm run [flags]

Flags:
  -c, --command string                                                        Command to send to instance
  -h, --help                                                                  help for run
  -i, --instance-id string                                                    Target Instance ID
  -n, --name string                                                           Name of EC2 instance
  -t, --tags '[{"Name": "tag:Name", "Values": ["instance1", "instance2"]}]'   List of tags to match against. E.g. '[{"Name": "tag:Name", "Values": ["instance1", "instance2"]}]'

Global Flags:
  -p, --profile string   AWS profile
```

## Next steps
- Update `--tags` flag to only accept JSON compliant strings
- Add ability to use StartSession to get a shell on an instance
    - Would prefer to implement natively using the SDK's `StartSession` method, rather than calling the aws cli.
- Add unit tests
- Add GitHub actions to provided release artifacts