package usecase

import (
	"nyala-backend/model"
	"nyala-backend/pkg/logruslogger"
	"nyala-backend/server/request"
	"nyala-backend/usecase/viewmodel"
	"time"
)

// OrderDetailUC ...
type OrderDetailUC struct {
	*ContractUC
}

// Create ...
func (uc OrderDetailUC) Create(data *request.OrderDetailRequest) (res viewmodel.OrderDetailVM, err error) {
	ctx := "OrderUC.Create"

	now := time.Now().UTC()
	res = viewmodel.OrderDetailVM{
		OrderID:   data.OrderID,
		ProductID: data.ProductID,
		Qty:       data.Qty,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}

	m := model.NewOrderDetailModel(uc.DB)
	res.OrderDetailID, err = m.Store(res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}
