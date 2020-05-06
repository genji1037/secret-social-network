package service

import "secret-social-network/dgraph"

func Link(uid1, uid2 string, point float64) error {
	return dgraph.User{}.Link(uid1, uid2, point)
}
