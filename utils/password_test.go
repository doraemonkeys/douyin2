package utils

import "testing"

func TestBcryptHash(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "1", args: args{str: "123456"}},
		{name: "2", args: args{str: "1234567"}},
		{name: "3", args: args{str: "Q3a4s5d6"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := BcryptHash(tt.args.str)
			if !BcryptMatch(hash, tt.args.str) {
				t.Errorf("BcryptHash() = %v, Match() = false", hash)
			}
		})
	}
}
