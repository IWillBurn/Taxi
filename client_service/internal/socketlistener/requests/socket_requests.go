package socketrequests

type SocketRequest struct {
	Key  string            `json:"key"`
	Data map[string]string `json:"data"`
}
