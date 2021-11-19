package controller

func errorResponse(err error) string {
	return err.Error()
}
