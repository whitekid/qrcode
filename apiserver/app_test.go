package apiserver

import (
	"testing"

	"github.com/stretchr/testify/require"

	"qrcodeapi/pkg/helper"
)

func TestParseInt(t *testing.T) {
	type args struct {
		value        string
		defaultValue int
		minValue     int
		maxValue     int
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{"default value", args{"", 10, 1, 100}, 10},
		{"cut max", args{"200", 10, 1, 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, helper.ParseIntDef(tt.args.value, tt.args.defaultValue, tt.args.minValue, tt.args.maxValue))
		})
	}
}
