package storage

type UserResp struct {
	Name  string     `json:"n,omitempty"`
	Links []UserResp `json:"l,omitempty"`

	Point map[string]float64 `json:"p,omitempty"`
}

func (u UserResp) Walk(us []UserResp, depth int, fn func(u UserResp, depth int)) {
	for _, u := range us {
		fn(u, depth)
		if len(u.Links) > 0 {
			nextDepth := depth + 1
			u.Walk(u.Links, nextDepth, fn)
		}
	}
}
