package client

type UIDInfoResp struct {
	UIDs []UIDInfo `json:"uids"`
}

type UIDInfo struct {
	OpenID string `json:"open_id"`
	UID    string `json:"uid"`
}
