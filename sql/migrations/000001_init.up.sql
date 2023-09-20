CREATE TABLE orders (
   id          varchar(100) PRIMARY KEY,
   price       DECIMAL(10,2),
   tax         DECIMAL(10,2),
   final_price DECIMAL(10,2)
);  
