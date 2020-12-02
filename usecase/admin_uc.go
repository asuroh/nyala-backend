package usecase

import (
	"errors"
	"kriyapeople/helper"
	"kriyapeople/model"
	"kriyapeople/pkg/bcrypt"
	"kriyapeople/pkg/interfacepkg"
	"kriyapeople/pkg/logruslogger"
	"kriyapeople/pkg/str"
	"kriyapeople/server/request"
	"kriyapeople/usecase/viewmodel"
	"strings"
	"time"
)

// AdminUC ...
type AdminUC struct {
	*ContractUC
}

// BuildBody ...
func (uc AdminUC) BuildBody(data *model.UserEntity, res *viewmodel.UserVM, isShowPassword bool) {

	res.ID = data.ID
	res.Information.UserName = data.UserName.String
	res.Information.Email = data.Email.String
	res.Information.Password = str.ShowString(isShowPassword, data.Password.String)
	res.RoleID = data.RoleID.String
	res.RoleName = data.Role.Name.String
	res.Information.Status.IsActive = data.Status.Bool
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// Login ...
func (uc AdminUC) Login(data request.UserLoginRequest) (res viewmodel.JwtVM, err error) {
	ctx := "AdminUC.Login"

	if len(data.Password) < 8 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "password_length", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	admin, err := uc.FindByEmail(data.Email, true)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_by_email", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	isMatch := bcrypt.CheckPasswordHash(data.Password, admin.Information.Password)
	if !isMatch {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "invalid_password", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	if !admin.Information.Status.IsActive {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "inactive_admin", uc.ReqID)
		return res, errors.New(helper.InactiveAdmin)
	}

	// Jwe the payload & Generate jwt token
	payload := map[string]interface{}{
		"id":   admin.ID,
		"role": "admin",
	}
	jwtUc := JwtUC{ContractUC: uc.ContractUC}
	err = jwtUc.GenerateToken(payload, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwt", uc.ReqID)
		return res, errors.New(helper.InternalServer)
	}

	return res, err
}

// FindAll ...
func (uc AdminUC) FindAll(search string, page, limit int, by, sort string) (res []viewmodel.UserVM, pagination viewmodel.PaginationVM, err error) {
	ctx := "AdminUC.FindAll"

	if !str.Contains(model.AdminBy, by) {
		by = model.DefaultAdminBy
	}
	if !str.Contains(SortWhitelist, strings.ToLower(sort)) {
		sort = DescSort
	}

	limit = uc.LimitMax(limit)
	limit, offset := uc.PaginationPageOffset(page, limit)

	m := model.NewAdminModel(uc.DB)
	data, count, err := m.FindAll(search, offset, limit, by, sort)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, pagination, err
	}
	pagination = PaginationRes(page, count, limit)

	for _, r := range data {
		temp := viewmodel.UserVM{}
		uc.BuildBody(&r, &temp, false)
		res = append(res, temp)
	}

	return res, pagination, err
}

// FindByID ...
func (uc AdminUC) FindByID(id string, isShowPassword bool) (res viewmodel.UserVM, err error) {
	ctx := "AdminUC.FindByID"

	m := model.NewAdminModel(uc.DB)
	data, err := m.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}
	uc.BuildBody(&data, &res, isShowPassword)

	return res, err
}

// FindByEmail ...
func (uc AdminUC) FindByEmail(email string, isShowPassword bool) (res viewmodel.UserVM, err error) {
	ctx := "AdminUC.FindByEmail"

	m := model.NewAdminModel(uc.DB)
	data, err := m.FindByEmail(email)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	uc.BuildBody(&data, &res, isShowPassword)

	return res, err
}

// CheckDetails ...
func (uc AdminUC) CheckDetails(data *request.UserRequest, oldData *viewmodel.UserVM) (err error) {
	ctx := "AdminUC.CheckDetails"

	admin, _ := uc.FindByEmail(data.Information.Email, false)
	if admin.ID != "" && admin.ID != oldData.ID {
		logruslogger.Log(logruslogger.WarnLevel, data.Information.Email, ctx, "duplicate_email", uc.ReqID)
		return errors.New(helper.DuplicateEmail)
	}

	if data.Information.Password == "" && oldData.Information.Password == "" {
		logruslogger.Log(logruslogger.WarnLevel, data.Information.Email, ctx, "empty_password", uc.ReqID)
		return errors.New(helper.InvalidPassword)
	}

	// Decrypt password input
	if data.Information.Password == "" {
		data.Information.Password = oldData.Information.Password
	}

	// Encrypt password
	data.Information.Password, err = bcrypt.HashPassword(data.Information.Password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "encrypt_password", uc.ReqID)
		return err
	}

	return err
}

// Create ...
func (uc AdminUC) Create(data *request.UserRequest) (res viewmodel.UserVM, err error) {
	ctx := "AdminUC.Create"

	err = uc.CheckDetails(data, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	information := viewmodel.UserDataVM{
		UserName: data.Information.UserName,
		Email:    data.Information.Email,
	}

	now := time.Now().UTC()
	res = viewmodel.UserVM{
		RoleID:      data.RoleID,
		Information: information,
		Data:        interfacepkg.Marshall(data.Information),
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}
	m := model.NewAdminModel(uc.DB)
	res.ID, err = m.Store(res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Update ...
func (uc AdminUC) Update(id string, data *request.UserRequest) (res viewmodel.UserVM, err error) {
	ctx := "AdminUC.Update"

	oldData, err := uc.FindByID(id, true)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_user", uc.ReqID)
		return res, err
	}

	err = uc.CheckDetails(data, &oldData)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	information := viewmodel.UserDataVM{
		UserName: data.Information.UserName,
		Email:    data.Information.Email,
	}

	now := time.Now().UTC()
	res = viewmodel.UserVM{
		RoleID:      data.RoleID,
		Information: information,
		Data:        interfacepkg.Marshall(data.Information),
		UpdatedAt:   now.Format(time.RFC3339),
	}
	m := model.NewAdminModel(uc.DB)
	res.ID, err = m.Update(id, res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Delete ...
func (uc AdminUC) Delete(id string) (res viewmodel.UserVM, err error) {
	ctx := "AdminUC.Delete"

	now := time.Now().UTC()
	m := model.NewAdminModel(uc.DB)
	res.ID, err = m.Destroy(id, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}
