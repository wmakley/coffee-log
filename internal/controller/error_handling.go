package controller

/// errorResponse TODO
func errorResponse(err error) string {
	return err.Error()
}
