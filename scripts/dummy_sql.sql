USE marketplace_test;

CREATE TABLE Users (
    id            INT PRIMARY KEY AUTO_INCREMENT,
    username      VARCHAR(50) UNIQUE NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin      BOOLEAN DEFAULT FALSE,
    profile_photo VARCHAR(255) DEFAULT NULL,
    bio           TEXT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Products (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(100) NOT NULL,
    price       INT NOT NULL,
    type        VARCHAR(50) NOT NULL,
    description TEXT DEFAULT NULL,
    vendor       VARCHAR(50) NOT NULL,
    image_url   VARCHAR(255) DEFAULT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_vendor FOREIGN KEY (vendor) REFERENCES Users(username) ON DELETE CASCADE
);


ALTER TABLE Products
ADD COLUMN Genre VARCHAR(100);

ALTER TABLE `Products`
MODIFY column  `Genre` VARCHAR(100) DEFAULT "none";

DROP table `Users`;
DROP table `Products`;


CREATE TABLE Cart (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    user_id     INT NOT NULL,
    product_id  INT NOT NULL,
    added_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cart_user FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    CONSTRAINT fk_cart_product FOREIGN KEY (product_id) REFERENCES Products(id) ON DELETE CASCADE
);


DROP table `Cart`;


CREATE TABLE Playlists (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(100) NOT NULL,
    user_id     INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_playlist_user FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE TABLE PlaylistItems (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    playlist_id INT NOT NULL,
    product_id  INT NOT NULL,
    added_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_playlistitem_playlist FOREIGN KEY (playlist_id) REFERENCES Playlists(id) ON DELETE CASCADE,
    CONSTRAINT fk_playlistitem_product FOREIGN KEY (product_id) REFERENCES Products(id) ON DELETE CASCADE
);

CREATE TABLE Purchases (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    user_id     INT NOT NULL,
    product_id  INT NOT NULL,
    purchased_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_purchase (user_id, product_id),
    CONSTRAINT fk_purchase_user FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    CONSTRAINT fk_purchase_product FOREIGN KEY (product_id) REFERENCES Products(id) ON DELETE CASCADE
);



CREATE TABLE Comments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    comment TEXT NOT NULL,
    likes_product BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_comment_user FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_product FOREIGN KEY (product_id) REFERENCES Products(id) ON DELETE CASCADE
);