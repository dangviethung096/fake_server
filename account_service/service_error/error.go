package service_error

import "core"

var (
	ERROR_NIL_PARAM  = core.NewError(0, "Nil param")
	ERROR_QUERY_FAIL = core.NewError(1, "Query fail")
)
