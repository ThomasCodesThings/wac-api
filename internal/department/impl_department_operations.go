package department

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
  "go.mongodb.org/mongo-driver/bson"
  "github.com/ThomasCodesThings/wac-api/internal/db_service"
  "strconv"
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

	// GetDepartments - Provides the list of departments
	GetDepartments(ctx *gin.Context)

	// AddDepartment - Add new department
	AddDepartment(ctx *gin.Context)

	// DeleteDepartment - Delete department
	DeleteDepartment(ctx *gin.Context)

}

// implDepartmentAPI struct definition
type implDepartmentAPI struct{
}

// NewDepartmentAPI creates a new instance of DepartmentAPI
func NewDepartmentAPI() DepartmentAPI {
	return &implDepartmentAPI{
	}
}

func (this *implDepartmentAPI) AddRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.Handle(http.MethodPost, "/operations", this.AddOperation)
	routerGroup.Handle(http.MethodDelete, "/operations/:operationId", this.DeleteOperation)
	routerGroup.Handle(http.MethodPut, "/operations/:operationId", this.EditOperation)
	routerGroup.Handle(http.MethodGet, "/operations/:operationId", this.GetOperation)
	routerGroup.Handle(http.MethodGet, "/operations", this.GetOperations)

	routerGroup.Handle(http.MethodGet, "/departments", this.GetDepartments)
	routerGroup.Handle(http.MethodGet, "/departments/:departmentId", this.GetDepartment)
	routerGroup.Handle(http.MethodGet, "/departments/:departmentId/operations", this.GetDepartmentOperations)
	routerGroup.Handle(http.MethodPost, "/departments/", this.AddDepartment)
	routerGroup.Handle(http.MethodDelete, "/departments/:departmentId", this.DeleteDepartment) 
}

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_conditions.go
func (this *implDepartmentAPI) GetOperations(ctx *gin.Context) {
	db := db_service.NewMongoService[Operation](db_service.MongoServiceConfig{})
    documents, err := db.FindDocuments(ctx, "operation")

    switch err {
    case nil:
        ctx.JSON(http.StatusOK, documents)
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Ambulance not found",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to get list of operations",
                "error":   err.Error(),
            })
    }
}

func (this *implDepartmentAPI) GetOperation(ctx *gin.Context) {
	db := db_service.NewMongoService[Operation](db_service.MongoServiceConfig{})
	//ctx.Set("db_service", dbService)
	
	operationId := ctx.Param("operationId")
	document, err := db.FindDocument(ctx, operationId, "operation")

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, document)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Operation not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to get operation",
				"error":   err.Error(),
			})
	}
}

func (this *implDepartmentAPI) AddOperation(ctx *gin.Context) {
	db := db_service.NewMongoService[Operation](db_service.MongoServiceConfig{})
	
	duration, err := strconv.Atoi(ctx.PostForm("duration"))

	var operation Operation
	operation.Id = uuid.New().String()
	operation.Firstname = ctx.PostForm("firstname")
	operation.Lastname = ctx.PostForm("lastname")
	operation.Department = ctx.PostForm("department")
	operation.AppointmentDate = ctx.PostForm("appointmentDate")
	operation.Duration = int32(duration)

	err = db.CreateDocument(ctx, operation.Id, &operation, "operation")

	switch err {
	case nil:
		ctx.JSON(http.StatusCreated, operation)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Operation already exists",
				"error":   err.Error(),
			},
		)

	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create operation",
				"error":   err.Error(),
			})
	}
}

func (this *implDepartmentAPI) EditOperation(ctx *gin.Context) {
	db := db_service.NewMongoService[Operation](db_service.MongoServiceConfig{})
	
	operationId := ctx.Param("operationId")
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	department := ctx.PostForm("department")
	appointmentDate := ctx.PostForm("appointmentDate")
	duration, err := strconv.Atoi(ctx.PostForm("duration"))

	var operation Operation
	operation.Id = operationId
	operation.Firstname = firstname
	operation.Lastname = lastname
	operation.Department = department
	operation.AppointmentDate = appointmentDate
	operation.Duration = int32(duration)

	err = db.UpdateDocument(ctx, operationId, &operation, "operation")

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, operation)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Operation not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update operation",
				"error":   err.Error(),
			})
	}
}

func (this *implDepartmentAPI) DeleteOperation(ctx *gin.Context) {
	db := db_service.NewMongoService[Operation](db_service.MongoServiceConfig{})
	
	operationId := ctx.Param("operationId")
	err := db.DeleteDocument(ctx, operationId, "operation")

	switch err {
	case nil:
		ctx.JSON(http.StatusNoContent, nil)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Operation not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete operation",
				"error":   err.Error(),
			})
	}
}	

func (this *implDepartmentAPI) GetDepartmentOperations(ctx *gin.Context) {
	
	departmentId := ctx.Param("departmentId")
	department, err := db_service.NewMongoService[Department](db_service.MongoServiceConfig{}).FindDocument(ctx, departmentId, "department")

	filter := bson.D{{Key: "department", Value: department.Name}}
	operations, err := db_service.NewMongoService[Operation](db_service.MongoServiceConfig{}).FindDocumentsByQuery(ctx, filter, "operation")

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, operations)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Operations not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to get list of operations by department",
				"error":   err.Error(),
			})
	}
}

func (this *implDepartmentAPI) GetDepartment(ctx *gin.Context) {
	db := db_service.NewMongoService[Department](db_service.MongoServiceConfig{})

	departmentId := ctx.Param("departmentId")
	department, err := db.FindDocument(ctx, departmentId, "department")

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, department)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Department not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to get department",
				"error":   err.Error(),
			})
	}
}

func (this *implDepartmentAPI) GetDepartments(ctx *gin.Context) {
	db := db_service.NewMongoService[Department](db_service.MongoServiceConfig{})

	departments, err := db.FindDocuments(ctx, "department")

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, departments)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Departments not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to get list of departments",
				"error":   err.Error(),
			})
	}
}

func (this *implDepartmentAPI) AddDepartment(ctx *gin.Context) {
	db := db_service.NewMongoService[Department](db_service.MongoServiceConfig{})

	var department Department
	department.Id = uuid.New().String()
	department.Name = ctx.PostForm("name")
	department.PricePerHour, _ = strconv.ParseFloat(ctx.PostForm("pricePerHour"), 32)

	err := db.CreateDocument(ctx, department.Id, &department, "department")

	switch err {
	case nil:
		ctx.JSON(http.StatusCreated, department)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Department already exists",
				"error":   err.Error(),
			},
		)

	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create department",
				"error":   err.Error(),
			})
	}
}

func (this *implDepartmentAPI) DeleteDepartment(ctx *gin.Context) {
	db := db_service.NewMongoService[Department](db_service.MongoServiceConfig{})
	
	departmentId := ctx.Param("departmentId")
	err := db.DeleteDocument(ctx, departmentId, "department")

	switch err {
	case nil:
		ctx.JSON(http.StatusNoContent, nil)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Department not found",
				"error":   err.Error(),
			},
		)

	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete department",
				"error":   err.Error(),
			})			
	}
}

