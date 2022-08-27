package breakglass

import "testing"

func TestBastion_AddToRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	type args struct {
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			if err := b.AddToRole(tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("AddToRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBastion_RemoveFromRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	type args struct {
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			if err := b.RemoveFromRole(tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBastion_GetRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			got, err := b.GetRole()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}
