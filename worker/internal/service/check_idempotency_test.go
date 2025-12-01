package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_checker_Execute1(t *testing.T) {
	type args struct {
		prev int
		curr int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "base test",
			args: args{1, 2},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := checker{}
			assert.True(t, ch.Execute(tt.args.prev, tt.args.curr))
		})
	}
}
