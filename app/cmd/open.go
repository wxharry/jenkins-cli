package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

type OpenOption struct {
	Name   string
	Config bool
}

var openOption OpenOption

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.PersistentFlags().StringVarP(&openOption.Name, "name", "n", "", "Open a specific Jenkins by name")
	openCmd.PersistentFlags().BoolVarP(&openOption.Config, "config", "c", false, "Open the configuration page of Jenkins")
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open your Jenkins with a browse",
	Long:  `Open your Jenkins with a browse`,
	Run: func(cmd *cobra.Command, args []string) {
		var jenkins *JenkinsServer

		if openOption.Name == "" {
			jenkins = getCurrentJenkins()
		} else {
			jenkins = findJenkinsByName(openOption.Name)
		}

		if jenkins != nil && jenkins.URL != "" {
			url := jenkins.URL
			if openOption.Config {
				url = fmt.Sprintf("%s/configure", url)
			}
			open(url)
		} else {
			log.Fatalf("No URL found with Jenkins %s", openOption.Name)
		}
	},
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
