package usecase

import (
	"errors"
	"nyala-backend/helper"
	"nyala-backend/model"
	"nyala-backend/pkg/logruslogger"
	"nyala-backend/pkg/str"
	"nyala-backend/server/request"
	"nyala-backend/usecase/viewmodel"
	"time"
)

// CustomerUC ...
type CustomerUC struct {
	*ContractUC
}

// BuildBody ...
func (uc CustomerUC) BuildBody(data *model.CustomerEntity, res *viewmodel.CustomerVM, isShowPassword bool) {
	res.CustomerID = data.CustomerID
	res.CustomerName = data.CustomerName
	res.Email = data.Email
	res.PhoneNumber = data.PhoneNumber
	res.Dob = data.Dob
	res.Sex = data.Sex
	res.Password = str.ShowString(isShowPassword, data.Password)
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// Login ...
func (uc CustomerUC) Login(data request.CustomerLoginRequest) (res viewmodel.JwtVM, err error) {
	ctx := "CustomerUC.Login"

	if len(data.Password) < 8 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "password_length", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	customer, err := uc.FindByEmailOrPhoneNumber(data.Username, data.Username, true)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_by_email_or_phone_number", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	data.Password = uc.Hmacsha.Encrypt(data.Password)
	if customer.Password != data.Password {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "invalid_password", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	// Jwe the payload & Generate jwt token
	payload := map[string]interface{}{
		"id": customer.CustomerID,
	}
	jwtUc := JwtUC{ContractUC: uc.ContractUC}
	err = jwtUc.GenerateToken(payload, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwt", uc.ReqID)
		return res, errors.New(helper.InternalServer)
	}

	return res, err
}

// RefreshToken ...
func (uc CustomerUC) RefreshToken(customerID string) (res viewmodel.JwtVM, err error) {
	ctx := "CustomerUC.RefreshToken"

	// Jwe the payload & Generate jwt token
	payload := map[string]interface{}{
		"id": customerID,
	}
	jwtUc := JwtUC{ContractUC: uc.ContractUC}
	err = jwtUc.GenerateToken(payload, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwt", uc.ReqID)
		return res, errors.New(helper.InternalServer)
	}

	return res, err
}

// FindByEmailOrPhoneNumber ...
func (uc CustomerUC) FindByEmailOrPhoneNumber(email, phoneNumber string, isShowPassword bool) (res viewmodel.CustomerVM, err error) {
	ctx := "CustomerUC.FindByEmailOrMail"

	m := model.NewCustomerModel(uc.DB)
	data, err := m.FindByEmailOrPhoneNumber(email, phoneNumber)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	uc.BuildBody(&data, &res, isShowPassword)

	return res, err
}

// FindByID ...
func (uc CustomerUC) FindByID(id string, isShowPassword bool) (res viewmodel.CustomerVM, err error) {
	ctx := "CustomerUC.FindByID"

	m := model.NewCustomerModel(uc.DB)
	data, err := m.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	uc.BuildBody(&data, &res, isShowPassword)

	return res, err
}

// CheckDetails ...
func (uc CustomerUC) CheckDetails(data *request.CustomerRequest, oldData *viewmodel.CustomerVM) (err error) {
	ctx := "CustomerUC.CheckDetails"

	data.SexBool, err = helper.SexStringToBool(data.Sex)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, data.Sex, ctx, "check_sex", uc.ReqID)
		return err
	}

	customer, _ := uc.FindByEmailOrPhoneNumber(data.Email, data.PhoneNumber, false)
	if customer.CustomerID != "" && customer.CustomerID != oldData.CustomerID {
		if customer.Email == data.Email {
			logruslogger.Log(logruslogger.WarnLevel, data.Email, ctx, "duplicate_email", uc.ReqID)
			return errors.New(helper.DuplicateEmail)
		}

		if customer.PhoneNumber == data.PhoneNumber {
			logruslogger.Log(logruslogger.WarnLevel, data.PhoneNumber, ctx, "duplicate_phone_number", uc.ReqID)
			return errors.New(helper.DuplicatePhoneNumber)
		}

		logruslogger.Log(logruslogger.WarnLevel, data.Email+" and "+data.PhoneNumber, ctx, "duplicate_email_and_phone_number", uc.ReqID)
		return errors.New(helper.DuplicateEmail + "_and_" + helper.DuplicatePhoneNumber)
	}

	if data.Password == "" && oldData.Password == "" {
		logruslogger.Log(logruslogger.WarnLevel, data.Password, ctx, "empty_password", uc.ReqID)
		return errors.New(helper.InvalidPassword)
	}

	// Decrypt password input
	if data.Password == "" {
		data.Password = oldData.Password
	}

	// Encrypt password
	data.Password = uc.Hmacsha.Encrypt(data.Password)

	return err
}

// Register ...
func (uc CustomerUC) Register(data *request.CustomerRequest) (res viewmodel.CustomerVM, err error) {
	ctx := "CustomerUC.Register"

	err = uc.CheckDetails(data, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.CustomerVM{
		CustomerName: data.CustomerName,
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		Dob:          data.Dob,
		Sex:          data.SexBool,
		Password:     data.Password,
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
	}
	m := model.NewCustomerModel(uc.DB)
	res.CustomerID, err = m.Store(res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}
