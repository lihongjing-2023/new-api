package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureMySQLDSNTimeoutParams(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want []string
	}{
		{
			name: "appends all timeouts when absent",
			dsn:  "app_user:pwd@tcp(10.0.0.13:3306)/myapp_db?charset=utf8mb4&parseTime=true",
			want: []string{"&timeout=", "&readTimeout=", "&writeTimeout="},
		},
		{
			name: "uses ? separator when the DSN has no query params yet",
			dsn:  "app_user:pwd@tcp(10.0.0.13:3306)/myapp_db",
			want: []string{"?timeout=", "&readTimeout=", "&writeTimeout="},
		},
		{
			name: "preserves explicitly configured timeouts",
			dsn:  "app_user:pwd@tcp(10.0.0.13:3306)/myapp_db?timeout=30s&readTimeout=60s&writeTimeout=90s",
			want: []string{"?timeout=30s", "&readTimeout=60s", "&writeTimeout=90s"},
		},
		{
			name: "appends only the missing timeout params",
			dsn:  "app_user:pwd@tcp(10.0.0.13:3306)/myapp_db?readTimeout=60s",
			want: []string{"?readTimeout=60s", "&timeout=", "&writeTimeout="},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ensureMySQLDSNTimeoutParams(tt.dsn)
			for _, param := range tt.want {
				assert.Contains(t, got, param)
			}
		})
	}
}

func TestEnsureMySQLDSNTimeoutParamsEnvOverride(t *testing.T) {
	t.Setenv("SQL_READ_TIMEOUT", "45s")
	got := ensureMySQLDSNTimeoutParams("app_user:pwd@tcp(10.0.0.13:3306)/myapp_db?charset=utf8mb4")
	assert.Contains(t, got, "&readTimeout=45s")
	assert.NotContains(t, got, "&readTimeout=120s")
}
