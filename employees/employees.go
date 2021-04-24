package employees

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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
	employeeData, err := readEmployees()
	if err != nil {
		log.Fatal("Error reading Employees: ", err.Error())
	}
	// Close the database connection pool after program executes
	defer db.Close()

	return c.JSON(http.StatusOK, employeeData)
}

// ReadEmployees reads all employee records
func readEmployees() ([]TUserResponse, error) {
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

func signup(c echo.Context) error {
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

	// TODO: CREAR USUARIO

	ctx := context.Background()

	// Check if database is alive.
	errPing := db.PingContext(ctx)
	if errPing != nil {
		return echo.ErrServiceUnavailable
	}

	u := new(TUserRequest)
    if err = c.Bind(u); err != nil {
      return echo.ErrBadRequest
    }

	hashedPswr := GetMD5Hash(u.Password)
	tsql := fmt.Sprintf("INSERT INTO Users VALUES ('" + u.Username + "','" + hashedPswr + "' );")

	_, sqlerr := db.Exec(tsql)
	if sqlerr != nil {
		return echo.ErrInternalServerError
	}

	// Close the database connection pool after program executes
	defer db.Close()
	res:= "Usuario " + u.Username + " creado."
	return c.JSON(http.StatusCreated, res)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
 }

