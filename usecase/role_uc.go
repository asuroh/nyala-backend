package usecase

import (
	"kriyapeople/model"
	"kriyapeople/pkg/logruslogger"
	"kriyapeople/pkg/str"
	"kriyapeople/usecase/viewmodel"
	"strings"
)

// RoleUC ...
type RoleUC struct {
	*ContractUC
}

// BuildBody ...
func (uc RoleUC) BuildBody(data *model.RoleEntity, res *viewmodel.RoleVM) {
	res.ID = data.ID
	res.Code = data.Code.String
	res.Name = data.Name.String
	res.Status = data.Status.Bool
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// SelectAll ...
func (uc RoleUC) SelectAll(search, by, sort string) (res []viewmodel.RoleVM, err error) {
	ctx := "RoleUC.SelectAll"

	if !str.Contains(model.RoleBy, by) {
		by = model.DefaultRoleBy
	}
	if !str.Contains(SortWhitelist, strings.ToLower(sort)) {
		sort = DescSort
	}

	m := model.NewRoleModel(uc.DB)
	data, err := m.SelectAll(search, by, sort)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	for _, r := range data {
		temp := viewmodel.RoleVM{}
		uc.BuildBody(&r, &temp)
		res = append(res, temp)
	}

	return res, err
}

// FindByID ...
func (uc RoleUC) FindByID(id string) (res viewmodel.RoleVM, err error) {
	ctx := "RoleUC.FindByID"

	m := model.NewRoleModel(uc.DB)
	data, err := m.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}
	uc.BuildBody(&data, &res)

	return res, err
}

// FindByCode ...
func (uc RoleUC) FindByCode(code string) (res viewmodel.RoleVM, err error) {
	ctx := "RoleUC.FindByCode"

	m := model.NewRoleModel(uc.DB)
	data, err := m.FindByCode(code)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}
	uc.BuildBody(&data, &res)

	return res, err
}
