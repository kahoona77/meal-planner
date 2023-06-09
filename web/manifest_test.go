package web

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestEmptyManifest(t *testing.T) {
	manifest := EmptyManifest()
	assert.NotNil(t, manifest)
	assert.Len(t, manifest, 0)
}

func TestManifest_File(t *testing.T) {
	tests := []struct {
		name string
		m    Manifest
		key  string
		want string
	}{
		{
			name: "with existing key",
			m:    Manifest{"test": {File: "path/file.js", Src: "src/file.ts"}},
			key:  "test",
			want: "path/file.js",
		},
		{
			name: "with non-existing key",
			m:    Manifest{"test": {File: "path/file.js", Src: "src/file.ts"}},
			key:  "foo",
			want: "",
		},
		{
			name: "with empty manifest",
			m:    EmptyManifest(),
			key:  "foo",
			want: "",
		},
		{
			name: "with empty entry",
			m:    Manifest{"test": {}},
			key:  "test",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.File(tt.key); got != tt.want {
				t.Errorf("File() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseManifest(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    Manifest
		wantErr bool
	}{
		{
			name: "existing file - no error",
			file: "data/manifest.json",
			want: Manifest{
				"src/main.css": {File: "assets/main-1f4a213e.css", Src: "src/main.css"},
				"src/main.tsx": {File: "assets/main-7ce63eb4.js", Src: "src/main.tsx"},
			},
			wantErr: false,
		},
		{
			name:    "non existing file",
			file:    "data/foobar.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "existing file - non parsable json",
			file:    "data/manifest.error",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFs := os.DirFS("../test")
			got, err := ParseManifest(tt.file, testFs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseManifest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
