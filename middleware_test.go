package ott

import (
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"reflect"
	"testing"
)

func TestEchoFirstTraceNodeInfo(t *testing.T) {
	type args struct {
		propagator propagation.TextMapPropagator
	}
	tests := []struct {
		name string
		args args
		want func(http.Handler) http.Handler
	}{
		{"TestEchoFirstTraceNodeInfo",
			args{propagator: nil},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EchoFirstTraceNodeInfo(tt.args.propagator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EchoFirstTraceNodeInfo() got --> %v, want --> %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}
