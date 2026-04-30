package method

import "github.com/ChopX4/raketka/payment/internal/model"

func (s *service) validatePayRequest(req model.PayOrderRequest) error {
	if !model.IsValidUUID(req.OrderUuid) {
		return model.ErrBadRequest
	}

	if !model.IsValidUUID(req.UserUuid) {
		return model.ErrBadRequest
	}

	if !req.PaymentMethod.IsValid() {
		return model.ErrBadRequest
	}

	return nil
}
