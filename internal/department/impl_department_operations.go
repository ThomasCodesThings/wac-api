package department

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

type DepartmentAPI interface {
	// internal registration of api routes
	AddRoutes(routerGroup *gin.RouterGroup)

	// AddOperation - Add new operation to waiting list
	AddOperation(ctx *gin.Context)

	// DeleteOperation - Delete operation from waiting list
	DeleteOperation(ctx *gin.Context)

	// EditOperation - Edit operation in waiting list
	EditOperation(ctx *gin.Context)

	// GetDepartmentOperations - Provides list of Operations in department
	GetDepartmentOperations(ctx *gin.Context)

	// GetOperation - Provides the operation waiting list by ID
	GetOperation(ctx *gin.Context)

	// GetOperations - Provides the operations waiting list
	GetOperations(ctx *gin.Context)
}

// implDepartmentAPI struct definition
type implDepartmentAPI struct{}

// NewDepartmentAPI creates a new instance of DepartmentAPI
func NewDepartmentAPI() DepartmentAPI {
	return &implDepartmentAPI{}
}

func (this *implDepartmentAPI) AddRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.Handle(http.MethodPost, "/operations", this.AddOperation)
	routerGroup.Handle(http.MethodDelete, "/operations/:operationId", this.DeleteOperation)
	routerGroup.Handle(http.MethodPut, "/operations/:operationId", this.EditOperation)
	routerGroup.Handle(http.MethodGet, "/department/:departmentId", this.GetDepartmentOperations)
	routerGroup.Handle(http.MethodGet, "/operations/:operationId", this.GetOperation)
	routerGroup.Handle(http.MethodGet, "/operations", this.GetOperations)
}

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_conditions.go
func (this *implDepartmentAPI) GetOperations(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadGateway)
}

func (this *implDepartmentAPI) GetOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadGateway)
}

func (this *implDepartmentAPI) AddOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadGateway)
}

func (this *implDepartmentAPI) EditOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadGateway)
}

func (this *implDepartmentAPI) DeleteOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadGateway)
}

func (this *implDepartmentAPI) GetDepartmentOperations(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadGateway)
}
