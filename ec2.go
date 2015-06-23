package main

import (
	"fmt"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const STACK_ID_KEY string = "aws:cloudformation:stack-name"
const APP_KEY string = "Application"
const ENV_KEY string = "Environment"

type EC2Instance struct {
	id string
	stackId string
	ipAddress string
}

type CloudFormation struct {
	EC2Instances []EC2Instance
	name string
}

func instances(app *string, env *string) {
	svc := ec2.New(&aws.Config{Region: "us-east-1"})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			match := make(map[string]bool)
			match["app"] = false
			match["env"] = false
			var stack_id *string
			for _, element := range inst.Tags {
				/*
				iterate through all the tags to find what we're looking for
				It's totally possible that there are no tags so I have to be
				super defensive here.
				 */
//				fmt.Println(*element.Key, ": ", *element.Value)
				if *element.Key == APP_KEY && *element.Value == *app {
					match["app"] = true
				} else if *element.Key == ENV_KEY && *element.Value == *env {
					match["env"] = true
				}

				/*
				Since we're already iterating over tags here might as well
				grab some that we care about for later
				 */
				if *element.Key == STACK_ID_KEY {
					stack_id = element.Value
				}

			}
			var ip string
			if match["app"] && match["env"] {
				if inst.PrivateIPAddress != nil {
					ip = *inst.PrivateIPAddress
				} else {
					if inst.PublicIPAddress != nil {
						ip = *inst.PublicIPAddress
					}
				}
			}

			if ip != "" {
				fmt.Println("Instance: ", *inst.InstanceID)
				if stack_id != nil {
					fmt.Println("Stack-ID: ", *stack_id)
				}
				fmt.Println("\tIP: ", ip)
			}
		}
	}
}

func printer(instance_data *sting)

func main() {
	appPtr := flag.String("app", "", "The app to query in tags")
	envPtr := flag.String("env", "", "The env to query in tags")
	flag.Parse()

	if appPtr == nil {
		panic("You must pass an app name with -app")
	}

	if envPtr == nil {
		panic("You must pass an env name with -env")
	}

	instances(appPtr, envPtr)
}
