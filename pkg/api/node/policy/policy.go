package policy

import (
	"juno/pkg/api/node"
	"juno/pkg/api/user"
)

func CanUpdate(u *user.User, n *node.Node) bool {
	return u.ID == n.OwnerID
}

func CanDelete(u *user.User, n *node.Node) bool {
	return u.ID == n.OwnerID
}

func CanRead(u *user.User, n *node.Node) bool {
	return u.ID == n.OwnerID
}

func CanCreate(u *user.User) bool {
	return true
}
