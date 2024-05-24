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

import (
	"time"
)

type Operations struct {

	// Unique id of the patient in waiting list
	PatientId int64 `json:"patientId"`

	// Name of patient in waiting list
	Firstname string `json:"firstname"`

	// Surname of patient in waiting list
	Surname string `json:"surname"`

	// Department of patient in waiting list
	Department string `json:"department"`

	// Date of appointment of patient in waiting list
	AppointmentDate time.Time `json:"appointmentDate"`

	// Duration of appointment of patient in waiting list
	Duration int32 `json:"duration"`
}