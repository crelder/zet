package parse

import (
	"reflect"
	"testing"
)

func TestReference(t *testing.T) {
	tests := []struct {
		name       string
		references string
		want       []string
	}{
		{"Happy case",
			`@article{pandza2010,
	author = {Pandza, Krsto and Thorpe, Richard},
	date-added = {2019-04-25 08:49:22 +0200},
	date-modified = {2019-04-25 08:49:24 +0200},
	journal = {British Journal of Management},
	number = {1},
	pages = {171--186},
	publisher = {Wiley Online Library},
	title = {Management as design, but what kind of design? An appraisal of the design science analogy for management},
	volume = {21},
	year = {2010}}`,
			[]string{"pandza2010"},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reference(tt.references); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reference() = %v, want %v", got, tt.want)
			}
		})
	}
}
