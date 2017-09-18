This is forking version, I change on url check from HEAD into GET request.
If you need original version please check https://github.com/dimiro1/health

# Go Health Check

An easy to use, extensible health check library for Go applications.

**Table of Contents**

- [Example](#example)
- [Motivation](#motivation)
- [Inspiration](#inspiration)
- [Installation](#Installation)
- [API](#api)
- [Testing](#testing)
- [Implementing custom checkers](#implementing-custom-checkers)
- [Implemented health check indicators](#implemented-health-check-indicators)
- [LICENSE](#license)

# Example

```go
package main

import (
    "net/http"
    "database/sql"
    "time"

    "github.com/dhiemaz/health"
    "github.com/dhiemaz/health/url"
    "github.com/dhiemaz/health/db"
    "github.com/dhiemaz/health/redis"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    database, _ := sql.Open("mysql", "/")
	mysql := db.NewMySQLChecker(database)
    timeout := 5 * time.Second
    
    companies := health.NewCompositeChecker()
    companies.AddChecker("Microsoft", url.NewChecker("https://www.microsoft.com/"))
    companies.AddChecker("Oracle", url.NewChecker("https://www.oracle.com/"))
    companies.AddChecker("Google", url.NewChecker("https://www.google.com/"))

    handler := health.NewHandler()
    handler.AddChecker("Go", url.NewCheckerWithTimeout("https://golang.org/", timeout))
    handler.AddChecker("Big Companies", companies)
    handler.AddChecker("MySQL", mysql)
    handler.AddChecker("Redis", redis.NewChecker("tcp", ":6379"))

    http.Handle("/health/", handler)
    http.ListenAndServe(":8080", nil)
}
```

```sh
$ curl localhost:8080/health/
```

If everything is ok the server must respond with HTTP Status 200 OK and have following json in the body.

```json
{
    "Big Companies": {
        "Google": {
            "code": 200,
            "status": "UP"
        },
        "Microsoft": {
            "code": 200,
            "status": "UP"
        },
        "Oracle": {
            "code": 200,
            "status": "UP"
        },
        "status": "UP"
    },
    "Go": {
        "code": 200,
        "status": "UP"
    },
    "MySQL": {
        "status": "UP",
        "version": "10.1.9-MariaDB"
    },
    "Redis": {
        "status": "UP",
        "version": "3.0.5"
    },
    "status": "UP"
}
```

The server responds with HTTP Status 503 Service Unavailable if the ckeck is Down and the json response could be something like this.

```json
{
    "Big Companies": {
        "Google": {
            "code": 200,
            "status": "UP"
        },
        "Microsoft": {
            "code": 200,
            "status": "UP"
        },
        "Oracle": {
            "code": 200,
            "status": "UP"
        },
        "status": "UP"
    },
    "Go": {
        "code": 200,
        "status": "UP"
    },
    "MySQL": {
        "status": "DOWN",
        "error": "Error 1044: Access denied for user ''@'localhost' to database 'invalid-database'",
    },
    "Redis": {
        "status": "UP",
        "version": "3.0.5"
    },
    "status": "DOWN"
}
```

# Motivation

It is very important to verify the status of your system, not only the system itself, but all its dependencies, 
If your system is not Up you can easily know what is the cause of the problem only looking the health check.

Also it serves as a kind of basic integration test between the systems.

# Inspiration

I took a lot of ideas from the [spring framework](http://spring.io/).

# Installation

This package is a go getable package.

```sh
$ go get github.com/dhiemaz/health
```

# API

The API is stable and I do not have any plans to break compatibility, but I recommend you to vendor this dependency in your project, as it is a good practice.

# Testing

You have to install the test dependencies.

```sh
$ go get gopkg.in/DATA-DOG/go-sqlmock.v1
$ go get github.com/rafaeljusto/redigomock
```

or you can go get this package with the -t flag

```sh
go get -t github.com/dimiro1/health
```

# Implementing custom checkers

The key interface is `health.Checker`, you only have to implement a type that satisfies that interface.

```go
type Checker interface {
	Check() Health
}
```

Here is an example of Disk Space usage (unix only).

```go
package main

import (
    "syscall"
    "os"
)

type DiskSpaceChecker struct {
	Dir       string
	Threshold uint64
}

func NewDiskSpaceChecker(dir string, threshold uint64) DiskSpaceChecker {
	return DiskSpaceChecker{Dir: dir, Threshold: threshold}
}

func (d DiskSpaceChecker) Check() health.Health {
	health := health.NewHealth()

	var stat syscall.Statfs_t

	wd, err := os.Getwd()

	if err != nil {
        health.Down().AddInfo("error", err.Error()) // Why the check is Down
        return health
	}

	syscall.Statfs(wd, &stat)

	diskFreeInBytes := stat.Bavail * uint64(stat.Bsize)

	if diskFreeInBytes < d.Threshold {
		health.Down()
	} else {
        health.Up()
    }

    health.
        AddInfo("free", diskFreeInBytes).
        AddInfo("threshold", d.Threshold)

	return health
}
```

## Important

The **status** key in the json has priority over a **status** key added by a Checker, so if some checker adds a **status** key to the json, it will not be rendered  

# Implemented health check indicators

| Health         | Description                            | Package                                              |
|----------------|----------------------------------------|------------------------------------------------------|
| url.Checker    | Check the connection with some URL     | https://github.com/dhiemaz/health/tree/master/url    |
| db.Checker     | Check the connection with the database | https://github.com/dhiemaz/health/tree/master/db     |
| redis.Checker  | Check the connection with the redis    | https://github.com/dhiemaz/health/tree/master/redis  |

# LICENSE

The MIT License (MIT)

Copyright (c) 2016 Claudemiro

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
