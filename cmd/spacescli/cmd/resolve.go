// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ava-labs/spacesvm/client"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve [options] space/key",
	Short: "Reads a value at space/key",
	RunE:  resolveFunc,
}

// TODO: move all this to a separate client code
func resolveFunc(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "expected exactly 1 argument, got %d", len(args))
		os.Exit(128)
	}
	cli := client.New(uri, requestTimeout)
	_, v, err := cli.Resolve(args[0])
	if err != nil {
		return err
	}

	color.Yellow("%s=>%q", args[0], v)
	return nil
}