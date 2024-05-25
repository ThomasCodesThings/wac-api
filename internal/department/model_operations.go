/*
 * Waiting List Api
 *
 * Department List management for Web-In-Cloud system
 *
 * API version: 1.0.0
 * Contact: your_email@example.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package department

type Operation struct {

	// Unique id of the patient in waiting list
	Id string `json:"operationId"`

	// Name of patient in waiting list
	Firstname string `json:"firstname"`

	// Surname of patient in waiting list
	Lastname string `json:"surname"`

	// Department of patient in waiting list
	Department string `json:"department"`

	// Date of appointment of patient in waiting list
	AppointmentDate string `json:"appointmentDate"`

	// Duration of appointment of patient in waiting list
	Duration int32 `json:"duration"`
}

type Department struct {
	Id string `json:"departmentId"`
	// Name of department
	Name string `json:"name"`
}
