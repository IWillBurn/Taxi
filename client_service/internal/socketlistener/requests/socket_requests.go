package socketrequests

type SocketRequest struct {
	Key  string            `json:"type"`
	Data map[string]string `json:"data"`
}
