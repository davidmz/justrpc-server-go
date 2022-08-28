package justrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Version struct {
	Major uint
	Minor uint
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

func (v *Version) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}
	v1, err := ParseVersion(text)
	if err != nil {
		return err
	}
	v.Major = v1.Major
	v.Minor = v1.Minor
	return nil
}

func (v *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

var vRe = regexp.MustCompile(`^([1-9]\d*).(0|[1-9]\d*)$`)

func ParseVersion(text string) (*Version, error) {
	m := vRe.FindStringSubmatch(text)
	if m == nil {
		return nil, errors.New("invalid version format, expected {major}.{minor}")
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	return &Version{uint(major), uint(minor)}, nil
}
