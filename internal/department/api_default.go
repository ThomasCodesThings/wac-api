package department

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

type DefaultAPI interface {
	// internal registration of api routes
	addRoutes(routerGroup *gin.RouterGroup)

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

// partial implementation of DefaultAPI - all functions must be implemented in add on files
type implDefaultAPI struct{}

func newDefaultAPI() DefaultAPI {
	return &implDefaultAPI{}
}

func (this *implDefaultAPI) addRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.Handle(http.MethodPost, "/operations", this.AddOperation)
	routerGroup.Handle(http.MethodDelete, "/operations/:operationId", this.DeleteOperation)
	routerGroup.Handle(http.MethodPut, "/operations/:operationId", this.EditOperation)
	routerGroup.Handle(http.MethodGet, "/department/:departmentId", this.GetDepartmentOperations)
	routerGroup.Handle(http.MethodGet, "/operations/:operationId", this.GetOperation)
	routerGroup.Handle(http.MethodGet, "/operations", this.GetOperations)
}

// AddOperation - Add new operation to waiting list
func (this *implDefaultAPI) AddOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// DeleteOperation - Delete operation from waiting list
func (this *implDefaultAPI) DeleteOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// EditOperation - Edit operation in waiting list
func (this *implDefaultAPI) EditOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// GetDepartmentOperations - Provides list of Operations in department
func (this *implDefaultAPI) GetDepartmentOperations(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// GetOperation - Provides the operation waiting list by ID
func (this *implDefaultAPI) GetOperation(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// GetOperations - Provides the operations waiting list
func (this *implDefaultAPI) GetOperations(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotImplemented)
}
