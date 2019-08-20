package semver

import (
	"regexp"
	"strconv"
	"strings"

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

	ErrVersionInvalid = errors.New("Invalid version")
)

func Default() Version {
	return Version{
		Major:      defaultMajor,
		Minor:      defaultMinor,
		Patchlevel: defaultPatchlevel,
		PreRelease: defaultPreRelease,
		Metadata:   defaultMetadata,
	}
}

type Version struct {
	Major      string
	Minor      string
	Patchlevel string
	PreRelease string
	Metadata   string

	prefix string
}

type option func(v *Version) option

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
