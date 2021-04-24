package employees

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Datos de conexión con Base de Datos
var server = "localhost"
var port = 1433
var user = ""
var password = ""
var database = "Northwind"

var db *sql.DB

func connectionSql(c echo.Context) error {
	var err error

	// Create connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
		return echo.ErrInternalServerError
	}
	log.Printf("Connected!\n")

	// Read employees
	employeeData, err := ReadEmployees()
	if err != nil {
		log.Fatal("Error reading Employees: ", err.Error())
	}
	// Close the database connection pool after program executes
	defer db.Close()

	return c.JSON(http.StatusOK, employeeData)
}

// ReadEmployees reads all employee records
func ReadEmployees() ([]TUserResponse, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsql := fmt.Sprintf("SELECT EmployeeID, FirstName, Title FROM Employees;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []TUserResponse
	var element TUserResponse
	
	// Iterate through the result set.
	for rows.Next() {
		var firstName, title string
		var employeeId int

		// Get values from row.
		err := rows.Scan(&employeeId, &firstName, &title)
		if err != nil {
			return nil, err
		}
		element = TUserResponse{
			EmployeeID: employeeId,
			FirstName: firstName,
			Title: title,
		}
		res = append(res, element)

	}
	return res, nil
}