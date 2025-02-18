// Copyright (c) 2018-2020, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cli

import (
	"fmt"

	"github.com/RyaxTech/singularity/docs"
	"github.com/RyaxTech/singularity/internal/app/singularity"
	"github.com/RyaxTech/singularity/pkg/sylog"
	"github.com/spf13/cobra"
)

// PluginUninstallCmd takes the name of a plugin and uninstalls it from the
// plugin directory.
//
// singularity plugin uninstall <name>
var PluginUninstallCmd = &cobra.Command{
	PreRun: CheckRootOrUnpriv,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := singularity.UninstallPlugin(name)
		if err != nil {
			sylog.Fatalf("Failed to uninstall plugin %q: %s.", name, err)
		}
		fmt.Printf("Uninstalled plugin %q.\n", name)
	},
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),

	Use:     docs.PluginUninstallUse,
	Short:   docs.PluginUninstallShort,
	Long:    docs.PluginUninstallLong,
	Example: docs.PluginUninstallExample,
}
