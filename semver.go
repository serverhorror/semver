package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	separator           = "."
	preReleaseSeparator = "-"
	metadataSeparator   = "+"
)

var (
	defaultMajor      = "0"
	defaultMinor      = "0"
	defaultPatchlevel = "0"
	defaultPreRelease = ""
	defaultMetadata   = ""
	defaultBuildTime  = "<unknown>" // time.Now().Format(time.RFC3339Nano)

	ErrVersionInvalid = errors.New("Invalid version")
)

func Default() Version {
	return Version{
		Major:      defaultMajor,
		Minor:      defaultMinor,
		Patchlevel: defaultPatchlevel,
		PreRelease: defaultPreRelease,
		Metadata:   defaultMetadata,

		buildTime: defaultBuildTime,
	}
}

type Version struct {
	Major      string
	Minor      string
	Patchlevel string
	PreRelease string
	Metadata   string

	buildTime string
	prefix    string
}

type option func(v *Version) option

// BuildTime sets the Versions prefix to p.
func BuildTime(t time.Time) option {
	return func(v *Version) option {
		previous := v.buildTime
		v.buildTime = t.Format(time.RFC3339Nano)
		return Prefix(previous)
	}
}

// Prefix sets the Versions prefix to p.
func Prefix(p string) option {
	return func(v *Version) option {
		previous := v.prefix
		v.prefix = p
		return Prefix(previous)
	}
}

// Option sets the options specified.
// It returns an option to restore the last arg's previous value.
func (v *Version) Option(opts ...option) (previous option) {
	for _, opt := range opts {
		previous = opt(v)
	}
	return previous
}

type Validator interface {
	Validate() error
}

func (v Version) Validate() error {
	const pattern = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	re := regexp.MustCompile(pattern)
	match := re.MatchString(v.String())
	if !match {
		return ErrVersionInvalid
	}
	return nil
}

func (v Version) VerboseString() string {
	major := fmt.Sprintf("Major: %#v", v.Major)
	minor := fmt.Sprintf("Minor: %#v", v.Minor)
	patchlevel := fmt.Sprintf("Patchlevel: %#v", v.Patchlevel)
	preRelease := fmt.Sprintf("Pre Release: %#v", v.PreRelease)
	metadata := fmt.Sprintf("Metadata: %#v", v.Metadata)
	buildTime := fmt.Sprintf("Build Time: %#v", v.buildTime)

	elems := []string{major, minor, patchlevel, preRelease, metadata, buildTime}

	return fmt.Sprintf("%s\n", strings.Join(elems, "\n"))
}
func (v Version) String() string {
	vers := strings.Join([]string{v.Major, v.Minor, v.Patchlevel}, separator)
	if v.PreRelease != "" {
		vers = strings.Join([]string{vers, v.PreRelease}, preReleaseSeparator)
	}
	if v.Metadata != "" {
		vers = strings.Join([]string{vers, v.Metadata}, metadataSeparator)
	}
	return strings.Join([]string{v.prefix, vers}, "")
}

type Versions []Version

// Len is the number of elements in the collection.
func (v Versions) Len() int {
	return len(v)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (v Versions) Less(i int, j int) bool {
	iVersion := v[i]
	iVersionMajor, err := strconv.ParseInt(iVersion.Major, 10, 64)
	if err != nil {
		return true
	}
	iVersionMinor, err := strconv.ParseInt(iVersion.Minor, 10, 64)
	if err != nil {
		return true
	}
	iVersionPatchlevel, err := strconv.ParseInt(iVersion.Patchlevel, 10, 64)
	if err != nil {
		return true
	}

	jVersion := v[j]
	jVersionMajor, err := strconv.ParseInt(jVersion.Major, 10, 64)
	if err != nil {
		return true
	}
	jVersionMinor, err := strconv.ParseInt(jVersion.Minor, 10, 64)
	if err != nil {
		return true
	}
	jVersionPatchlevel, err := strconv.ParseInt(jVersion.Patchlevel, 10, 64)
	if err != nil {
		return true
	}

	if iVersionMajor > jVersionMajor {
		return true
	}

	if iVersionMinor > jVersionMinor {
		return true
	}
	if iVersionPatchlevel > jVersionPatchlevel {
		return true
	}

	return false
}

// Swap swaps the elements with indexes i and j.
func (v Versions) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}
