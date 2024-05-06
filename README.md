# Cats Social App

Requirement: [Project Sprint Cats Social](https://openidea-projectsprint.notion.site/Cats-Social-9e7639a6a68748c38c67f81d9ab3c769)

Make new migration script

```migrate create -ext sql -dir migration -seq init```

Migrate database

up : ```migrate -database "postgres://postgres:password@localhost:5432/catsx?sslmode=disable" -path db/migrations up```

down : ```migrate -database "postgres://postgres:password@localhost:5432/catsx?sslmode=disable" -path db/migrations down```


## Run & Build Cats App

Run for debugging

```
make run
```

Build app

```
make build
```



