package usecase

import (
	"nyala-backend/model"
	"nyala-backend/pkg/logruslogger"
	"nyala-backend/pkg/number"
	"nyala-backend/server/request"
	"nyala-backend/usecase/viewmodel"
	"strconv"
	"time"
)

// OrderUC ...
type OrderUC struct {
	*ContractUC
}

// GenerateCode ...
func (uc OrderUC) GenerateCode(date time.Time) (res string, err error) {
	ctx := "OrderUC.CheckDetails"

	m := model.NewOrderModel(uc.DB)
	count, err := m.Count()
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = "PO-" + strconv.Itoa(count+1) + "/" + number.IntToRoman(int(date.Month())) + "/" + strconv.Itoa(date.Year())

	return res, err
}

// Create ...
func (uc OrderUC) Create(data *request.OrderRequest, customerID string) (res viewmodel.OrderVM, err error) {
	ctx := "OrderUC.Create"

	tx := model.SQLDBTx{DB: uc.DB}
	dbTx, err := tx.TxBegin()
	uc.Tx = dbTx.DB
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "tx", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	orderNumber, err := uc.GenerateCode(now)
	if err != nil {
		uc.Tx.Rollback()
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "generate_order_number", uc.ReqID)
		return res, err
	}

	res = viewmodel.OrderVM{
		CustomerID:      customerID,
		OrderNumber:     orderNumber,
		PaymentMethodID: data.PaymentMethodID,
		OrderDate:       now.Format(time.RFC3339),
		CreatedAt:       now.Format(time.RFC3339),
		UpdatedAt:       now.Format(time.RFC3339),
	}

	m := model.NewOrderModel(uc.DB)
	res.OrderID, err = m.Store(res, now)
	if err != nil {
		uc.Tx.Rollback()
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
			uc.Tx.Rollback()
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
			return res, err
		}

		res.OrderDetail = append(res.OrderDetail, orderDetail)
	}
	uc.Tx.Commit()

	return res, err
}
