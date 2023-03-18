
# Golang Api with Docker for Digitalizing Historical OHLC Price Data

## Background

 A sample of data below is centralized and digitalized from a csv file.

**Example of the CSV data** \

| UNIX          | SYMBOL  | OPEN           | HIGH           | LOW            | CLOSE          |
|     ---       |   ---   |    ---         |     ---        |     ---        |       ---      |
| 1644719700000 | BTCUSDT | 42123.29000000 | 42148.32000000 | 42120.82000000 | 42146.06000000 |
| 1644719640000 | BTCUSDT | 42113.08000000 | 42126.32000000 | 42113.07000000 | 42123.30000000 |
| 1644719580000 | BTCUSDT | 42120.80000000 | 42130.23000000 | 42111.01000000 | 42113.07000000 |
| 1644719520000 | BTCUSDT | 42114.47000000 | 42123.31000000 | 42102.22000000 | 42120.80000000 |
| 1644719460000 | BTCUSDT | 42148.23000000 | 42148.24000000 | 42114.04000000 | 42114.48000000 |

## On this project the major external libraries used are

- Go gin http framework:  [Go gin documentation](https://gin-gonic.com/docs/)
- Go gorm:  [ORM library for Golang](https://gorm.io/docs/)
- Postgress database: [With pgx as its driver](https://github.com/jackc/pgx)

## Endpoints

1. **POST /data**
  This is a `<multipart/form-data>` content type http post request having [csv_file  => *.csv] as key value payload with value as csv file *containing the data sample mentioned above*

  Example: POST [http://127.0.0.1:8090/data](http://127.0.0.1:8090/data)

  **Successful response payload**
    `
      {
        "data": {
            "csvLinesRead": "200",
            "totalSavedRows": "200",
        },
        "status": "success"
      }
    `

- csvLinesRead: Total number of rows on the csv file (excluding the head).
- totalSavedRows: "Total number of rows successfully saved in the database".

2. **GET /data**
  This is a get request to query the OHLC saved data.

  Example: GET [http://127.0.0.1:8090/data](http://127.0.0.1:8090/data)

  **Url Query**

- search: Value to use to perform databse full search
- limit: Value of number of items to request per request
- page: Value of current page, default is 1.
- ptype: Value to determine the type of pagination object returned with the response data

- *Request with search query* [http://127.0.0.1:8090/data?search=1644719700000](http://127.0.0.1:8090/data?search=1644719700000)

- *Request with search, limit and ptype queries* [http://127.0.0.1:8090/data?search=1644719700000&limit=100&ptype=full](http://127.0.0.1:8090/data?search=1644719700000&limit=100&ptype=full)

- *Request with page and limit queries* [http://127.0.0.1:8090/data?page=2&limit=1000](http://127.0.0.1:8090/data?page=2&limit=1000)

### Response Examples

  1. [http://127.0.0.1:8090/data?search=1644719700000&limit=100&ptype=full](http://127.0.0.1:8090/data?search=1644719700000&limit=3&ptype=full)
  **Successful response payload**

    {
      "data": [
          {
              "unix": 1644719460000,
              "symbol": "BTCUSDT",
              "open": 1644719500000,
              "high": 1644719500000,
              "low": 1644719500000,
              "close": 1644719500000
          },
          {
              "unix": 1644719460000,
              "symbol": "BTCUSDT",
              "open": 1644719500000,
              "high": 1644719500000,
              "low": 1644719500000,
              "close": 1644719500000
          },
          {
              "unix": 1644719460000,
              "symbol": "BTCUSDT",
              "open": 1644719500000,
              "high": 1644719500000,
              "low": 1644719500000,
              "close": 1644719500000
          }
      ],
      "message": "Data successfully fetched",
      "pagination": {
        "current_page_url": "http://127.0.0.1:8090/data?search=1644719460000&ptype=full&ptype=full&page=1&limit=3",
        "current_page": 1,
        "total_pages": 54,
        "per_page": 3,
        "limit": 3,
        "previous_page": 0,
        "next_page": 2,
        "current_page_total": 3,
        "previous_page_url": "",
        "next_page_url": "http://127.0.0.1:8090/data?search=1644719460000&ptype=full&ptype=full&page=2&limit=3",
        "last_page_url": "http://127.0.0.1:8090/data?search=1644719460000&ptype=full&ptype=full&page=54&limit=3",
        "total": 160
      },
      "status": "success"
    }

  2. [http://127.0.0.1:8090/data?page=2&limit=1000](http://127.0.0.1:8090/data?page=2&limit=1000)
    **Successful response payload**
    `
    {
      "data": [
        {
          "unix": 1644719460000,
          "symbol": "BTCUSDT",
          "open": 1644719500000,
          "high": 1644719500000,
          "low": 1644719500000,
          "close": 1644719500000
        },
        {
          "unix": 1644719460000,
          "symbol": "BTCUSDT",
          "open": 1644719500000,
          "high": 1644719500000,
          "low": 1644719500000,
          "close": 1644719500000
        },
        {
          "unix": 1644719460000,
          "symbol": "BTCUSDT",
          "open": 1644719500000,
          "high": 1644719500000,
          "low": 1644719500000,
          "close": 1644719500000
        }
      ],
      "message": "Data successfully fetched",
      "status": "success"
    }

## App Information

This app is created and testes on linux with docker. To run this app on windows
please visit [this link](https://www.thewindowsclub.com/how-to-run-sh-or-shell-script-file-in-windows-10) for more instructions.

**Folder structure**

production\
├── docker-compose.yml\
├── .env\
├── project\
│   ├── controller\
│   │   ├── create.go\
│   │   └── get.go\
│   ├── Dockerfile\
│   ├── .env\
│   ├── .gitignore\
│   ├── go.mod\
│   ├── go.sum\
│   ├── main.go\
│   ├── middleware\
│   │   ├── cors.go\
│   │   └── timeout.go\
│   ├── model\
│   │   ├── db.go\
│   │   └── ohlc.go\
│   ├── router\
│   │   └── router.go\
│   ├── services\
│   │   ├── error_response.go\
│   │   └── paginator.go\
│   └── test\
│       ├── create_test.go\
│       └── get_test.go\
├── psql_data\
├── readme.md\
└── run.sh

- project: Contains the golang files and application logic
- psql_data: Is the database volume folder
- .env: Is the environment variables for postgres database setup
- .prod.env: Is the production version environment variables for postgres database setup
- docker-compose.local.yml: Docker compose file for local version
- docker-compose.yml: Docker compose file for production version
- run.sh: Bash script to build the app container and start the app.
- readme.md: Application's documentation file

### Project Folder content (files)

- Controller folder where the logic of the app is located
- Model folder where the gorm database instance and Ohlc model are located
- Sercives contains helper functions
- Middleware contains middleware function for cors and
- Router contains app endpoints and instanciated in the main.go
- Test contains test files
- main.go is the enttry file of the app and its http server
- Dockfile is the docker file for the app image

### Logic Note

  Considering the fact that a very large data would be processed, speed, efficiency, memory consumption are all considered.
  A worker pool is orchestrated for saving the csv data in chunks and bacthes using database transaction for a rollback if an error is encounter in any of the spawned goroutine pool.
  The number of worker depends on the number of system's CPU and file size.

**Database mode**\
  Database transanction database used to enable rollback if there's any error during insertion.
  
**Test Sample**

- Transaction mode is a it slower, tested with a csv file of 296MB containing 4,000,001 lines of data OHLC data sample, including the header. Processed withing the average time of 1m 26s.

- Normal is faster, tested with a csv file of 403MB containing of 5,455,001 lines of data OHLC data sample, including the header. Processed withing the average time of 1m 20s.

## Run the app

cd into the app root folder from your terminal an run the command `bash run.sh`.
If you are using windows OS, please visit [this link](https://www.thewindowsclub.com/how-to-run-sh-or-shell-script-file-in-windows-10) for more instructions on how to run  shell script on windows.

## Unit Test

  Go gin does not provfide enough information for testing on gin engine *ServeHTTP* approach especially for modularised project.
  
  Native golang hhtp client request is utilzed to test the endpoints. The test files are located in the test folder or module in side the project folder.

  **Note:** \
  Before testing, the app or the container must be running actively because requests are sent to the local url and app exposed port.
  If the app port is changed in the docker-compose configuration, const port 8090 used needs to be changed in *create_test.go* file inside test.

  No testing database is setup for testing, the same database in the app is used for pupolating the database. A temporary csv file is generated in the *create_test.go.TestSaveSCV*, the number of rows is 50,000 at a time, subsequent running of test successfully would keep poplating the database.

  **Test command** \
  cd into project folder and and run ` go run -v ./test ` to run the test.

### IMPORTANT NOTICE

  If the root directory is this project is *production* or *prod_bundle.sh* is not foot at the root, it means you have the production bundle; please kindly ignore this section as the documentation strictly addressed the production bundle already.

  This project contains the local setup and production setup. Everything discussed in the documentation up is for production bundle after running ` bash prod_bundle.sh ` to generate the production folder and the setup. Run ` bash local.sh ` for the development setup as it allows hot reload with the container volume. The container would be built and started.
  