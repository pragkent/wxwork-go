package wxwork

import (
	"errors"

	"github.com/pragkent/wxwork-go/internal"
)

// TargetSet determines the targets of an action.
type TargetSet struct {
	users   internal.StringSet
	parties internal.IntSet
	tags    internal.IntSet
}

// AddUser adds an user to target set.
func (t *TargetSet) AddUser(user string) *TargetSet {
	t.users = append(t.users, user)
	return t
}

// AddUsers adds users to target set.
func (t *TargetSet) AddUsers(users []string) *TargetSet {
	for _, u := range users {
		t.users = append(t.users, u)
	}

	return t
}

// AddParty adds party to target set.
func (t *TargetSet) AddParty(party int) *TargetSet {
	t.parties = append(t.parties, party)
	return t
}

// AddParties adds parties to target set.
func (t *TargetSet) AddParties(parties []int) *TargetSet {
	for _, p := range parties {
		t.parties = append(t.parties, p)
	}

	return t
}

// AddTag adds tag to target set.
func (t *TargetSet) AddTag(tag int) *TargetSet {
	t.tags = append(t.tags, tag)
	return t
}

// AddTags adds tags to target set.
func (t *TargetSet) AddTags(tags []int) *TargetSet {
	for _, p := range tags {
		t.tags = append(t.tags, p)
	}

	return t
}

// Validate validates target set.
//
// One of users / parties / tags must be set.
func (t TargetSet) Validate() error {
	if len(t.users) == 0 && len(t.parties) == 0 && len(t.tags) == 0 {
		return errors.New("wxwork: empty target set")
	}

	return nil
}
