package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/sharedmemory"
)

type CCount struct {
	Country string
	Count   int64
}

type People struct {
	Persons []Person `json:"people"`
}

type Person struct {
	ID         int            `json:"id"`
	Title      sql.NullString `json:"title"`
	FirstName  string         `json:"firstname"`
	MiddleName sql.NullString `json:"middlename"`
	LastName   string         `json:"lastname"`
	Suffix     sql.NullString `json:"suffix"`
	Scode      string         `json:"scode"`
	Ccode      string         `json:"ccode"`
	State      string         `json:"state"`
	Country    string         `json:"country"`
}

// server, database, driver configuration
var server, database, driver = "lpc:localhost", "AdventureWorks2022", "mssql" // "sqlserver" or "mssql"
// protocols: https://learn.microsoft.com/en-us/sql/sql-server/connect-to-database-engine?view=sql-server-ver16&tabs=sqldb

// trusted connection, and encryption configuraiton
var trusted_connection, encrypt = true, true

// db is global variable to pass between functions
var db *sql.DB

// use background context globally to pass between functions
var ctx = context.Background()

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	mssql, err := ConnectMSSQL(ctx, db, driver, server, database, trusted_connection, encrypt)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		SelectPeople(w, r, mssql)
	})
	mux.HandleFunc("/ccount", func(w http.ResponseWriter, r *http.Request) {
		SelectCountryCount(w, r, mssql)
	})

	mux.HandleFunc("/country/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "code %s", r.PathValue("code"))
		remainder := r.URL.Path[len("/country/"):]
		countryCode, _, _ := strings.Cut(remainder, "/")
		if countryCode == "" {
			SelectPeople(w, r, mssql)
		} else {
			SelectCountry(w, r, mssql, countryCode)
		}
	})

	http.ListenAndServe(":3000", mux)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func ConnectMSSQL(
	ctx context.Context,
	db *sql.DB,
	driver string,
	server string,
	database string,
	trusted_connection bool,
	encrypt bool) (*sql.DB, error) {
	var err error

	connString := fmt.Sprintf("server=%s;database=%s;TrustServerCertificate=%v;encrypt=%v", server, database, trusted_connection, encrypt)
	db, err = sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected!\n")

	return db, nil
}

func SelectPeople(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// defer db.Close()

	// err = db.Ping(ctx)
	// if err != nil {
	// 	log.Fatalf("error after pinging PostgreSQL database: %v\n", err)
	// }

	query := `SELECT Person.Person.BusinessEntityID
,Person.Person.Title
,Person.Person.FirstName
,Person.Person.MiddleName
,Person.Person.LastName
,Person.Person.Suffix
,Person.StateProvince.StateProvinceCode
,Person.StateProvince.CountryRegionCode
,Person.StateProvince.Name
,Person.CountryRegion.Name
FROM Person.Person
JOIN Person.BusinessEntityAddress ON Person.Person.BusinessEntityID = Person.BusinessEntityAddress.BusinessEntityID
JOIN Person.Address ON Person.BusinessEntityAddress.AddressID = Person.Address.AddressID
JOIN Person.StateProvince ON Person.Address.StateProvinceID = Person.StateProvince.StateProvinceID
JOIN Person.CountryRegion ON Person.StateProvince.CountryRegionCode = Person.CountryRegion.CountryRegionCode;`
	tsql := fmt.Sprintf(query)

	// Execute query
	rows, err := db.Query(tsql)
	if err != nil {
		log.Fatal("Error reading table: " + err.Error())
	}
	defer rows.Close()

	people := []Person{}
	// Iterate through the result set.
	for rows.Next() {
		var person Person
		// Get values from row.
		err = rows.Scan(
			&person.ID,
			&person.Title,
			&person.FirstName,
			&person.MiddleName,
			&person.LastName,
			&person.Suffix,
			&person.Scode,
			&person.Ccode,
			&person.State,
			&person.Country,
		)

		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
		}

		people = append(people, person)
		// person, err = json.MarshalIndent(person, "", "\t")
		// people.AddPerson(person)
	}
	json.NewEncoder(w).Encode(people)
	// c.IndentedJSON(http.StatusOK, people)
}

