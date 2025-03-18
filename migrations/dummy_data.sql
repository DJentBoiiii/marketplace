CREATE TABLE Users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255),
  email VARCHAR(255),
  password VARCHAR(255)
);

DELETE Users;

SELECT * FROM marketplace_dev.users;


CREATE TABLE Products (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100),
  price INTEGER,
  type VARCHAR(7),
  description TEXT,
  vendor INT,
  FOREIGN KEY (vendor) REFERENCES Users(id),
  product_path VARCHAR(255),
  product_img VARCHAR(255)
);

DROP TABLE Products;