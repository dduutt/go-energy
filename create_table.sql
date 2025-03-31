CREATE TABLE energy (
    id INT AUTO_INCREMENT PRIMARY KEY,-- 主键，自增
    CODE INT NOT NULL,-- 对应 Code 字段
    workshop VARCHAR (255),-- 对应 WorkShop 字段
    room VARCHAR (255),-- 对应 Room 字段
    NAME VARCHAR (255),-- 对应 Name 字段
    VALUE DECIMAL (10,2),-- 对应 Value 字段
    bytes VARCHAR (255),-- 对应 Bytes 字段
    protocol VARCHAR (255),-- 对应 Protocol 字段
    ip VARCHAR (255),-- 对应 IP 字段
    PORT INT,-- 对应 Port 字段
    slave_or_area VARCHAR (255),-- 对应 SlaveOrArea 字段
    START INT,-- 对应 Start 字段
    size INT,-- 对应 Size 字段
    data_type VARCHAR (255),-- 对应 DataType 字段
    byte_order VARCHAR (255),-- 对应 IsLittleEndian 字段
    magnification FLOAT-- 对应 Magnifications 字段
);