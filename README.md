# RESTful API using go lang

This repo is using the [go](https://go.dev/) programming lanuage and the [Adventure Works](https://learn.microsoft.com/en-us/sql/samples/adventureworks-install-configure?view=sql-server-ver16&tabs=ssms) sample database to create a RESTful api.

# Prerequisites

This project assumes one is using Windows OS and has the following already installed:
1. [Microsoft SQL Server](https://www.microsoft.com/en-ca/sql-server/sql-server-downloads)
2. [SSMS](https://learn.microsoft.com/en-us/sql/ssms/download-sql-server-management-studio-ssms?view=sql-server-ver16) for restoring Adventure Works database
3. [Adventure Works](https://learn.microsoft.com/en-us/sql/samples/adventureworks-install-configure?view=sql-server-ver16&tabs=ssms) 2022 OLTP sample database
   - Download 'AdventureWorks2022.bak'
   - Restore database using SSMS
4. [Go](https://go.dev/dl/) programming language
5. Installing [pure Go](https://github.com/microsoft/go-mssqldb) database driver for Go's 'database/sql' package
	- repo above is forked from [denisenkom](https://github.com/denisenkom/go-mssqldb)

# Run App

clone this repo using:
> git clone https://github.com/davidhintelmann/AdventureWorksRESTful.git

now run app using:
> go run main.go

or build the app and run it:
> go build -o api.exe main.go
> 
> .\api.exe

You can now open a browser and navigate to 'localhost:3000'. Endpints are at 'localhost:3000/{endpoint}'

# Endpoints

After cloning this repo and running the app, you can access endpints at 'localhost:3000/{endpoint}'

Endpoints:
-   /people
-   /country/{COUNTRY CODE}
-   /ccount

## people

This endpoint will return information about people's names, country, etc.

The '/people' endpoint, found at 'localhost:3000/people' will return the following SQL Query:

```sql
SELECT Person.Person.BusinessEntityID
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
JOIN Person.CountryRegion ON Person.StateProvince.CountryRegionCode = Person.CountryRegion.CountryRegionCode;
```

## ccount

This endpoint will return information about how many businesses are in each country.

The '/ccount' endpoint, found at 'localhost:3000/ccount' will return the following SQL Query:

```sql
SELECT [Person].[CountryRegion].[Name] AS "Country"
,COUNT([Person].[Person].[BusinessEntityID]) AS "Business Sum"
FROM [Person].[Person]
JOIN [Person].[BusinessEntityAddress] ON [Person].[BusinessEntityAddress].[BusinessEntityID] = [Person].[Person].[BusinessEntityID]
JOIN [Person].[Address] ON [Person].[Address].[AddressID] = [Person].[BusinessEntityAddress].[AddressID]
JOIN [Person].[StateProvince] ON [Person].[StateProvince].[StateProvinceID] = [Person].[Address].[StateProvinceID]
JOIN [Person].[CountryRegion] ON [Person].[CountryRegion].[CountryRegionCode] = [Person].[StateProvince].[CountryRegionCode]
GROUP BY [Person].[CountryRegion].[Name]
ORDER BY COUNT([Person].[Person].[BusinessEntityID]) DESC
```

## country

This endpoint will return similar information as the /people endpoint but it will filter based on country code.

The '/country' endpoint, found at 'localhost:3000/country/{COUNTRY CODE}' will return the following SQL Query:

```sql
SELECT Person.Person.BusinessEntityID
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
WHERE  Person.CountryRegion.CountryRegionCode = '{COUNTRY CODE}';
```

**Note:** {COUNTRY CODE} in the above sql query is a [two letter country code](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Officially_assigned_code_elements)

Navigating to 'localhost:3000/country/CA' will only return people from Canada.