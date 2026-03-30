package handlers

type CheckEmailResponse struct {
	Exists bool `json:"exists"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
