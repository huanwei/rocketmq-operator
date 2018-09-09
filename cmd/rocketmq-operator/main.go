package main

import (
	"fmt"
	"os"

	"github.com/huanwei/rocketmq-operator/cmd/rocketmq-operator/app"
	operatoropts "github.com/huanwei/rocketmq-operator/pkg/options"
	"github.com/huanwei/rocketmq-operator/pkg/version"
)

const (
	configPath = "/etc/rocketmq-operator/broker-config.yaml"
)

func main() {
	fmt.Fprintf(os.Stderr, "Starting rocketmq-operator version '%s'\n", version.GetBuildVersion())

	opts, err := operatoropts.NewOperatorOpts(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config: %v\n", err)
		os.Exit(1)
	}
	if err := app.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
