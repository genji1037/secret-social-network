package service

import "secret-social-network/storage"

func Link(uid1, uid2 string, point float64) error {
	return storage.User{}.Link(uid1, uid2, point)
}
