package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Replace with your own connection parameters
var server = "localhost"
var port = 1433
var user = ""
var password = ""
var database = "Northwind"

var db *sql.DB

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Jafari"
	claims["type"] = "Super Admin"
	claims["app"] = "Aplicacion 1"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Esta ruta es pública para todos")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	appId := claims["app"].(string)
	compare := strings.Compare(appId, "Aplicacion 1")
	if compare == 0 {
		name := claims["name"].(string)
		return c.String(http.StatusOK, "Welcome "+name+"! You are in "+ appId +".")
	} else {
		return echo.ErrUnauthorized
	}
}

func connectionSql () {
	var err error

    // Create connection string
    connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
        server, user, password, port, database)

    // Create connection pool
    db, err = sql.Open("sqlserver", connString)
    if err != nil {
        log.Fatal("Error creating connection pool: " + err.Error())
    }
    log.Printf("Connected!\n")

	// Read employees
    count, err := ReadEmployees()
    if err != nil {
        log.Fatal("Error reading Employees: ", err.Error())
    }
    log.Printf("Read %d row(s) successfully.\n", count)

    // Close the database connection pool after program executes
    defer db.Close()

	selectVersion()
}

// ReadEmployees reads all employee records
func ReadEmployees() (int, error) {
    ctx := context.Background()

    // Check if database is alive.
    err := db.PingContext(ctx)
    if err != nil {
        return -1, err
    }

    tsql := fmt.Sprintf("SELECT EmployeeID, FirstName, Title FROM Employees;")

    // Execute query
    rows, err := db.QueryContext(ctx, tsql)
    if err != nil {
        return -1, err
    }

    defer rows.Close()

    var count int

    // Iterate through the result set.
    for rows.Next() {
        var FirstName, Title string
        var EmployeeId int

        // Get values from row.
        err := rows.Scan(&EmployeeId, &FirstName, &Title)
        if err != nil {
            return -1, err
        }

        fmt.Printf("ID: %d, Name: %s, Title: %s\n", EmployeeId, FirstName, Title)
        count++
    }

    return count, nil
}

// Gets and prints SQL Server version
func selectVersion(){
    // Use background context
    ctx := context.Background()

    // Ping database to see if it's still alive.
    // Important for handling network issues and long queries.
    err := db.PingContext(ctx)
    if err != nil {
        log.Fatal("Error pinging database: " + err.Error())
    }

    var result string

    // Run query and scan for result
    err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
    if err != nil {
        log.Fatal("Scan failed:", err.Error())
    }
    fmt.Printf("%s\n", result)
}

func main() {
	e := echo.New()
	log.Printf("CONECTANDO...\n")
	connectionSql ()
	log.Printf("TERMINADO DE CONECTAR\n")
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":1323"))
}
