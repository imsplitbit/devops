package ec2

import (
	"fmt"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const STACK_ID_KEY string = "aws:cloudformation:stack-name"
const APP_KEY string = "Application"
const ENV_KEY string = "Environment"

//type EC2Instance interface {
//	New() Instance
//}
//
//type EC2Instances interface {
//	instanceForAppAndEnv() map[string]Instances
//}

type Instance struct {
	id string
	stackId string
	ipAddress string
}

type Instances struct {
	instances []Instance
	name string
}

func (i Instances) instancesForAppAndEnv(app *string, env *string) map[string]Instances {
	var stacks = make(map[string]Instances)
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
			var instance Instance
			var stackId *string
			var ip string

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
					ip = *inst.PrivateIPAddress
				} else {
					if inst.PublicIPAddress != nil {
						ip = *inst.PublicIPAddress
					}
				}
			}

			/*
			If there is no ip for the instance we can't get to it
			anyway so don't print it.
			 */
			if ip != "" {
				instance.id = *inst.InstanceID

				if stackId != nil {
					instance.stackId = *stackId
				}
				instance.ipAddress = ip
				if _, isPresent := stacks[*stackId]; !isPresent {
					var stack Instances
					stacks[*stackId] = stack
				}

				/*
				The contents of map members are immutable so we copy them out,
				change them, then set the old member to the new value.
				 */
				var stack Instances
				stack = stacks[*stackId]
				stack.instances = append(stack.instances, instance)
				stacks[*stackId] = stack
			}
		}
	}
	return stacks
}

func (i Instance) Print(instance_data *Instance) {
	fmt.Println("\tInstance: ", instance_data.id)
	fmt.Println("\t\tIP: ", instance_data.ipAddress)
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

	var instances Instances
	var instance Instance
	stacks := instances.instancesForAppAndEnv(appPtr, envPtr)
	for key, value := range stacks {
		fmt.Println("Stack-ID: ", key)
		for _, instances := range value.instances {
			instance.Print(&instances)
		}
	}
}
