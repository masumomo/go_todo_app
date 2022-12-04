package config

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		injectConfig *Config
		wantEnv      string
		wantPort     int
		wantErr      bool
	}{
		{
			name:     "Success Default New Config",
			wantEnv:  "dev",
			wantPort: 80,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Env, tt.wantEnv) {
				t.Errorf("Env = %v, want %v", got.Env, tt.wantEnv)
			}
			if !reflect.DeepEqual(got.Port, tt.wantPort) {
				t.Errorf("Port = %v, want %v", got.Port, tt.wantPort)
			}
		})
	}
}
