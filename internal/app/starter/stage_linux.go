// Copyright (c) 2018-2020, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package starter

import (
	"os"
	"syscall"

	"github.com/RyaxTech/singularity/internal/pkg/runtime/engine"
	starterConfig "github.com/RyaxTech/singularity/internal/pkg/runtime/engine/config/starter"
	"github.com/RyaxTech/singularity/pkg/sylog"
)

// StageOne validates and prepares container configuration which is
// used during container creation. Updated (possibly) engine configuration
// is wrote back into a shared sconfig so that new values will appear
// in next stages of engine execution and in master process.
//
// Any privileges gained from SUID flow or capabilities in
// extended attributes are already dropped by this moment.
func StageOne(sconfig *starterConfig.Config, e *engine.Engine) {
	sylog.Debugf("Entering stage 1\n")

	if err := e.PrepareConfig(sconfig); err != nil {
		sylog.Fatalf("%s\n", err)
	}

	if err := sconfig.Write(e.Common); err != nil {
		sylog.Fatalf("%s", err)
	}

	os.Exit(0)
}

// StageTwo performs container execution.
func StageTwo(masterSocket int, e *engine.Engine) {
	sylog.Debugf("Entering stage 2\n")

	// master socket allows communications between
	// stage 2 and master process, typically used for
	// synchronization or for sending state

	// call engine operation StartProcess, at this stage
	// we are in a container context, chroot was done.
	// The privileges are those applied by the container
	// configuration, in the case of Singularity engine
	// and if run as a user, there is no privileges set
	if err := e.StartProcess(masterSocket); err != nil {
		// write data to just tell master to not execute PostStartProcess
		// in case of failure
		if _, err := syscall.Write(masterSocket, []byte("f")); err != nil {
			sylog.Errorf("fail to send data to master: %s", err)
		}
		sylog.Fatalf("%s\n", err)
	}
}
