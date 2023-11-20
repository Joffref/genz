package testing

import "testing"

func TestAssertOutputIsEqual(t *testing.T) {
	type args struct {
		expected []byte
		actual   []byte
	}
	tests := map[string]struct {
		args    args
		wantErr bool
	}{
		"Same content should not return an error": {
			args:    args{expected: []byte("foo"), actual: []byte("foo")},
			wantErr: false,
		},
		"Different content should return an error": {
			args:    args{expected: []byte("foo"), actual: []byte("bar")},
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := assertOutputIsEqual("", "", tc.args.expected, tc.args.actual, false); (err != nil) != tc.wantErr {
				t.Errorf("assertOutputIsEqual() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
