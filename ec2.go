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
	id *string
	stackId *string
	ipAddress *string
}

type CloudFormation struct {
	instances []EC2Instance
	name *string
}

func instances(app *string, env *string) {
	var stacks = make(map[string]CloudFormation)
	/*
	Make a new connection to AWS' EC2 API
	 */
	svc := ec2.New(&aws.Config{Region: "us-east-1"})

	/*
	Get a listing of all instances
	 */
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	/*
	There are multiple reservations in the response so we must iterate
	through those.
	 */
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			var instance EC2Instance
			var stackId *string
			var ip *string

			match := make(map[string]bool)
			match["app"] = false
			match["env"] = false

			for _, element := range inst.Tags {
				/*
				iterate through all the tags to find what we're looking for
				It's totally possible that there are no tags so I have to be
				super defensive here.
				 */
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
					stackId = element.Value
				}

			}

			/*
			If the instances has a private ip address, use that.
			If not, use the public instead.
			 */
			if match["app"] && match["env"] {
				if inst.PrivateIPAddress != nil {
					ip = inst.PrivateIPAddress
				} else {
					if inst.PublicIPAddress != nil {
						ip = inst.PublicIPAddress
					}
				}
			}

			/*
			If there is no ip for the instance we can't get to it
			anyway so don't print it.
			 */
			if *ip != "" {
				instance.id = inst.InstanceID

				if stackId != nil {
					instance.stackId = stackId
				}
				instance.ipAddress = ip
				if value, isPresent := stacks[*stackId]; !isPresent {
					var stack CloudFormation
					stacks[*stackId] = stack
				}
				var instances []EC2Instance
				instances = stacks[*stackId].instances
				instances = append(instances, instance)
				stacks[*stackId].instances = instances
				instancePrinter(&instance)
			}
		}
	}
}

func instancePrinter(instance_data *EC2Instance) {
	fmt.Println("Instance: ", instance_data.id)
	fmt.Println("Stack-ID: ", instance_data.stackId)
	fmt.Println("\tIP: ", instance_data.ipAddress)
	fmt.Println()
}

func main() {
	appPtr := flag.String("app", "", "The app to query in tags")
	envPtr := flag.String("env", "", "The env to query in tags")
	flag.Parse()

	if appPtr == nil || *appPtr == "" {
		panic("You must pass an app name with -app")
	}

	if envPtr == nil || *envPtr == "" {
		panic("You must pass an env name with -env")
	}

	instances(appPtr, envPtr)
}
