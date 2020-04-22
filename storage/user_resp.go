package storage

type UserResp struct {
	Name  string     `json:"name,omitempty"`
	Links []UserResp `json:"links,omitempty"`

	Point map[string]float64 `json:"links|point,omitempty"`
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
