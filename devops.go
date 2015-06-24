package devops

import (
	"flag"
	"fmt"
	"github.com/imsplitbit/devops/ec2"
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

	stacks := instancesForAppAndEnv(appPtr, envPtr)
	for key, value := range stacks {
		fmt.Println("Stack-ID: ", key)
		for _, instances := range value.instances {
			instancePrinter(&instances)
		}
	}
}
