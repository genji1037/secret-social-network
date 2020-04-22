package service

import (
	"secret-social-network/storage"
)

func Link(uid1, uid2 string, point float64) error {
	u1 := &storage.User{
		Name: uid1,
	}
	err := u1.QueryByName()
	if err != nil {
		return err
	}

	u2 := &storage.User{
		Name: uid2,
	}
	err = u2.QueryByName()
	if err != nil {
		return err
	}

	if len(u1.Uid) == 0 {
		u1 = &storage.User{
			Name: uid1,
		}
		err := u1.Save()
		if err != nil {
			return err
		}
	}

	if len(u2.Uid) == 0 {
		u2 = &storage.User{
			Name: uid2,
		}
		err := u2.Save()
		if err != nil {
			return err
		}
	}

	u1.Links = []storage.User{
		{
			Uid:   u2.Uid,
			Point: point,
		},
	}
	err = u1.Save()
	if err != nil {
		return err
	}

	u2.Links = []storage.User{
		{
			Uid:   u1.Uid,
			Point: point,
		},
	}
	err = u2.Save()
	return err
}
