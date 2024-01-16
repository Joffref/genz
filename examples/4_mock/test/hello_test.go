package test

import "testing"

func TestHelloMock_SayHelloTo(t *testing.T) {
	type fields struct {
		SayHelloToFunc func(param0 string) string
	}
	type args struct {
		name string
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
				name: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &HelloMock{
				SayHelloToFunc: tt.fields.SayHelloToFunc,
			}
			if got := m.SayHelloTo(tt.args.name); got != tt.want {
				t.Errorf("SayHelloTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
