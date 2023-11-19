package testing

import "testing"

func TestAssertOutputIsEqual(t *testing.T) {
	type args struct {
		expected []byte
		actual   []byte
		verbose  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Same content should not return an error",
			args:    args{expected: []byte("foo"), actual: []byte("foo"), verbose: false},
			wantErr: false,
		},
		{
			name:    "Different content should return an error",
			args:    args{expected: []byte("foo"), actual: []byte("bar"), verbose: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := assertOutputIsEqual(tt.args.expected, tt.args.actual, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("assertOutputIsEqual() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
