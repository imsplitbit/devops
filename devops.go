package devops

import (
	"flag"
	"fmt"
	"./ec2"
)

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

	var instance ec2.Instance
	var instances ec2.Instances
	stacks := instances.instancesForAppAndEnv(appPtr, envPtr)
	for key, value := range stacks {
		fmt.Println("Stack-ID: ", key)
		for _, instances := range value.instances {
			instance.Print(&instances)
		}
	}
}
