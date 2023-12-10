package test

import "testing"

func TestHelloMock_SayHelloTo(t *testing.T) {
	type fields struct {
		SayHelloToFunc func(param0 string) string
	}
	type args struct {
		param0 string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "dummy_test",
			fields: fields{
				SayHelloToFunc: func(param0 string) string {
					return ""
				},
			},
			args: args{
				param0: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &HelloMock{
				SayHelloToFunc: tt.fields.SayHelloToFunc,
			}
			if got := m.SayHelloTo(tt.args.param0); got != tt.want {
				t.Errorf("SayHelloTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
