package api

// IndexResponse defines an HTTP response object
type IndexResponse struct {
}

// ExecuteRequest defines an HTTP response object
type ExecuteRequest struct {
	Args string `json:"args"`
}

// ExecuteResponse defines an HTTP response object
type ExecuteResponse struct {
	Centry   string `json:"centry"`
	Result   string `json:"result"`
	ExitCode int    `json:"exitCode"`
}
