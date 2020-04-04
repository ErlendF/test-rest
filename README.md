# Test-Rest

This is a REST API for a very simple message board. It supports and uses a MySQL or PostgreSQL database for persistent storage (set by the --dbtype flag, defaults to postgres).

## Endpoints


### /post
*GET/POST*

- A GET requests retrieves all posts stored in the database and related comments.
- A POST request makes a new post. The body of the POST request needs to be in the following format:
```
{
    "content": "This is an example post."
}
```


### /comment
*POST*

Used for making new comments on a post. The body of the POST request needs to be in the following format:
```
{
    "post": 1,
    "content": "This is an example comment on post 1."
}
```


## Setup

Configuration may be provided by flags, environment variables or from a configuration file. They are prioritized in that order. 

Available configuration:

| Name             | Type   | Flag                  | Env              | Cfg file        |
|------------------|--------|-----------------------|------------------|-----------------|
| Config file      | String | -c, --config          | CONFIG           | -               |
| DB host          | string | --dbHost              | DB_HOST          | dbHost          |
| DB name          | string | --dbName              | DB_NAME          | dbName          |
| DB password      | string | --dbPassword          | DB_PASSWORD      | dbPassword      |
| DB port          | int    | --dbPort              | DB_PORT          | dbPort          |
| DB SSL mode      | string | --dbSslMode           | DB_SSLMODE       | dbSslMode       |
| DB type          | string | --dbType              | DB_TYPE          | dbType          |
| DB user          | string | --dbUser              | DB_USER          | dbUser          |
| Help             | bool   | -h, --help            | -                | -               |
| JSON formatter   | bool   | -j, --jsonFormatter   | JSON_FORMATTER   | jsonFormatter   |
| Port             | int    | -p, --port            | PORT             | port            |
| Shutdown timeout | int    | -s, --shutdownTimeout | SHUTDOWN_TIMEOUT | shutdownTimeout |
| Verbose          | bool   | -v, --verbose         | VERBOSE          | verbose         |

## Security

This repository is only for testing and is not meant for an production environment. As such, for example the database password WILL BE LOGGED if verbose is set. Additionaly there is currently no way to delete existing posts or comments, and no way to identify the individuals who posted them.