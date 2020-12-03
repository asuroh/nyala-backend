package usecase

import (
	"fmt"
	"nyala-backend/model"
	"nyala-backend/pkg/logruslogger"
	"nyala-backend/server/request"
	"nyala-backend/usecase/viewmodel"
	"time"
)

// OrderUC ...
type OrderUC struct {
	*ContractUC
}

// CheckDetails ...
func (uc OrderUC) CheckDetails(data *request.OrderRequest, oldData *viewmodel.OrderVM) (err error) {
	ctx := "OrderUC.CheckDetails"
	fmt.Println(ctx)
	data.OrderNumber = "PO-123/IX/2020"
	return err
}

// Create ...
func (uc OrderUC) Create(data *request.OrderRequest, customerID string) (res viewmodel.OrderVM, err error) {
	ctx := "OrderUC.Create"

	err = uc.CheckDetails(data, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.OrderVM{
		CustomerID:      customerID,
		OrderNumber:     data.OrderNumber,
		PaymentMethodID: data.PaymentMethodID,
		OrderDate:       now.Format(time.RFC3339),
		CreatedAt:       now.Format(time.RFC3339),
		UpdatedAt:       now.Format(time.RFC3339),
	}

	m := model.NewOrderModel(uc.DB)
	res.OrderID, err = m.Store(res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	orderDetailUC := OrderDetailUC{ContractUC: uc.ContractUC}
	for _, row := range data.OrderDetail {
		detail := request.OrderDetailRequest{
			OrderID:   res.OrderID,
			ProductID: row.ProductID,
			Qty:       row.Qty,
		}

		orderDetail, err := orderDetailUC.Create(&detail)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
			return res, err
		}

		res.OrderDetail = append(res.OrderDetail, orderDetail)
	}

	return res, err
}
