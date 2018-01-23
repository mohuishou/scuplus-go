package scujwc

import (
	"log"
	"testing"
)

func TestJwc_getEvaList(t *testing.T) {
	tests := []struct {
		name    string
		j       *Jwc
		wantErr bool
	}{
		{
			name:    "test",
			j:       j,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetEvaList()
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Jwc.getEvaList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestJwc_Evaluate(t *testing.T) {
	//补充相关信息
	eva := &Evaluate{}
	log.Println(eva)
	eva.Comment = "很好的老师"
	eva.Star = 5
	type args struct {
		evaluate *Evaluate
	}
	tests := []struct {
		name    string
		j       *Jwc
		args    args
		wantErr bool
	}{
		{
			name: "test",
			j:    j,
			args: args{
				evaluate: &eva,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Evaluate(tt.args.evaluate); (err != nil) != tt.wantErr {
				t.Errorf("Jwc.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
