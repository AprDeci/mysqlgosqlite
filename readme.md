# mysqlgosqlite

## 依赖

[github.com/jarvanstack/mysqldump](https://github.com/jarvanstack/mysqldump)

基于 [UN1Q-com/mysql2sqlite: Converts MySQL dump to SQLite3 compatible dump](https://github.com/UN1Q-com/mysql2sqlite) 修改 

- 所有没有主键的表但存在id字段则id为主键
- 调整索引转换规则



## 所需环境

- awk
- sqlite3



## 用法

