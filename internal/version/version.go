package version

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

// Info creates a formattable struct for output
type Info struct {
	Version string `json:"Version,omitempty"`
	Commit  string `json:"Commit,omitempty"`
	Date    string `json:"Date,omitempty"`
}

// New will create a pointer to a new version object
func New(version string, commit string, date string) *Info {
	return &Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}

// ToJSON converts the Info into a JSON String
func (v *Info) ToJSON() string {
	bytes, _ := json.Marshal(v)
	return string(bytes) + "\n"
}

// ToYAML converts the Info into a JSON String
func (v *Info) ToYAML() string {
	bytes, _ := yaml.Marshal(v)
	return string(bytes)
}