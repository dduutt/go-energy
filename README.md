# go-energy 能源数据采集系统

`go-energy` 是一个高性能的能源数据采集工具，采用 Go 语言编写。它支持通过 Modbus TCP 和 Siemens S7-200 Smart 协议并发采集工业设备数据，并将采集到的能源数据持久化存储到 MySQL/MariaDB 数据库中。

## 核心特性

- **多协议支持**：内置支持 Modbus TCP 和 S7-200 Smart 协议。
- **高并发采集**：基于 IP 地址进行设备分组，利用 Go 协程实现并发读取，极大提高采集效率。
- **配置驱动**：所有设备信息和数据库配置均通过 Excel 文件管理，无需修改代码。
- **数据持久化**：将采集时间、设备编号、原始字节及解析后的数值完整存储，方便后续分析。
- **自动重试**：具备完善的超时和重试机制，确保在网络波动情况下的采集成功率。

## 快速开始

### 1. 数据库初始化
在您的 MySQL/MariaDB 数据库中创建一个名为 `energy` 的数据库（或您自定义的名称），并运行项目根目录下的 `create_table.sql` 来创建数据表。

```sql
-- create_table.sql 内容参考
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
```

### 2. 准备配置文件
编辑项目根目录下的 `数据采集配置.xlsx` 文件：

- **配置表**：填写数据库连接信息（HOST, PORT, USER, PASSWORD, DATABASE）。
- **设备表**：填写采集设备明细。

### 3. 运行程序
确保您已根据目标环境编译好二进制文件（或使用已有的 `go-energy` 程序）：

```bash
./go-energy
```

## 详细配置说明 (`数据采集配置.xlsx`)

### 配置表 (Sheet: 配置表)
| 键 (Key) | 说明 |
| :--- | :--- |
| HOST | 数据库主机地址 |
| PORT | 数据库端口 |
| USER | 数据库用户名 |
| PASSWORD | 数据库密码 |
| DATABASE | 数据库名称 |

### 设备表 (Sheet: 设备表)
| 列名 | 说明 |
| :--- | :--- |
| 编号 | 设备的唯一标识码 (Code) |
| 协议 | `modbus_tcp` 或 `s7_200_smart` |
| IP/PORT | 设备的通讯地址和端口 |
| 从站/区域 | Modbus 的 SlaveId 或 S7 的变量地址 (如 V1000) |
| 地址/长度 | 采集的起始地址和读取长度 |
| 数据类型 | 支持 `float32`, `float64`, `int16`, `int32`, `uint16`, `uint32` 等 |
| 倍率 | 采集数值的缩放比例 |
| 字节序 | `大端` (默认) 或 `小端` |

## 部署建议

本项目通常作为定时任务运行。以下是一个 Linux Cron 示例（每小时执行一次）：

```bash
# 编辑 crontab
crontab -e

# 添加以下内容 (假设程序路径为 /home/user/go-energy)
0 * * * * cd /home/user/go-energy && ./go-energy >> /dev/null 2>&1
```

## 日志查看
程序运行后会在 `logs/` 目录下生成以日期命名的日志文件，例如 `logs/2026-05-07.txt`。

- `[excel error]`：Excel 文件打开或读取配置错误。
- `[db error]`：数据库连接失败。
- `[group error]`：某一组设备（同 IP）连接超时。
- `[meter error]`：具体设备读取数据失败。
- `[parse error]`：字节解析为数值时出错。
