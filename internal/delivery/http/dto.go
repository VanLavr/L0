package http

type Response struct {
	Error   string `json:"error_details"`
	Content any    `json:"data"`
}
