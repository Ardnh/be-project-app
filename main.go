package main

import (
	"os"
	"project-app/app"
	"project-app/helper"
	"project-app/routes"

	_ "project-app/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

// @title           Project APP API
// @version         1.0
// @description     API Documentation for Project APP API.

// @contact.name   Muhammad Ardan Hilal
// @contact.url    ardn.h79@gmail.com
// @contact.email  ardn.h79@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {

	newApp := fiber.New()
	db := app.DbConnection()
	validate := validator.New(validator.WithRequiredStructEnabled())

	routes.SetupRoutes(newApp, db, validate)

	port := os.Getenv("APP_PORT")
	err := newApp.Listen(port)
	helper.PanicIfError(err)
}

// type EmployeeAttendance struct {
//     EmployeeID             string         `json:"employee_id" gorm:"column:employee_id"`
//     FullName               string         `json:"full_name" gorm:"column:full_name"`
//     JobPosition            string         `json:"job_position" gorm:"column:job_position"`
//     Date                   time.Time      `json:"date" gorm:"column:date"`
//     Shift                  string         `json:"shift" gorm:"column:shift"`
//     ScheduleCheckIn        time.Time      `json:"schedule_check_in" gorm:"column:schedule_check_in"`
//     ScheduleCheckOut       time.Time      `json:"schedule_check_out" gorm:"column:schedule_check_out"`
//     AttendanceCode         sql.NullString `json:"attendance_code" gorm:"column:attendance_code"`
//     TimeOffCode            sql.NullString `json:"time_off_code" gorm:"column:time_off_code"`
//     CheckIn                sql.NullTime   `json:"check_in" gorm:"column:check_in"`
//     CheckOut               sql.NullTime   `json:"check_out" gorm:"column:check_out"`
//     ScheduleBreakStart     time.Time      `json:"schedule_break_start" gorm:"column:schedule_break_start"`
//     ScheduleBreakEnd       time.Time      `json:"schedule_break_end" gorm:"column:schedule_break_end"`
//     ActualBreakStart       sql.NullTime   `json:"actual_break_start" gorm:"column:actual_break_start"`
//     ActualBreakEnd         sql.NullTime   `json:"actual_break_end" gorm:"column:actual_break_end"`
//     LateIn                 string         `json:"late_in" gorm:"column:late_in"`
//     EarlyOut               string         `json:"early_out" gorm:"column:early_out"`
//     ScheduleWorkingHour    string         `json:"schedule_working_hour" gorm:"column:schedule_working_hour"`
//     ActualWorkingHour      string         `json:"actual_working_hour" gorm:"column:actual_working_hour"`
//     RealWorkingHour        string         `json:"real_working_hour" gorm:"column:real_working_hour"`
//     OvertimeDurationBefore string         `json:"overtime_duration_before" gorm:"column:overtime_duration_before"`
//     OvertimeDurationAfter  string         `json:"overtime_duration_after" gorm:"column:overtime_duration_after"`
//     FingerCheckIn          sql.NullTime   `json:"finger_check_in" gorm:"column:finger_check_in"`
//     FingerCheckOut         sql.NullTime   `json:"finger_check_out" gorm:"column:finger_check_out"`
//     UpdatedAt              sql.NullTime   `json:"updated_at" gorm:"column:updated_at"`
// }
