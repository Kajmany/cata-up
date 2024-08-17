package cfg

import (
	"net/url"
	"reflect"
	"testing"
)

func TestRepoURLUnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *repoURL
		wantErr bool
	}{
		{
			name:  "Valid GitHub URL",
			input: "https://github.com/cataclysmbnteam/Cataclysm-BN",
			want: &repoURL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/cataclysmbnteam/Cataclysm-BN",
				},
			},
			wantErr: false,
		},
		{
			name:  "Valid GitHub URL with trailing slash",
			input: "https://github.com/CleverRaven/Cataclysm-DDA/",
			want: &repoURL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/CleverRaven/Cataclysm-DDA/",
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			input:   "not a url",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Non-GitHub URL",
			input:   "https://gitlab.com/user/repo",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid path",
			input:   "https://github.com/user",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty input",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			// Could potentially allow this to be omitted in future
			name:    "No scheme",
			input:   "github.com/CleverRaven/Cataclysm-DDA/",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got repoURL
			err := got.UnmarshalText([]byte(tt.input))

			if (err != nil) != tt.wantErr {
				t.Errorf("repoURL.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(&got, tt.want) {
					t.Errorf("repoURL.UnmarshalText() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
