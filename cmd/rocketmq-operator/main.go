package main

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	utilflag "k8s.io/apiserver/pkg/util/flag"

	"k8s.io/apiserver/pkg/util/logs"

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

	opts.AddFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()
	goflag.CommandLine.Parse([]string{})

	logs.InitLogs()
	defer logs.FlushLogs()

	pflag.VisitAll(func(flag *pflag.Flag) {
		glog.V(2).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})

	if err := app.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
