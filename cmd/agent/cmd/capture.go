package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//nolint:exhaustivestruct,gochecknoglobals
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)
}

func capture() {
	var interval string
	pflag.StringVarP(&interval, "interval", "i", "1m", "Duration between two frames")

	counter := 0
	for {
		err := captureFrame(counter)
		if err != nil {
			return
		}

		counter++
	}

}

func captureFrame(counter int) error {
	args := []string{
		"-t", "250",
		"-vf",
		"-q", "100",
		"-o", fmt.Sprintf("/tmp/timelapse/frame%d.jpg", counter),
	}

	cmd := exec.Command("raspistill", args...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to capture frame #%d: %w", counter, err)
	}

	return nil
}
