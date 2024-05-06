# EniQilo Store

Requirement: [Project Sprint EniQilo Store](https://openidea-projectsprint.notion.site/EniQilo-Store-93d69f62951c4c8aaf91e6c090127886)


## Database Migration

Database migration must use [golang-migrate](https://github.com/golang-migrate/migrate) as a tool to manage database migration

- **Short Tutorial:**
    - Direct your terminal to your project folder first
    - Initiate folder
        
        ```bash
        mkdir db/migrations
        
        ```
        
    - Create migration
        
        ```bash
        migrate create -ext sql -dir db/migrations add_user_table
        
        ```
        
        This command will create two new files named `add_user_table.up.sql` and `add_user_table.down.sql` inside the `db/migrations` folder
        
        - `.up.sql` can be filled with database queries to create / delete / change the table
        - `.down.sql` can be filled with database queries to perform a `rollback` or return to the state before the table from `.up.sql` was created
    - Execute migration
        
        ```bash
        migrate -database "postgres://postgres:password@host:5432/postgres?sslmode=disable" -path ./db/migrations -verbose up
        
        ```
        
    - Rollback migration (one migration)
        
        ```bash
        migrate -database "postgres://postgres:password@host:5432/postgres?sslmode=disable" -path db/migrations -verbose down
        
        ```
        
    - Rollback migration (all migration)
        
        ```bash
        migrate -database "postgres://postgres:password@host:5432/postgres?sslmode=disable" -path db/migrations -verbose drop
        ```


## Run & Build EniQilo Store

Run for debugging

```
make run
```

Build app

```
make build
```



