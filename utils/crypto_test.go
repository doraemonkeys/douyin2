package utils

import (
	"reflect"
	"testing"
)

func TestHexStrToBytes(t *testing.T) {
	type args struct {
		hexStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{hexStr: "123456"}, want: []byte{0x12, 0x34, 0x56}, wantErr: false},
		{name: "2", args: args{hexStr: "12345"}, want: nil, wantErr: true},
		{name: "3", args: args{hexStr: "1234567"}, want: nil, wantErr: true},
		{name: "4", args: args{hexStr: "12345678"}, want: []byte{0x12, 0x34, 0x56, 0x78}, wantErr: false},
		{name: "5", args: args{hexStr: "FF"}, want: []byte{0xFF}, wantErr: false},
		{name: "6", args: args{hexStr: "ff"}, want: []byte{0xFF}, wantErr: false},
		{name: "7", args: args{hexStr: "Ff"}, want: []byte{0xFF}, wantErr: false},
		{name: "8", args: args{hexStr: "XXX"}, want: nil, wantErr: true},
		{name: "9", args: args{hexStr: "1"}, want: nil, wantErr: true},
		{name: "10", args: args{hexStr: "A1DFCB"}, want: []byte{0xA1, 0xDF, 0xCB}, wantErr: false},
		{name: "11", args: args{hexStr: "A1DFC"}, want: nil, wantErr: true},
		{name: "12", args: args{hexStr: "AAFF88dDAc"}, want: []byte{0xAA, 0xFF, 0x88, 0xDD, 0xAC}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexStrToBytes(tt.args.hexStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexStrToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexStrToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
