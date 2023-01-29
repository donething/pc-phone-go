package javlib

import (
	"testing"
)

func TestRename(t *testing.T) {
	// 需要重命名的路径
	Rename([]string{`D:\Downloads\Thunder`}, nil)
}

func Test_obtainFullName(t *testing.T) {
	type args struct {
		fanhao string
	}
	tests := []struct {
		name         string
		args         args
		wantFullName string
		wantErr      bool
	}{
		{
			"SW-362",
			args{"SW-362"},
			"SW-362 同室にお見舞いにくるのは生足女子ばっかりで、モテない僕は日替わり無防備パンチラにチ○ポだけ元気になり、こっそりカーテン越しに痴漢してみた。見て見ぬフリしてた看護師さんもお見舞い女子たちも僕の勃起チ○ポ",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFullName, err := obtainFanhaoName(tt.args.fanhao)
			if err != nil || gotFullName != tt.wantFullName {
				t.Errorf("obtainFanhaoName() = %v, want %v", gotFullName, tt.wantFullName)
			}
		})
	}
}
