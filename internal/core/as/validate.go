package as

import (
	"fmt"
	"github.com/crelder/zet/internal/core/bl"
	"github.com/crelder/zet/internal/core/port"
)

// Validator shows any inconsistencies the zettelkasten has.
type Validator struct {
	Repo port.Repo
}

func NewValidator(r port.Repo) Validator {
	return Validator{
		Repo: r,
	}
}

// Val returns different types of consistency errors that your zettelkasten has.
// If there are no errors, it returns nil.
func (v Validator) Val() []bl.ValidatorErr {
	z, err := v.Repo.GetZettel()
	if err != nil {
		fmt.Printf("An error ocurred while validating %v", err)
	}

	zk := bl.NewZk(z, nil)

	return zk.Validate()
}
