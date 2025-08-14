package response

type HttpErrorDto struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
