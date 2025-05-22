package handler

type DiscountHandler interface {
}

type discountHandler struct {
}

func NewDiscountHandler() discountHandler {
	return discountHandler{}
}
