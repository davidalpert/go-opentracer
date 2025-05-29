package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/opentracer/internal/version"
	"github.com/spf13/cobra"
)

// VersionOptions is a struct to support version command
type VersionOptions struct {
	*printers.PrinterOptions
	VersionDetail *version.DetailStruct
}

// NewVersionOptions returns initialized VersionOptions
func NewVersionOptions(s printers.IOStreams) *VersionOptions {
	return &VersionOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(s).WithDefaultOutput("text"),
		VersionDetail:  &version.Detail,
	}
}

// NewCmdVersion creates the version command
func NewCmdVersion(s printers.IOStreams) *cobra.Command {
	o := NewVersionOptions(s)
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

	o.AddPrinterFlags(cmd.Flags())

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
	return o.WriteOutput(o.VersionDetail)
}
