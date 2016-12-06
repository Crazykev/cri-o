package cliconfig

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/kubernetes-incubator/cri-o/server"
	"github.com/opencontainers/runc/libcontainer/selinux"
	"github.com/urfave/cli"
)

const (
	ociConfigPath       = "/etc/ocid/ocid.conf"
	ocidRoot            = "/var/lib/ocid"
	conmonPath          = "/usr/libexec/ocid/conmon"
	pausePath           = "/usr/libexec/ocid/pause"
	seccompProfilePath  = "/etc/ocid/seccomp.json"
	apparmorProfileName = "ocid-default"
)

// DefaultConfig returns the default configuration for ocid.
func DefaultConfig() *server.Config {
	return &server.Config{
		RootConfig: server.RootConfig{
			Root:         ocidRoot,
			SandboxDir:   filepath.Join(ocidRoot, "sandboxes"),
			ContainerDir: filepath.Join(ocidRoot, "containers"),
			LogDir:       "/var/log/ocid/pods",
		},
		APIConfig: server.APIConfig{
			Listen: "/var/run/ocid.sock",
		},
		RuntimeConfig: server.RuntimeConfig{
			Runtime: "/usr/bin/runc",
			Conmon:  conmonPath,
			ConmonEnv: []string{
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			},
			SELinux:         selinux.SelinuxEnabled(),
			SeccompProfile:  seccompProfilePath,
			ApparmorProfile: apparmorProfileName,
		},
		ImageConfig: server.ImageConfig{
			Pause:    pausePath,
			ImageDir: filepath.Join(ocidRoot, "store"),
		},
	}
}

// MergeConfig merge given server config with config specified in cli context.
func MergeConfig(config *server.Config, ctx *cli.Context) error {
	// Don't parse the config if the user explicitly set it to "".
	if path := ctx.GlobalString("config"); path != "" {
		if err := config.FromFile(path); err != nil {
			if ctx.GlobalIsSet("config") || !os.IsNotExist(err) {
				return err
			}

			// We don't error out if --config wasn't explicitly set and the
			// default doesn't exist. But we will log a warning about it, so
			// the user doesn't miss it.
			logrus.Warnf("default configuration file does not exist: %s", ociConfigPath)
		}
	}

	// Override options set with the CLI.
	if ctx.GlobalIsSet("conmon") {
		config.Conmon = ctx.GlobalString("conmon")
	}
	if ctx.GlobalIsSet("containerdir") {
		config.ContainerDir = ctx.GlobalString("containerdir")
	}
	if ctx.GlobalIsSet("pause") {
		config.Pause = ctx.GlobalString("pause")
	}
	if ctx.GlobalIsSet("root") {
		config.Root = ctx.GlobalString("root")
	}
	if ctx.GlobalIsSet("sandboxdir") {
		config.SandboxDir = ctx.GlobalString("sandboxdir")
	}
	if ctx.GlobalIsSet("listen") {
		config.Listen = ctx.GlobalString("listen")
	}
	if ctx.GlobalIsSet("runtime") {
		config.Runtime = ctx.GlobalString("runtime")
	}
	if ctx.GlobalIsSet("selinux") {
		config.SELinux = ctx.GlobalBool("selinux")
	}
	if ctx.GlobalIsSet("seccomp-profile") {
		config.SeccompProfile = ctx.GlobalString("seccomp-profile")
	}
	if ctx.GlobalIsSet("apparmor-profile") {
		config.ApparmorProfile = ctx.GlobalString("apparmor-profile")
	}
	return nil
}
