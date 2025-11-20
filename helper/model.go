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

func ToProductResponse(product domain.Product) web.ProductResponse {
	return web.ProductResponse{
		ProductID:     product.ProductID,
		ProductName:   product.ProductName,
		PurchasePrice: product.PurchasePrice,
		SellingPrice:  product.SellingPrice,
		StockQuantity: product.StockQuantity,
		CategoryID:    product.CategoryID,
		SupplierID:    product.SupplierID,
	}
}

func ToProductResponses(products []domain.Product) []web.ProductResponse {
	productResponses := make([]web.ProductResponse, 0)

	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}
	return productResponses
}

func ToProductUpdateResponse(product domain.ProductUpdate) web.ProductUpdateResponse {
	return web.ProductUpdateResponse{
		ProductID:     product.ProductID,
		ProductName:   product.ProductName,
		PurchasePrice: product.PurchasePrice,
		SellingPrice:  product.SellingPrice,
		StockQuantity: product.StockQuantity,
		CategoryID:    product.CategoryID,
		SupplierID:    product.SupplierID,
	}
}

func ToProductUpdateResponses(products []domain.ProductUpdate) []web.ProductUpdateResponse {
	productResponses := make([]web.ProductUpdateResponse, 0)

	for _, product := range products {
		productResponses = append(productResponses, ToProductUpdateResponse(product))
	}
	return productResponses
}

func ToInventoryLogResponse(inventoryLog domain.InventoryLog) web.InventoryLogResponse {
	return web.InventoryLogResponse{
		LogID:          inventoryLog.LogID,
		ProductID:      inventoryLog.ProductID,
		UserID:         inventoryLog.UserID,
		ChangeQuantity: inventoryLog.ChangeQuantity,
		Reason:         inventoryLog.Reason,
		CreatedAt:      inventoryLog.CreatedAt,
	}
}

func ToTransactionItemResponse(detail domain.TransactionDetailWithProduct) web.TransactionItemResp {
	return web.TransactionItemResp{
		ProductID:   detail.ProductID,
		ProductName: detail.ProductName,
		Quantity:    detail.Quantity,
		Price:       detail.PriceAtSale,
		SubTotal:    detail.SubTotal,
	}
}

func ToTransactionItemResponses(details []domain.TransactionDetailWithProduct) []web.TransactionItemResp {
	var responses []web.TransactionItemResp
	for _, detail := range details {
		responses = append(responses, ToTransactionItemResponse(detail))
	}
	return responses
}

func ToTransactionResponse(transaction domain.TransactionWithTotal, items []web.TransactionItemResp) web.TransactionResponse {
	return web.TransactionResponse{
		TransactionID: transaction.TransactionID,
		UserID:        transaction.UserID,
		TotalAmount:   transaction.TotalAmount,
		CreatedAt:     transaction.CreatedAt,
		Items:         items,
	}
}

func ToTransactionResponses(transactions []domain.TransactionWithTotal) []web.TransactionResponse {
	var responses []web.TransactionResponse
	for _, trx := range transactions {
		responses = append(responses, ToTransactionResponse(trx, nil))
	}
	return responses
}
