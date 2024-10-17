package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T){
	tests := map[string]struct{
		envVars map[string]string
		expected Config
		expectedErr string
	}{
		"success": {
			envVars: map[string]string{
				"DATABASE_HOST": "localhost",
				"DATABASE_PORT": "5432",
				"DATABASE_USER": "user",
				"DATABASE_PASSWORD": "password",
				"DATABASE_NAME": "name",
			},
			expected: Config{
				DB_Host: "localhost",
				DB_Port: 5432,
				DB_User: "user",
				DB_Password: "password",
				DB_Name: "name",
				DB_RetryDuration: "3s",
				HTTP_Domain: "localhost",
				HTTP_Port: "8000",
			},
		},
		"missing env var": {
			envVars: map[string]string{
				"DATABASE_HOST": "localhost",
				"DATABASE_PORT": "5432",
				"DATABASE_USER": "user",
				"DATABASE_PASSWORD": "password",
			},
			expectedErr: "required key DATABASE_NAME missing value",
		},
		"invalid port": {
			envVars: map[string]string{
				"DATABASE_HOST": "localhost",
				"DATABASE_PORT": "port",
				"DATABASE_USER": "user",
				"DATABASE_PASSWORD": "password",
				"DATABASE_NAME": "name",
			},
			expectedErr: "invalid value for DATABASE_PORT: strconv.Atoi: parsing \"port\": invalid syntax",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T){
			for key, value := range tc.envVars {
				err := os.Setenv(key, value)
				assert.NoError(t, err)
			}

			defer func() {
				for key := range tc.envVars {
					err := os.Unsetenv(key)
					assert.NoError(t, err)
				}
			}()

			cfg, err := New()
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expected, cfg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, cfg)
			}
		})
	}
	
}