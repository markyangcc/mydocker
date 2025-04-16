package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CommitCOntainer(containerName string, imagePath string) error {

	rootfs := filepath.Join(OverlayfsRoot, containerName, "merged")

	if _, err := os.Stat(rootfs); os.IsNotExist(err) {
		return fmt.Errorf("container runtime workspace does not exist: %s", rootfs)
	}

	// trime path
	cmd := exec.Command("tar", "-cf", imagePath, "-C", rootfs, ".")

	if err := cmd.Run(); err != nil {
		return err
	}

	// check
	if _, err := os.Stat(rootfs); os.IsNotExist(err) {
		return fmt.Errorf("check failed, commited tarball not found:%v", err)
	}

	return nil
}
