package cliconfig

import "github.com/urfave/cli"

// GetCommonFlags return a common flag list shared between ocid and kpod
func GetCommonFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: ociConfigPath,
			Usage: "path to configuration file",
		},
		cli.StringFlag{
			Name:  "conmon",
			Usage: "path to the conmon executable",
		},
		cli.StringFlag{
			Name:  "containerdir",
			Usage: "ocid container dir",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output for logging",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "path to ocid socket",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "",
			Usage: "set the log file path where internal debug information is written",
		},
		cli.StringFlag{
			Name:  "log-format",
			Value: "text",
			Usage: "set the format used by logs ('text' (default), or 'json')",
		},
		cli.StringFlag{
			Name:  "pause",
			Usage: "path to the pause executable",
		},
		cli.StringFlag{
			Name:  "root",
			Usage: "ocid root dir",
		},
		cli.StringFlag{
			Name:  "runtime",
			Usage: "OCI runtime path",
		},
		cli.StringFlag{
			Name:  "sandboxdir",
			Usage: "ocid pod sandbox dir",
		},
		cli.StringFlag{
			Name:  "seccomp-profile",
			Usage: "default seccomp profile path",
		},
		cli.StringFlag{
			Name:  "apparmor-profile",
			Usage: "default apparmor profile name (default: \"ocid-default\")",
		},
		cli.BoolFlag{
			Name:  "selinux",
			Usage: "enable selinux support",
		},
	}
}
