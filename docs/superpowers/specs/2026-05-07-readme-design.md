# README.md Design Specification - go-energy

## 1. Overview
`go-energy` is a high-concurrency data collection tool written in Go, designed to gather energy consumption data from various industrial devices and store it in a centralized MySQL/MariaDB database.

## 2. Structure
The `README.md` will follow a "User Guide" style, optimized for operators and system administrators.

### Sections:
- **Project Introduction**: Core value proposition (multi-protocol, high concurrency, Excel-driven).
- **Quick Start**:
    - Step 1: Initialize Database (mention `create_table.sql`).
    - Step 2: Configure Excel (mention `数据采集配置.xlsx`).
    - Step 3: Run the binary.
- **Configuration Guide (Deep Dive)**:
    - **Database Config**: Fields in "配置表" (HOST, PORT, USER, PASSWORD, DATABASE).
    - **Meter Config**: Detailed explanation of "设备表" columns:
        - Protocol: `modbus_tcp`, `s7_200_smart`.
        - Data Type: `float32`, `int16`, etc.
        - Byte Order: `大端`, `小端`.
        - Magnification: Handling scaling factors.
- **Technical Features**:
    - Supported Protocols.
    - Supported Data Types.
    - Concurrency Model (per-IP grouping).
    - Retry Logic (3 retries by default).
- **Deployment**:
    - Linux Cron setup example for hourly execution.
- **Maintenance & Troubleshooting**:
    - Log location (`logs/YYYY-MM-DD.txt`).
    - Error code meanings.

## 3. Formatting
- Use GitHub-flavored Markdown.
- Use code blocks for SQL and Bash commands.
- Use tables for Excel field descriptions.

## 4. Language
- Chinese (Simplified), as per project context and user preference.
