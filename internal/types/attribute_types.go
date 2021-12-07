package types

import (
	"encoding/json"
	"fmt"
)

// OpenTelemetryAttributeType represents the supported openTelemetry attribute value types
type OpenTelemetryAttributeType string

const (
	StringAttribute OpenTelemetryAttributeType = "string"
	BoolAttribute                              = "bool"
	IntAttribute                               = "int"
	Int32Attribute                             = "int32"
	Int64Attribute                             = "int64"
)

func (at *OpenTelemetryAttributeType) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	attrType := OpenTelemetryAttributeType(s)
	switch attrType {
	case StringAttribute, BoolAttribute, IntAttribute, Int32Attribute, Int64Attribute:
		*at = attrType
		return nil
	}
	return fmt.Errorf("invalid OpenTelemetryAttributeType: %s", attrType)
}
