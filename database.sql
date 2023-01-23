CREATE DATABASE GINMSSQL

GO 

USE GINMSSQL

GO

CREATE TABLE Persons (
    id int IDENTITY(1,1) PRIMARY KEY,
    username varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    age int NOT NULL
)

INSERT INTO Persons(username, email, age) VALUES(N'Tran Dung', N'trandungksnb00@gmail.com', 21)

SELECT * FROM Persons