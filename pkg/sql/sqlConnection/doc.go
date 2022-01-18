// Package sqlConnection define method to create connection to SQL database.
//  Open return a pointer on database using GORM using environment variables :
//  - DB_SQL_DRIVER : SQL driver used by database (postgresql, mysql, etc...) (default : postgres)
//  - DB_HOST : hostname database (default : localhost)
//  - DB_PORT : port database (default : 5432)
//  - DB_USER : user (default : postgres)
//  - DB_NAME : database name (default : postgres)
//  - DB_PASSWORD : password (default : postgres)
//  - DB_EXTRA_CONFIGS : extra configs connection(default : sslmode=disable TimeZone=Pacific/Noumea)
//  - DB_DEBUG : enable or disable debug log (default : true)
//  - DB_LOG : enable or disable query database log (default : true)
package sqlConnection
