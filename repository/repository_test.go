// This file contains the repository implementation layer.
package repository

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestNewRepository(t *testing.T) {
	type args struct {
		opts NewRepositoryOptions
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		{
			name: "success",
			args: args{
				opts: NewRepositoryOptions{
					Dsn: "postgres://postgres:postgres@localhost:5432/database?sslmode=disable",
				},
			},
			want: &Repository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepository(tt.args.opts); got == nil {
				t.Errorf("NewRepository() = nil")
			}
		})
	}
}
