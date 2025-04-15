package container

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func parseVolumeArgs(volume string) (host, container string, err error) {
	splited := strings.Split(volume, ":")
	if len(splited) != 2 {
		return "", "", errors.New("invalid arguments")
	}

	host = splited[0]
	container = splited[1]

	if _, err := os.Stat(host); os.IsNotExist(err) {
		return "", "", fmt.Errorf("host path does not exist: %s", err)
	}

	return host, container, nil
}

func VolumeMount(containerName string, volume string) error {
	var err error

	hostPath, containerPath, err := parseVolumeArgs(volume)
	if err != nil {
		return err
	}

	logrus.Infof("mounting volume, mount %v(host) to %v(container)", hostPath, containerPath)

	err = ensurePathExists(hostPath)
	if err != nil {
		return fmt.Errorf("error create hostPath:%v", err)
	}

	containerPathinHost := filepath.Join(OverlayfsRoot, containerName, "merged", containerPath)
	err = ensurePathExists(containerPathinHost)
	if err != nil {
		return fmt.Errorf("error create containerPathinHost:%v", err)
	}

	if err := mountVolume(hostPath, containerPathinHost); err != nil {
		return err
	}

	return nil
}

func UnmountVolume(containerName string, volume string) error {
	_, containerPath, err := parseVolumeArgs(volume)
	if err != nil {
		return err
	}
	containerPathinHost := filepath.Join(OverlayfsRoot, containerName, "merged", containerPath)

	if err := unmountVolume(containerPathinHost); err != nil {
		return err
	}

	return nil
}

func mountVolume(src string, dst string) error {

	if err := syscall.Mount(src, dst, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mounting failed: %w", err)
	}

	return nil
}

func unmountVolume(dst string) error {

	if err := syscall.Unmount(dst, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("mounting failed: %w", err)
	}

	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete volume dir:%v", err)
	}

	return nil
}
