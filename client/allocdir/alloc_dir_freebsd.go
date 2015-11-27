package allocdir

import (
	"errors"
	"syscall"
)

// Bind mounts the shared directory into the task directory. Must be root to
// run.
func (d *AllocDir) mountSharedDir(taskDir string) error {
	return errors.New("Unimplemented")
}

func (d *AllocDir) unmountSharedDir(dir string) error {
	return syscall.Unmount(dir, 0)
}
