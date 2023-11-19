package test

import "testing"

func TestHuman_Validate(t *testing.T) {
	type fields struct {
		Firstname string
		Lastname  string
		Age       uint
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Firstname should start with a capital letter",
			fields:  fields{Firstname: "foo", Lastname: "bar", Age: 42},
			wantErr: true,
		},
		{
			name:    "Lastname must be set",
			fields:  fields{Firstname: "Foo", Lastname: "", Age: 42},
			wantErr: true,
		},
		{
			name:    "Lastname should start with a capital letter",
			fields:  fields{Firstname: "Foo", Lastname: "bar", Age: 42},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Human{
				Firstname: tt.fields.Firstname,
				Lastname:  tt.fields.Lastname,
				Age:       tt.fields.Age,
			}
			if err := v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
