package partition

import (
	"github.com/jackc/pgx/v5"
	"testing"
)

func TestPostgresDB_ArchiveDateRange(t *testing.T) {
	type fields struct {
		ConnConfig *pgx.ConnConfig
	}
	type args struct {
		dateRange DateRange
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
			pgd := PostgresDB{
				ConnConfig: tt.fields.ConnConfig,
			}
			if err := pgd.ArchiveDateRange(tt.args.dateRange); (err != nil) != tt.wantErr {
				t.Errorf("ArchiveDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresDB_AttachNewDateRange(t *testing.T) {
	type fields struct {
		ConnConfig *pgx.ConnConfig
	}
	type args struct {
		table TrackedTable
	}

	// Open one DB connection; pass it to the tests?
	// does it need to be closed; or will it memory leak?
	connString := ""
	myconn, err := pgx.ParseConfig(connString)
	if err != nil {
		t.Fatalf("CONN ERR: %s", err.Error())
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy #1", fields{ConnConfig: myconn}, args{table: TrackedTable{
			Schema: "",
			Name:   "",
			Ranges: nil,
		}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgd := PostgresDB{
				ConnConfig: tt.fields.ConnConfig,
			}
			// These are planned already?? pass the latest plan
			if err := pgd.AttachNewDateRange(tt.args.table); (err != nil) != tt.wantErr {
				t.Errorf("AttachNewDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresDB_CreateNewPartitionSlice(t *testing.T) {
	type fields struct {
		ConnConfig *pgx.ConnConfig
	}
	type args struct {
		policy Policy
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
			pgd := PostgresDB{
				ConnConfig: tt.fields.ConnConfig,
			}
			if err := pgd.CreateNewPartitionSlice(tt.args.policy); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewPartitionSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresDB_DetachDateRange(t *testing.T) {
	type fields struct {
		ConnConfig *pgx.ConnConfig
	}
	type args struct {
		table TrackedTable
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
			pgd := PostgresDB{
				ConnConfig: tt.fields.ConnConfig,
			}
			if err := pgd.DetachDateRange(tt.args.table); (err != nil) != tt.wantErr {
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
