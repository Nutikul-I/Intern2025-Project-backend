package handler

type CategoryHandler interface {
}

type categoryHandler struct {
}

func NewCategoryHandler() categoryHandler {
	return categoryHandler{}
}
