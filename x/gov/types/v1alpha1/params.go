package v1alpha1

import (
	yaml "gopkg.in/yaml.v2"
)

// String satisfies the gogoproto Message interface for CustomParams.
// The rest of the v1alpha1 Go package is unused; only this method is
// retained because gov.pb.go embeds CustomParams as a proto message.
func (cp CustomParams) String() string {
	out, _ := yaml.Marshal(cp)
	return string(out)
}