// func (peeps *People) AddPerson(per Person) {
// 	peeps.Persons = append(peeps.Persons, per)
// }

func SelectCountry(w http.ResponseWriter, r *http.Request, db *sql.DB, country string) {
	// defer db.Close()

	// err = db.Ping(ctx)
	// if err != nil {
	// 	log.Fatalf("error after pinging PostgreSQL database: %v\n", err)
	// }

	query := `SELECT Person.Person.BusinessEntityID
,Person.Person.Title
,Person.Person.FirstName
,Person.Person.MiddleName
,Person.Person.LastName
,Person.Person.Suffix
,Person.StateProvince.StateProvinceCode
,Person.StateProvince.CountryRegionCode
,Person.StateProvince.Name
,Person.CountryRegion.Name
FROM Person.Person
JOIN Person.BusinessEntityAddress ON Person.Person.BusinessEntityID = Person.BusinessEntityAddress.BusinessEntityID
JOIN Person.Address ON Person.BusinessEntityAddress.AddressID = Person.Address.AddressID
JOIN Person.StateProvince ON Person.Address.StateProvinceID = Person.StateProvince.StateProvinceID
JOIN Person.CountryRegion ON Person.StateProvince.CountryRegionCode = Person.CountryRegion.CountryRegionCode
WHERE  Person.CountryRegion.CountryRegionCode = '%s';`
	tsql := fmt.Sprintf(query, country)

	// Execute query
	rows, err := db.Query(tsql)
	if err != nil {
		log.Fatal("Error reading table: " + err.Error())
	}
	defer rows.Close()

	people := []Person{}
	// Iterate through the result set.
	for rows.Next() {
		var person Person
		// Get values from row.
		err = rows.Scan(
			&person.ID,
			&person.Title,
			&person.FirstName,
			&person.MiddleName,
			&person.LastName,
			&person.Suffix,
			&person.Scode,
			&person.Ccode,
			&person.State,
			&person.Country,
		)

		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
		}

		people = append(people, person)
		// person, err = json.MarshalIndent(person, "", "\t")
		// people.AddPerson(person)
	}
	json.NewEncoder(w).Encode(people)
	// c.IndentedJSON(http.StatusOK, people)
}

func SelectCountryCount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Check if database is alive.
	// err := db.PingContext(ctx)
	// if err != nil {
	// 	log.Fatal("Error pinging database: " + err.Error())
	// }

	query := `SELECT [Person].[CountryRegion].[Name] AS "Country"
	,COUNT([Person].[Person].[BusinessEntityID]) AS "Business Sum"
	FROM [Person].[Person]
	JOIN [Person].[BusinessEntityAddress] ON [Person].[BusinessEntityAddress].[BusinessEntityID] = [Person].[Person].[BusinessEntityID]
	JOIN [Person].[Address] ON [Person].[Address].[AddressID] = [Person].[BusinessEntityAddress].[AddressID]
	JOIN [Person].[StateProvince] ON [Person].[StateProvince].[StateProvinceID] = [Person].[Address].[StateProvinceID]
	JOIN [Person].[CountryRegion] ON [Person].[CountryRegion].[CountryRegionCode] = [Person].[StateProvince].[CountryRegionCode]
	GROUP BY [Person].[CountryRegion].[Name]
	ORDER BY COUNT([Person].[Person].[BusinessEntityID]) DESC`

	tsql := fmt.Sprintf(query)

	// Execute query
	rows, err := db.Query(tsql)
	if err != nil {
		log.Fatal("Error reading table: " + err.Error())
	}

	defer rows.Close()

	// var row_count int = 0
	var ccount []CCount

	// Iterate through the result set.
	for rows.Next() {
		// var count, country string
		var cc CCount

		// Get values from row.
		// err := rows.Scan(&count, &country)
		if err := rows.Scan(&cc.Country, &cc.Count); err != nil {
			log.Fatal("Error reading rows: " + err.Error())
		}
		ccount = append(ccount, cc)

		// fmt.Printf("Country: %s, Count: %v\n", cc.Country, cc.Count)
	}

	json.NewEncoder(w).Encode(ccount)
}
