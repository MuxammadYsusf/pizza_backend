package util

import "net/http"

const (
	HTTPOK        = http.StatusOK
	HTTPBadReq    = http.StatusBadRequest
	HTTPUnAuth    = http.StatusUnauthorized
	HTTPForbidden = http.StatusForbidden
	HTTPNotFound  = http.StatusNotFound
	HTTPServerErr = http.StatusInternalServerError
)

const (
	STATUS_NEW         = "NEW"
	STATUS_PAID        = "PAID"
	STATUS_IN_PROGRESS = "IN_PROGRESS"
	STATUS_DELIVERING  = "DELIVERING"
	STATUS_CANCELED    = "CANCELED"
	STATUS_DONE        = "DONE"
)

const (
	ROLE_ADMIN = "ADMIN"
	ROLE_USER  = "USER"
)	
