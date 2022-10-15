package partition

import "testing"

func TestArchiveDateRange(t *testing.T) {
	type args struct {
		dateRange DateRange
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ArchiveDateRange(tt.args.dateRange); (err != nil) != tt.wantErr {
				t.Errorf("ArchiveDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttachNewDateRange(t *testing.T) {
	type args struct {
		table TrackedTable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AttachNewDateRange(tt.args.table); (err != nil) != tt.wantErr {
				t.Errorf("AttachNewDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateNewPartitionSlice(t *testing.T) {
	type args struct {
		policy Policy
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNewPartitionSlice(tt.args.policy); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewPartitionSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetachDateRange(t *testing.T) {
	type args struct {
		table TrackedTable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DetachDateRange(tt.args.table); (err != nil) != tt.wantErr {
				t.Errorf("DetachDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_relationExists(t *testing.T) {
	type args struct {
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := relationExists(tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("relationExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("relationExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}
