package wxwork

import "github.com/pragkent/wxwork-go/internal"

// TargetSet determines the targets of an action.
type TargetSet struct {
	User  UserSet
	Party PartySet
	Tag   TagSet
}

// User list.
type UserSet internal.StringSet

// Pary list.
type PartySet internal.IntSet

// Tag list.
type TagSet internal.IntSet
