/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/spinner"
)

// createCmd represents the create command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the status of the runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config := container.Config{
			Image: "localdev-server",
			Cmd:   []string{"/repo/scripts/runtime/status.sh"},
			Env: []string{
				"LOCALDEV_REPO=/repo",
				"LFRDEV_DOMAIN=" + viper.GetString(constants.Const.TlsLfrdevDomain),
			},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocket() + ":/var/run/docker.sock",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}

		exitCode := spinner.Spin(
			spinner.SpinOptions{
				Doing: "Status", Done: "is running", On: "'localdev' runtime environment", Enable: !flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool, string) int) int {
				return docker.InvokeCommandInLocaldev("localdev-status", config, host, true, flags.Verbose, fior, "")
			})

		os.Exit(exitCode)
	},
}

func init() {
	runtimeCmd.AddCommand(statusCmd)
}
