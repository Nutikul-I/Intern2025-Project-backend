package handler

type MerchantHandler interface {
}

type merchantHandler struct {
}

func NewMerchantHandler() merchantHandler {
	return merchantHandler{}
}
