package utils

import(
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"sort"
	"strings"
)

var supportedObjectPrinterFormatMap = map[string]string{
	"text": "text",
	"t":    "text",
	"json": "json",
	"js":   "json",
	"j":    "json",
	"yaml": "yaml",
	"yml":  "yaml",
	"y":    "yaml",
}

var supportedObjectPrinterKeys = []string{}
var supportedObjectPrinterCategories = []string{}

// PrinterOptions is an abstraction over various types of output formatting
type PrinterOptions struct {
	OutputFormat string
}

// AddPrinterFlags adds flags to a cobra.Command
func (o *PrinterOptions) AddPrinterFlags(c *cobra.Command) {
	if o.OutputFormat == "" {
		o.OutputFormat = "text"
	}
	o.addObjectPrinterFlags(c)
}

// AddObjectPrinterFlags adds flags to a cobra.Command
func (o *PrinterOptions) addObjectPrinterFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.OutputFormat, "output", "o", o.OutputFormat, fmt.Sprintf("output format. One of: %s.", strings.Join(supportedObjectPrinterCategories, "|")))
}

// Validate asserts that the printer options are valid
func (o *PrinterOptions) Validate() error {
	if !StringInSlice(o.SupportedFormats(), o.OutputFormat) {
		return fmt.Errorf("invalid output format: %s\nvalid format values are: %v", o.OutputFormat, strings.Join(o.SupportedFormatCategories(), "|"))
	}
	return nil
}

// SupportedFormats returns the list of supported formats
func (o *PrinterOptions) SupportedFormats() []string {
	return supportedObjectPrinterKeys
}

// SupportedFormatCategories returns the list of supported formats
func (o *PrinterOptions) SupportedFormatCategories() []string {
	return supportedObjectPrinterCategories
}

// OutputFormatIsText returns the format category
func (o *PrinterOptions) OutputFormatIsText() bool {
	return strings.EqualFold(o.FormatCategory(), "text")
}

// FormatCategory returns the format category
func (o *PrinterOptions) FormatCategory() string {
	ExitIfErr(o.Validate())

	return supportedObjectPrinterFormatMap[o.OutputFormat]
}

func (o *PrinterOptions) FormatOutput(v interface{}) (string, string, error) {
	formatCategory := o.FormatCategory()
	output := ""

	if formatCategory != "" {
		// category-specific marshalling
		return marshalObjectToString(v, formatCategory)
	}

	return output, formatCategory, fmt.Errorf("could not map output=%s to a format category", o.OutputFormat)
}

// marshalObjectToString is the magic of the printer abstraction
func marshalObjectToString(v interface{}, formatCategory string) (string, string, error) {
	output := ""
	if formatCategory == "text" {
		switch v.(type) {
		case fmt.Stringer:
			return v.(fmt.Stringer).String(), formatCategory, nil
		case string:
			return v.(string), formatCategory, nil
		default:
			return fmt.Sprintf("%v\n", v), formatCategory, nil
		}
	} else if formatCategory == "yaml" {
		oB, _ := yaml.Marshal(v)
		output = string(oB)
	} else if formatCategory == "json" {
		oB, _ := json.MarshalIndent(v, "", "  ")
		output = string(oB) + "\n"
	} else {
		return output, formatCategory, fmt.Errorf("do not support format category %s", formatCategory)
	}

	return output, formatCategory, nil
}

func init() {
	// extract and sort supportedObjectPrinterKeys and Categories
	for k, c := range supportedObjectPrinterFormatMap {
		supportedObjectPrinterKeys = append(supportedObjectPrinterKeys, k)

		if !StringInSlice(supportedObjectPrinterCategories, c) {
			supportedObjectPrinterCategories = append(supportedObjectPrinterCategories, c)
		}
	}

	sort.Slice(supportedObjectPrinterKeys, func(i, j int) bool {
		return supportedObjectPrinterKeys[i] < supportedObjectPrinterKeys[j]
	})

	sort.Slice(supportedObjectPrinterCategories, func(i, j int) bool {
		return supportedObjectPrinterCategories[i] < supportedObjectPrinterCategories[j]
	})
}
