package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ensurePathExists(path string) error {

	// Create all directories (including parents) with 0755 permissions
	// Returns nil if the path already exists or was created successfully
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

// untar image tarball to rootfs dir
func CreateReadOnlyLayer(root string) error {
	tarball := filepath.Join(root, "rootfs.tar")

	err := ensurePathExists(root)
	if err != nil {
		return fmt.Errorf("error create rootfs path:%v", err)
	}

	if _, err := exec.Command("tar", "-xvf", tarball, "-C", root).CombinedOutput(); err != nil {
		return fmt.Errorf("error untar tarball:%v", err)
	}

	return nil
}

// create write layer dir
func CreateWriteLayer(root string) error {
	writePath := filepath.Join(root, "writelayer")

	err := ensurePathExists(writePath)
	if err != nil {
		return fmt.Errorf("error create write layer path:%v", err)
	}

	return nil
}

func CreateMountPoint(root string) error {
	mntPath := filepath.Join(root, "mnt")

	err := ensurePathExists(mntPath)
	if err != nil {
		return fmt.Errorf("error create mount path:%v", err)
	}

	cdir := "dirs=" + root + "writelayer=" + root + "rootfs"

	cmd := exec.Command("mount", "-t", "aufs", "-o", cdir, "none", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error mount container rootfs:%v", err)
	}

	return nil
}

func NewWorkSpace(root string) error {
	var err error

	err = CreateReadOnlyLayer(root)
	if err != nil {
		return fmt.Errorf("failed to create readonly layer:%v", err)
	}
	err = CreateWriteLayer(root)
	if err != nil {
		return fmt.Errorf("failed to create write layer:%v", err)
	}
	err = CreateMountPoint(root)
	if err != nil {
		return fmt.Errorf("failed to create mount point:%v", err)
	}

	return nil
}
