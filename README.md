# WB parsing service

### Initialize SQL structure

Connect to running instance of MS SQL Server
```bash
docker run -it --network ${MSSQL_NETWORK} emergn/mssql-tools -S mssql -U sa -P ${MSSQL_PASSWORD} -C
```

Inside mssql-tools container:
```sql
CREATE DATABASE DWH_test
GO
USE DWH_test;
GO
CREATE TABLE wbProduct( id INT PRIMARY KEY IDENTITY (1,1), modifiedDate DATETIME, nmID int, name NVARCHAR(200), price DECIMAL(16,2),);
GO
CREATE TABLE wbSellerPrice (id INT PRIMARY KEY IDENTITY (1,1), modifiedDate DATETIME, nmID INT NOT NULL REFERENCES wbProduct (id), seller NVARCHAR(100), price DECIMAL(16,2));
GO
```