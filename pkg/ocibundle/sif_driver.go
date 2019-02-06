package ocibundle

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/sylabs/sif/pkg/sif"
	"github.com/sylabs/singularity/pkg/ocibundle/tools"
)

type sifDriver struct {
	image      string
	bundlePath string
	writable   bool
	Driver
}

func (s *sifDriver) Create() error {
	if s.image == "" {
		return fmt.Errorf("image wasn't set, need one to create bundle")
	}

	flag := os.O_RDONLY
	if s.writable {
		flag = os.O_RDWR
	}
	file, err := os.OpenFile(s.image, flag, 0)
	if err != nil {
		return fmt.Errorf("can't open image %s: %s", s.image, err)
	}
	defer file.Close()

	fimg, err := sif.LoadContainerFp(file, !s.writable)
	if err != nil {
		return fmt.Errorf("could not load image fp: %v", err)
	}
	part, _, err := fimg.GetPartPrimSys()
	if err != nil {
		return fmt.Errorf("could not get primaty partitions: %v", err)
	}
	fstype, err := part.GetFsType()
	if err != nil {
		return fmt.Errorf("could not get fs type: %v", err)
	}
	if fstype != sif.FsSquash {
		return fmt.Errorf("unsuported image fs type: %v", fstype)
	}
	offset := uint64(part.Fileoff)
	size := uint64(part.Filelen)

	// create OCI bundle
	if err := tools.CreateBundle(s.bundlePath, nil); err != nil {
		return fmt.Errorf("failed to create OCI bundle: %s", err)
	}

	// associate SIF image with a block
	loop, err := tools.CreateLoop(file, offset, size)
	if err != nil {
		tools.DeleteBundle(s.bundlePath)
		return fmt.Errorf("failed to find loop device: %s", err)
	}

	if err := syscall.Mount(loop, tools.RootFs(s.bundlePath).Path(), "squashfs", syscall.MS_RDONLY, "errors=remount-ro"); err != nil {
		return fmt.Errorf("failed to mount SIF partition: %s", err)
	}

	if s.writable {
		if err := tools.CreateOverlay(s.bundlePath); err != nil {
			return fmt.Errorf("failed to create overlay: %s", err)
		}
	}
	return nil
}

func (s *sifDriver) Delete() error {
	if s.writable {
		if err := tools.DeleteOverlay(s.bundlePath); err != nil {
			return err
		}
	}
	// Umount rootfs
	if err := syscall.Unmount(tools.RootFs(s.bundlePath).Path(), syscall.MNT_DETACH); err != nil {
		return err
	}
	// delete bundle directory
	return tools.DeleteBundle(s.bundlePath)
}

// FromSif ...
func FromSif(image, bundle string, writable bool) (Driver, error) {
	var err error

	s := &sifDriver{
		writable: writable,
	}
	s.bundlePath, err = filepath.Abs(bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to determine bundle path: %s", err)
	}
	if image != "" {
		s.image, err = filepath.Abs(image)
		if err != nil {
			return nil, fmt.Errorf("failed to determine image path: %s", err)
		}
	}
	return s, nil
}
