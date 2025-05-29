package cmd

import (
	"fmt"
	"github.com/davidalpert/opentracer/internal/utils"
	"github.com/davidalpert/opentracer/internal/version"
	"github.com/spf13/cobra"
)

// VersionOptions is a struct to support version command
type VersionOptions struct {
	utils.PrinterOptions
	VersionDetail version.DetailStruct
}

// NewVersionOptions returns initialized VersionOptions
func NewVersionOptions() *VersionOptions {
	return &VersionOptions{
		VersionDetail: version.Detail,
	}
}

// NewCmdVersion creates the version command
func NewCmdVersion() *cobra.Command {
	o := NewVersionOptions()
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err

			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	o.AddPrinterFlags(cmd)

	return cmd
}

// Complete completes the VersionOptions
func (o *VersionOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Validate validates the VersionOptions
func (o *VersionOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run executes the command
func (o *VersionOptions) Run() error {
	if s, _, err := o.PrinterOptions.FormatOutput(o.VersionDetail); err != nil {
		return err
	} else {
		fmt.Println(s)
	}

	return nil
}
