CREATE TABLE energy (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code INT NOT NULL,
    workshop VARCHAR (255),
    room VARCHAR (255),
    name VARCHAR (255),
    value DECIMAL (10,2),
    datetime DATETIME,
    bytes VARCHAR (255),
    protocol VARCHAR (255),
    ip VARCHAR (255),
    PORT INT,
    slave_or_area VARCHAR (255),
    start INT,
    size INT,
    data_type VARCHAR (255),
    byte_order VARCHAR (255),
    magnification FLOAT
);