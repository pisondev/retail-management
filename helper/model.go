package helper

import (
	"retail-management/model/domain"
	"retail-management/model/web"
)

func ToUserRegisterResponse(user domain.User) web.UserRegisterResponse {
	return web.UserRegisterResponse{
		UserID:   user.UserID,
		Username: user.Username,
	}
}

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
	}
}

func ToUserResponses(users []domain.User) []web.UserResponse {
	userResponses := make([]web.UserResponse, 0)

	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return userResponses
}

func ToRoleResponse(role domain.Role) web.RoleResponse {
	return web.RoleResponse{
		RoleID:   role.RoleID,
		RoleName: role.RoleName,
	}
}

func ToRoleResponses(roles []domain.Role) []web.RoleResponse {
	roleResponses := make([]web.RoleResponse, 0)

	for _, role := range roles {
		roleResponses = append(roleResponses, ToRoleResponse(role))
	}
	return roleResponses
}

func ToCategoryResponse(category domain.Category) web.CategoryResponse {
	return web.CategoryResponse{
		CategoryID:   category.CategoryID,
		CategoryName: category.CategoryName,
	}
}

func ToCategoryResponses(categories []domain.Category) []web.CategoryResponse {
	categoryResponses := make([]web.CategoryResponse, 0)

	for _, category := range categories {
		categoryResponses = append(categoryResponses, ToCategoryResponse(category))
	}
	return categoryResponses
}

func ToSupplierResponse(supplier domain.Supplier) web.SupplierResponse {
	return web.SupplierResponse{
		SupplierID:   supplier.SupplierID,
		SupplierName: supplier.SupplierName,
		PhoneNumber:  supplier.PhoneNumber,
		Email:        supplier.Email,
	}
}

func ToSupplierResponses(suppliers []domain.Supplier) []web.SupplierResponse {
	supplierResponses := make([]web.SupplierResponse, 0)

	for _, supplier := range suppliers {
		supplierResponses = append(supplierResponses, ToSupplierResponse(supplier))
	}
	return supplierResponses
}
