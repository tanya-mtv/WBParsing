CREATE TABLE IF NOT EXISTS wbProduct
(
    id INT PRIMARY KEY IDENTITY (1,1),
    modifiedDate DATETIME,
    nmID int,
    name NVARCHAR(200),
    price DECIMAL(16,2)
);

CREATE TABLE IF NOT EXISTS wbSellerPrice
(
    id INT PRIMARY KEY IDENTITY (1,1),
    modifiedDate DATETIME,
    nmID INT NOT NULL REFERENCES wbProduct (id),
    seller NVARCHAR(100),
    price DECIMAL(16,2)
);