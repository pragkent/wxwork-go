package wxwork

import (
	"errors"

	"github.com/pragkent/wxwork-go/internal"
)

// TargetSet determines the targets of an action.
type TargetSet struct {
	Users   UserSet
	Parties PartySet
	Tags    TagSet
}

func (t *TargetSet) Validate() error {
	if t == nil {
		return errors.New("wxwork: empty target set")
	}

	if len(t.Users) == 0 && len(t.Parties) == 0 && len(t.Tags) == 0 {
		return errors.New("wxwork: empty target set")
	}

	return nil
}

// User list.
type UserSet internal.StringSet

// Party list.
type PartySet internal.IntSet

// Tag list.
type TagSet internal.IntSet
