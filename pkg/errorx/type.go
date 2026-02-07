package errorx

var (
	ErrInternal      = New(500, "Internal Server Error")
	ErrNotFound      = New(404, "Not Found")
	ErrForbidden     = New(403, "Forbidden")
	ErrUnauthorized  = New(401, "Unauthorized")
	ErrParamsInvalid = New(400, "Params Invalid")
)
