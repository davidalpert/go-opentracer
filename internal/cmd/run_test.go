package cmd

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"reflect"
	"testing"
)

func Test_rawTagToTypedAttribute(t *testing.T) {
	tests := []struct {
		haveContext context.Context
		rawTag      string
		want        attribute.KeyValue
		wantErr     bool
	}{
		{
			rawTag:      "a:b",
			haveContext: context.TODO(),
			want:        attribute.String("a", "b"),
		},
		{
			rawTag:      "a:true",
			haveContext: context.TODO(),
			want:        attribute.String("a", "true"),
		},
		{
			rawTag:      "a:true:bool",
			haveContext: context.TODO(),
			want:        attribute.Bool("a", true),
		},
		{
			rawTag:      "a:true:bool:something-else",
			haveContext: context.TODO(),
			want:        attribute.Bool("a", true),
		},
		{
			rawTag:      "a:4:int",
			haveContext: context.TODO(),
			want:        attribute.Int("a", 4),
		},
		{
			rawTag:      "a:4:int32",
			haveContext: context.TODO(),
			want:        attribute.Int("a", 4),
		},
		{
			rawTag:      "a:4:int64",
			haveContext: context.TODO(),
			want:        attribute.Int64("a", 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.rawTag, func(t *testing.T) {
			got, err := rawTagToTypedAttribute(tt.haveContext, tt.rawTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("rawTagToTypedAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rawTagToTypedAttribute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
