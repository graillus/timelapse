package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/graillus/timelapse/internal/api/client"
	"github.com/graillus/timelapse/internal/log"
	"github.com/spf13/cobra"
)

type captureOptions struct {
	interval    string
	serverURL   string
	storagePath string
}

//nolint:exhaustivestruct,gochecknoglobals
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		capture()
	},
}

var (
	opts      captureOptions //nolint:gochecknoglobals
	api       *client.Client //nolint:gochecknoglobals
	syncMutex sync.Mutex     //nolint:gochecknoglobals
)

//nolint:gochecknoinits
func init() {
	flags := captureCmd.Flags()
	flags.StringVarP(&opts.interval, "interval", "i", "1m", "Duration between two frames")
	flags.StringVarP(&opts.serverURL, "server-url", "s", "http://localhost:8990", "URL of the timelapse API server")

	rootCmd.AddCommand(captureCmd)
}

func capture() {
	opts.storagePath = "/tmp/timelapse"
	api = client.New(opts.serverURL + "/api")

	interval, err := time.ParseDuration(opts.interval)
	if err != nil {
		log.Fatalf("invalid interval %s: %v", opts.interval, err)
	}

	err = os.MkdirAll(opts.storagePath, os.ModePerm)
	if err != nil {
		log.Fatalf("unable to create directory %s: %s", opts.storagePath, err)
	}

	trigger := make(chan struct{})
	defer close(trigger)

	go func() {
		for range trigger {
			time.Now().Unix()

			err := captureFrame(fmt.Sprintf("%d", time.Now().Unix()))
			if err != nil {
				log.Errorf("Could not capture frame: %v", err)

				continue
			}

			go func() {
				err = syncFiles()
				if err != nil {
					log.Errorf("Error syncing files: %v", err)
				}
			}()
		}
	}()

	trigger <- struct{}{}

	//nolint:gosimple
	for {
		select {
		case <-time.After(interval):
			trigger <- struct{}{}
		}
	}
}

func captureFrame(id string) error {
	framePath := path.Join(
		opts.storagePath,
		fmt.Sprintf("frame_%s.jpg", id),
	)
	args := []string{
		"-t", "250",
		"-vf",
		"-q", "100",
		"-o", framePath,
	}

	cmd := exec.Command("raspistill", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to capture frame: %w", err)
	}

	if api == nil {
		return nil
	}

	return nil
}

func syncFiles() error {
	syncMutex.Lock()
	defer syncMutex.Unlock()

	err := filepath.Walk(opts.storagePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		log.Infof("Uploading frame %s...", info.Name())
		t := time.Now()
		_, err = api.PostFrame(path)
		if err != nil {
			return fmt.Errorf("failed to upload frame: %w", err)
		}
		log.Infof("Upload took %s", time.Since(t))

		err = os.Remove(path)
		if err != nil {
			log.Warnf("Failed to remove file %s: %v", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error listing files: %w", err)
	}

	return nil
}
