package http

type Response struct {
	Error   string `json:"error details"`
	Content any    `json:"data"`
}
