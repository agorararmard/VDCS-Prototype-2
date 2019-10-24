package main

import (
	"reflect"
	"testing"
)

func Test_addEval(t *testing.T) {
	type args struct {
		code   []string
		idx    int
		params []string
		typesA []string
		chName string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addEval(tt.args.code, tt.args.idx, tt.args.params, tt.args.typesA, tt.args.chName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addEval() = %v, want %v", got, tt.want)
			}
		})
	}
}
