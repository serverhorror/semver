package semver_test

import (
	"testing"

	. "github.com/serverhorror/semver"
)

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name string
		v    Version
		want string
	}{
		{name: "Default",
			v:    Default(),
			want: "0.0.0"},
		{name: "PreRelease",
			v:    Version{Major: "0", Minor: "0", Patchlevel: "0", PreRelease: "x"},
			want: "0.0.0-x"},
		{name: "PreRelease+and+Metadata",
			v:    Version{Major: "0", Minor: "0", Patchlevel: "0", PreRelease: "x", Metadata: "y"},
			want: "0.0.0-x+y"},
		{name: "Metadata",
			v:    Version{Major: "0", Minor: "0", Patchlevel: "0", Metadata: "y"},
			want: "0.0.0+y"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       Version
		wantErr bool
	}{
		{name: "Default",
			v:       Default(),
			wantErr: false},
		{name: "PreRelease",
			v:       Version{Major: "0", Minor: "0", Patchlevel: "0", PreRelease: "x"},
			wantErr: false},
		{name: "PreRelease+and+Metadata",
			v:       Version{Major: "0", Minor: "0", Patchlevel: "0", PreRelease: "x", Metadata: "y"},
			wantErr: false},
		{name: "Metadata",
			v:       Version{Major: "0", Minor: "0", Patchlevel: "0", Metadata: "y"},
			wantErr: false},
		{name: "Invalid Major",
			v:       Version{Major: "a", Minor: "0", Patchlevel: "0", Metadata: "y"},
			wantErr: true},
		{name: "Invalid Minor",
			v:       Version{Major: "0", Minor: "a", Patchlevel: "0", Metadata: "y"},
			wantErr: true},
		{name: "Invalid PatchLevel",
			v:       Version{Major: "0", Minor: "0", Patchlevel: "000", Metadata: "y"},
			wantErr: true},
		{name: "Invalid Prerelease",
			v:       Version{Major: "0", Minor: "0", Patchlevel: "0", PreRelease: ";,", Metadata: "y"},
			wantErr: true},
		{name: "Invalid Metadata",
			v:       Version{Major: "0", Minor: "0", Patchlevel: "0", Metadata: ",;"},
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Version.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVersions_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		v    Versions
		args args
		want bool
	}{
		{name: "i small j",
			v: Versions{
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
				Version{Major: "0", Minor: "0", Patchlevel: "0"}},
			args: args{i: 0, j: 1},
			want: true,
		},
		{name: "same version",
			v: Versions{
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
		{name: "higher minor version",
			v: Versions{
				Version{Major: "1", Minor: "1", Patchlevel: "0"},
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
			},
			args: args{i: 0, j: 1},
			want: true,
		},
		{name: "lower minor version",
			v: Versions{
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
				Version{Major: "1", Minor: "1", Patchlevel: "0"},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
		{name: "higher patchlevel",
			v: Versions{
				Version{Major: "1", Minor: "0", Patchlevel: "1"},
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
			},
			args: args{i: 0, j: 1},
			want: true,
		},
		{name: "lower patchlevel",
			v: Versions{
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
				Version{Major: "1", Minor: "0", Patchlevel: "1"},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
		{name: "same version with PreRelease",
			v: Versions{
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
				Version{Major: "1", Minor: "0", Patchlevel: "0", PreRelease: "rc1"},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
		{name: "same version with Metadata",
			v: Versions{
				Version{Major: "1", Minor: "0", Patchlevel: "0"},
				Version{Major: "1", Minor: "0", Patchlevel: "0", Metadata: "fefe"},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Versions.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}
