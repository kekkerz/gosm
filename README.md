# gosm
This is a rewrite of [ssm-connect](https://github.com/kekkerz/ssm-connect), which was written in Python. I'm currently in the process of learning Go, so I decided to re-write my original script here to try to get some experience with the language.

## Current functionality
- Send remote commands to instances targeted by their name.
```
[abevil@epsilon ~/git/gosm]$ go run main.go run -c hostname -n bitwarden
2023/12/11 15:46:17 ip-172-0-143-106.us-east-2.compute.internal
```

## Usage
```
[abevil@epsilon ~/git/gosm]$ go run main.go run --help
Usage:
  gosm run [flags]

Flags:
  -c, --command string   Command to send to instance
  -h, --help             help for run
  -n, --name string      Name of EC2 instance
  -p, --profile string   AWS profile
```

## Next steps
- Add ability to use StartSession to get a shell on an instance
    - Would prefer to implement natively using the SDK's `StartSession` method, rather than calling the aws cli.
- Add unit tests
- Add GitHub actions to provided release artifacts