# GoSQL Documentation

## Overview

GoSQL is a Go library designed to simplify database interactions, migrations, and query management.

## Packages

### Query Builder

The ConditionBuilder provide functions for dynamically constructing SQL conditions.

```go
import "github.com/mekramy/gosql/query"

func main() {
    cond := query.NewCondition(query.NumbericResolver)
    cond.And("name = ?", "John").
        AndClosure("age > ? AND age < ?", 9, 31).
        OrIf(false, "age IS NULL").
        OrClosureIf(true, "membership @in", "admin", "manager", "accountant")

    // Result: "name = $1 AND (age > $2 AND age < $3) OR (membership IN ($4, $5, $6))"
}
```

### Query Manager

The `query` package provides tools for managing and generating SQL queries.

#### Example

```go
import (
    "github.com/mekramy/gosql/query"
    "github.com/mekramy/gofs"
)
func main() {
    queriesFS := gofs.NewDir("database/queries")
    queryManager, _err_ := query.NewQueryManager(fs, query.WithRoot("database"))

    usersList :=  queryManager.Get("queries/users/users_list")
    usersTrash :=  queryManager.Get("queries/users/deleted users")
    customers, exists :=  queryManager.Find("queries/customers/list")
    customers, exists :=  queryManager.Find("queries/customers/deleted") // "", false
}
```

Query files style:

```sql
-- users.sql

-- { query: users_list }
SELECT * FROM users WHERE `deleted_at` IS NULL AND `name` LIKE ?;

-- { query: deleted users }
SELECT * FROM users WHERE deleted_at IS NOT NULL;


-- customers.sql
-- { query: list }
SELECT * from customers;
```

### Postgres Package

The `postgres` package provides tools for constructing and executing SQL commands specifically for PostgreSQL databases. Query placeholders must `?`.

```go
package main

import (
    "context"

    "github.com/mekramy/gosql/postgres"
    "github.com/jackc/pgx/v5/pgconn"
)

func main() {
    ctx := context.Background()
    config := postgres.NewConfig().
        Host("localhost").
        Port(5432).
        User("postgres").
        Password("password").
        Database("test").
        SSLMode("disable").
        MinConns(2)
    conn, err := postgres.New(
        ctx, config.Build(),
        func(c *pgxpool.Config) { c.MaxConns = 7 },
    )
    defer conn.Close(ctx)

    cmd := postgres.NewCmd(conn)
    result, err := cmd.Command("INSERT INTO users (name) VALUES (?)").Exec(ctx, "John Doe")
}
```

### MySQL Package

The `mysql` package provides tools for constructing and executing SQL commands specifically for MySQL databases.

```go
package main

import (
    "context"

    "github.com/mekramy/gosql/mysql"
)

func main() {
    conn, err := mysql.New(
        context.Background(),
        mysql.NewConfig().Database("test").Password("root").Build(),
    )
    if err != nil {
        log.Fatal(err)
    }

    cmd := mysql.NewCmd(conn)
    result, err := cmd.Command("INSERT INTO users (name) VALUES (?)").Exec(context.Background(), "John Doe")
}
```

### Migration Package

The `migration` package provides tools for managing database migrations by stage.

```go
package main

import (
    "log"

    "github.com/mekramy/gofs"
    "github.com/mekramy/gosql/migration"
    "github.com/mekramy/gosql/mysql"
)

func main() {
    conn := CreateConnection()
    fs := CreateFS()
    mig, err := migration.NewMigration(
        migration.NewMySQLSource(conn),
        fs,
        migration.WithRoot("migrations"),
    )

    err := mig.Up([]string{"table", "index", "seed"})
    if err != nil {
        log.Fatal(err)
    }
}
```

Migration files style:

```sql
-- 1741791024-create-users-table.sql
-- { up: table } table is sectin name
CREATE TABLE IF NOT EXISTS ...

-- { down: table }

-- { up: index }
...

-- { down: index }
...

-- { up: seed }
...

-- { down: seed }
...
```
