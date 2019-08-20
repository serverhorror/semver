package semver_test

import (
	"reflect"
	"sort"
	"testing"

	. "github.com/serverhorror/semver"
)

func TestVersions_Sort(t *testing.T) {
	tests := []struct {
		name string
		have Versions
		want Versions
	}{
		{name: "one entry",
			have: Versions{Default()},
			want: Versions{Default()},
		},
		{name: "major 1 and major 2",
			have: Versions{Version{Major: "2"}, Version{Major: "1"}},
			want: Versions{Version{Major: "1"}, Version{Major: "2"}},
		},
		{name: "major 1 and major 2",
			have: Versions{Version{Major: "2"}, Version{Major: "1"}},
			want: Versions{Version{Major: "1"}, Version{Major: "2"}},
		},
		{name: "major 1 and minor 2",
			have: Versions{Version{Major: "1", Minor: "2"}, Version{Major: "1", Minor: "1"}},
			want: Versions{Version{Major: "1", Minor: "1"}, Version{Major: "1", Minor: "2"}},
		},
		{name: "major 2 and major 2 w minor",
			have: Versions{Version{Major: "2"}, Version{Major: "2", Minor: "1"}},
			want: Versions{Version{Major: "2", Minor: "1"}, Version{Major: "2"}},
		},
		{name: "major 2 and major 2 w minor",
			have: Versions{Version{Major: "2"}, Version{Major: "2", Patchlevel: "1"}},
			want: Versions{Version{Major: "2", Patchlevel: "1"}, Version{Major: "2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.have)
			if !reflect.DeepEqual(tt.have, tt.want) {
				t.Errorf("sort.Sort(%T) error = %v, wantErr %v", tt.have, tt.have, tt.want)
			}
		})
	}
}
