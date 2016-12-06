package main

import (
	"fmt"
	"net"
	"os"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/kubernetes-incubator/cri-o/cmd/cliconfig"
	"github.com/kubernetes-incubator/cri-o/server"
	"github.com/opencontainers/runc/libcontainer/selinux"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"k8s.io/kubernetes/pkg/kubelet/api/v1alpha1/runtime"
)

type byName []cli.Flag

func (f byName) Len() int {
	return len(f)
}
func (f byName) Less(i, j int) bool {
	return f[i].GetName() < f[j].GetName()
}
func (f byName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func main() {
	app := cli.NewApp()
	app.Name = "ocid"
	app.Usage = "ocid server"
	app.Version = "0.0.1"
	app.Metadata = map[string]interface{}{
		"config": cliconfig.DefaultConfig(),
	}

	app.Flags = cliconfig.GetCommonFlags()

	// remove once https://github.com/urfave/cli/pull/544 lands
	sort.Sort(byName(app.Flags))
	sort.Sort(byName(configCommand.Flags))

	app.Commands = []cli.Command{
		configCommand,
	}

	app.Before = func(c *cli.Context) error {
		// Load the configuration file.
		config := c.App.Metadata["config"].(*server.Config)
		if err := cliconfig.MergeConfig(config, c); err != nil {
			return err
		}

		cf := &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000000000Z07:00",
			FullTimestamp:   true,
		}

		logrus.SetFormatter(cf)

		if c.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if path := c.GlobalString("log"); path != "" {
			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0666)
			if err != nil {
				return err
			}
			logrus.SetOutput(f)
		}

		switch c.GlobalString("log-format") {
		case "text":
			// retain logrus's default.
		case "json":
			logrus.SetFormatter(new(logrus.JSONFormatter))
		default:
			return fmt.Errorf("unknown log-format %q", c.GlobalString("log-format"))
		}

		return nil
	}

	app.Action = func(c *cli.Context) error {
		config := c.App.Metadata["config"].(*server.Config)

		if !config.SELinux {
			selinux.SetDisabled()
		}

		if _, err := os.Stat(config.Runtime); os.IsNotExist(err) {
			// path to runtime does not exist
			return fmt.Errorf("invalid --runtime value %q", err)
		}

		// Remove the socket if it already exists
		if _, err := os.Stat(config.Listen); err == nil {
			if err := os.Remove(config.Listen); err != nil {
				logrus.Fatal(err)
			}
		}
		lis, err := net.Listen("unix", config.Listen)
		if err != nil {
			logrus.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()

		service, err := server.New(config)
		if err != nil {
			logrus.Fatal(err)
		}

		runtime.RegisterRuntimeServiceServer(s, service)
		runtime.RegisterImageServiceServer(s, service)

		// after the daemon is done setting up we can notify systemd api
		notifySystem()

		if err := s.Serve(lis); err != nil {
			logrus.Fatal(err)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
