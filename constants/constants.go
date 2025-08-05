package constants

import "net/http"

const (
	OK        = http.StatusOK
	BadReq    = http.StatusBadRequest
	UnAuth    = http.StatusUnauthorized
	Forbidden = http.StatusForbidden
	Empty     = http.StatusNotFound
	Err       = http.StatusInternalServerError
)

const (
	OrderStatusDone       = "done"
	OrderStatusCanceled   = "canceled"
	OrderStatusInProgress = "in progress"
)
