package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/mrunalp/fileutils"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

// fatal prints the error's details if it is a libcontainer specific error type
// then exits the program with an exit status of 1.
func fatal(err error) {
	// make sure the error is written to the logger
	logrus.Error(err)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// saveState saves the rootfs and config.json used to create/run a container
func saveState(spec *specs.Spec, configPath string) error {
	rootfsPath := spec.Root.Path
	tmpDir, err := ioutil.TempDir("/tmp", "runc-save")
	if err != nil {
		return err
	}
	if err = fileutils.CopyFile(configPath, tmpDir); err != nil {
		return err
	}
	if err = fileutils.CopyDirectory(rootfsPath, tmpDir); err != nil {
		return err
	}
	return nil
}

// setupSpec performs initial setup based on the cli.Context for the container
func setupSpec(context *cli.Context) (*specs.Spec, error) {
	bundle := context.String("bundle")
	if bundle != "" {
		if err := os.Chdir(bundle); err != nil {
			return nil, err
		}
	}
	spec, err := loadSpec(specConfig)
	if err != nil {
		return nil, err
	}

	if err = saveState(spec, filepath.Join(bundle, specConfig)); err != nil {
		return nil, err
	}
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket != "" {
		setupSdNotify(spec, notifySocket)
	}
	if os.Geteuid() != 0 {
		return nil, fmt.Errorf("runc should be run as root")
	}
	return spec, nil
}
