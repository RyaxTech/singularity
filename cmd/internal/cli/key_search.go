// Copyright (c) 2020, Control Command Inc. All rights reserved.
// Copyright (c) 2017-2020, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cli

import (
	"context"
	"os"

	"github.com/RyaxTech/singularity/docs"
	"github.com/RyaxTech/singularity/internal/pkg/remote/endpoint"
	"github.com/RyaxTech/singularity/pkg/sylog"
	"github.com/RyaxTech/singularity/pkg/sypgp"
	"github.com/spf13/cobra"
	"github.com/sylabs/scs-key-client/client"
)

// KeySearchCmd is 'singularity key search' and look for public keys from a key server
var KeySearchCmd = &cobra.Command{
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		co, err := getKeyserverClientOpts(keyServerURI, endpoint.KeyserverSearchOp)
		if err != nil {
			sylog.Fatalf("Keyserver client failed: %s", err)
		}

		if err := doKeySearchCmd(cmd.Context(), args[0], co...); err != nil {
			sylog.Errorf("search failed: %s", err)
			os.Exit(2)
		}
	},

	Use:     docs.KeySearchUse,
	Short:   docs.KeySearchShort,
	Long:    docs.KeySearchLong,
	Example: docs.KeySearchExample,
}

func doKeySearchCmd(ctx context.Context, search string, co ...client.Option) error {
	// get keyring with matching search string
	return sypgp.SearchPubkey(ctx, search, keySearchLongList, co...)
}
