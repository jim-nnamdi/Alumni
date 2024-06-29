package model

type ConnectionRequest struct {
	Id        string `json:"connectionid"` /* relative connection id */
	Requester string `json:"requester"`    /* session key for request initiator */
	Accepter  string `json:"accepter"`     /* session key for request receiver */
}

type Connection struct {
	ConnectionRequestId string   `json:"connection_request_id"` /* connection initialized */
	ConnectionExpired   bool     `json:"connection_expired"`    /* did connection expire */
	ConnectedUsers      []string `json:"connected_users"`       /* which users are involved */
}
