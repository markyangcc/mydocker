package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	OverlayfsRoot = "/var/lib/mydocker/overlayfs"
)

func NewContainerRootfs(containerName string) error {
	var err error

	overlayfsDirs, err := CreateOverlayDirs(containerName)
	if err != nil {
		return fmt.Errorf("failed to create readonly layer:%v", err)
	}
	mountPoint, err := CreateMountPoint(containerName)
	if err != nil {
		return fmt.Errorf("failed to create write layer:%v", err)
	}

	lowerdir := overlayfsDirs[0]
	upperdir := overlayfsDirs[1]
	workdir := overlayfsDirs[2]

	// cp base rootfs into readonly dir
	err = CopyToReadonlyDir(containerName, lowerdir)
	if err != nil {
		return fmt.Errorf("failed to prepare readonly layer:%v", err)
	}

	err = MountOverlayfs(lowerdir, upperdir, workdir, mountPoint)
	if err != nil {
		return fmt.Errorf("failed to create mount point:%v", err)
	}

	return nil
}

// mount overlayfs
func MountOverlayfs(lowerdir, upperdir, workdir string, mountPoint string) error {
	logrus.Info("mounting overlayfs")
	// mount -t overlay overlay -o lowerdir=<lowerdir1>:<lowerdir2>:<...>,upperdir=<upperdir>,workdir=<workdir> <mountpoint>
	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerdir, upperdir, workdir)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", options, mountPoint)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error mount container rootfs:%v", err)
	}

	return nil
}

// mkdir for lowerdir,upperdir,workdir
func CreateOverlayDirs(containerName string) ([]string, error) {
	// Where we store contaienr's unpacked rootfs and it's live data
	root := filepath.Join(OverlayfsRoot, containerName)
	if err := ensurePathExists(root); err != nil {
		return nil, fmt.Errorf("error create path:%v", err)
	}

	lowerdir := filepath.Join(OverlayfsRoot, containerName, "lowerdir")
	upperdir := filepath.Join(OverlayfsRoot, containerName, "upper")
	workdir := filepath.Join(OverlayfsRoot, containerName, "work")

	dirs := []string{lowerdir, upperdir, workdir}
	for _, v := range dirs {
		if err := ensurePathExists(v); err != nil {
			return nil, fmt.Errorf("error create %v layer dir:%v", v, err)
		}
	}

	return []string{lowerdir, upperdir, workdir}, nil
}

func CreateMountPoint(containerName string) (string, error) {
	mergedPath := filepath.Join(OverlayfsRoot, containerName, "merged")

	err := ensurePathExists(mergedPath)
	if err != nil {
		return "", fmt.Errorf("error create merged path:%v", err)
	}

	return mergedPath, nil
}
func CopyToReadonlyDir(containerName string, lowerdir string) error {

	err := copyDir("/tmp/rootfs/.", lowerdir)
	if err != nil {
		// we need clean these
		if err := os.RemoveAll(lowerdir); err != nil {
			return fmt.Errorf("failed to cleanup rootfs dir, some not clean files lay on the filesystem:%v", err)
		}

		return fmt.Errorf("error in copy rootfs to dst:%v", err)
	}

	// Move "busybox rootfs" to var/lib/mydocker/overlayfs/busybox/
	// when run 'ls', it should like this

	// # ls /var/lib/mydocker/overlayfs/busybox/rootfs/
	// bin  dev  etc  home  lib  lib64  root  tmp  usr  var

	return nil
}

func copyDir(src string, dst string) error {
	output, err := exec.Command("cp", "-a", src, dst).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error in when cp, stdout:%v:%v", string(output), err)
	}

	return nil
}

func ensurePathExists(path string) error {

	// Create all directories (including parents) with 0755 permissions
	// Returns nil if the path already exists or was created successfully
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

func CleanupContainerRootfs(containerName string) error {
	root := filepath.Join(OverlayfsRoot, containerName)

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return fmt.Errorf("container directory does not exist: %s", root)
	}

	// umount
	mergedPath := filepath.Join(root, "merged")
	if err := unmountPath(mergedPath); err != nil {
		return fmt.Errorf("failed to unmount merged directory: %v", err)
	}

	// remove dir
	if err := os.RemoveAll(root); err != nil {
		return fmt.Errorf("failed to remove container directory: %v", err)
	}

	return nil
}

func unmountPath(path string) error {
	if err := exec.Command("umount", path).Run(); err == nil {
		return nil
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("path still exists after unmount: %s", path)
	}

	return nil
}
