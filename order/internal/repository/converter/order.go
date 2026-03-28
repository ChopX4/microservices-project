package converter

import (
	"github.com/ChopX4/raketka/order/internal/model"
	repoModel "github.com/ChopX4/raketka/order/internal/repository/model"
)

func OrderByUUIDToRepo(model model.OrderByUUID) repoModel.OrderByUUID {
	paymentMethod := PaymentMethodToRepo(model.PaymentMethod)
	if paymentMethod == "" {
		paymentMethod = "UNKNOWN"
	}

	status := OrderStatusToRepo(model.Status)
	if status == "" {
		status = "PENDING_PAYMENT"
	}

	return repoModel.OrderByUUID{
		OrderUUID:       model.OrderUUID,
		UserUUID:        model.UserUUID,
		PartUuids:       model.PartUuids,
		TotalPrice:      model.TotalPrice,
		TransactionUUID: model.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
	}
}

func OrderByUUIDToModel(repo repoModel.OrderByUUID) model.OrderByUUID {
	return model.OrderByUUID{
		OrderUUID:       repo.OrderUUID,
		UserUUID:        repo.UserUUID,
		PartUuids:       repo.PartUuids,
		TotalPrice:      repo.TotalPrice,
		TransactionUUID: repo.TransactionUUID,
		PaymentMethod:   PaymentMethodToModel(repo.PaymentMethod),
		Status:          OrderStatusToModel(repo.Status),
	}
}

func PaymentMethodToRepo(method model.PaymentMethod) repoModel.PaymentMethod {
	return repoModel.PaymentMethod(method)
}

func OrderStatusToRepo(status model.OrderStatus) repoModel.OrderStatus {
	return repoModel.OrderStatus(status)
}

func PaymentMethodToModel(method repoModel.PaymentMethod) model.PaymentMethod {
	return model.PaymentMethod(method)
}

func OrderStatusToModel(status repoModel.OrderStatus) model.OrderStatus {
	return model.OrderStatus(status)
}
