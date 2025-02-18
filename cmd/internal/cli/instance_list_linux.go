// Copyright (c) 2018-2020, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cli

import (
	"os"

	"github.com/RyaxTech/singularity/docs"
	"github.com/RyaxTech/singularity/internal/app/singularity"
	"github.com/RyaxTech/singularity/pkg/cmdline"
	"github.com/RyaxTech/singularity/pkg/sylog"
	"github.com/spf13/cobra"
)

func init() {
	addCmdInit(func(cmdManager *cmdline.CommandManager) {
		cmdManager.RegisterFlagForCmd(&instanceListUserFlag, instanceListCmd)
		cmdManager.RegisterFlagForCmd(&instanceListJSONFlag, instanceListCmd)
		cmdManager.RegisterFlagForCmd(&instanceListLogsFlag, instanceListCmd)
	})
}

// -u|--user
var instanceListUser string

var instanceListUserFlag = cmdline.Flag{
	ID:           "instanceListUserFlag",
	Value:        &instanceListUser,
	DefaultValue: "",
	Name:         "user",
	ShortHand:    "u",
	Usage:        `if running as root, list instances from "<username>"`,
	Tag:          "<username>",
	EnvKeys:      []string{"USER"},
}

// -j|--json
var instanceListJSON bool

var instanceListJSONFlag = cmdline.Flag{
	ID:           "instanceListJSONFlag",
	Value:        &instanceListJSON,
	DefaultValue: false,
	Name:         "json",
	ShortHand:    "j",
	Usage:        "print structured json instead of list",
	EnvKeys:      []string{"JSON"},
}

// -l|--logs
var instanceListLogs bool

var instanceListLogsFlag = cmdline.Flag{
	ID:           "instanceListLogsFlag",
	Value:        &instanceListLogs,
	DefaultValue: false,
	Name:         "logs",
	ShortHand:    "l",
	Usage:        "display location of stdout and sterr log files for instances",
	EnvKeys:      []string{"LOGS"},
}

// singularity instance list
var instanceListCmd = &cobra.Command{
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		name := "*"
		if len(args) > 0 {
			name = args[0]
		}

		uid := os.Getuid()
		if instanceListUser != "" && uid != 0 {
			sylog.Fatalf("Only root user can list user's instances")
		}

		err := singularity.PrintInstanceList(os.Stdout, name, instanceListUser, instanceListJSON, instanceListLogs)
		if err != nil {
			sylog.Fatalf("Could not list instances: %v", err)
		}
	},
	DisableFlagsInUseLine: true,

	Use:     docs.InstanceListUse,
	Short:   docs.InstanceListShort,
	Long:    docs.InstanceListLong,
	Example: docs.InstanceListExample,
}
