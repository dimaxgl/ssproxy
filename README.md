### SSProxy (Simple Socks Proxy)

This proxy based on **github.com/armon/go-socks5** package with some features:
- customizable user credentials store (_database and memory are implemented_)
- pluggable password security (TODO, only bcrypt is now supported)

###### Installation:
```bash
go get -u github.com/dimaxgl/ssproxy/proxy
```

Please take a look for _config.sample.yaml_, which helps you to understand proxy settings:

```yaml
# tcp listen address
listenAddress: 0.0.0.0:1111
# store settings
store:
  # store type
  type: database
  # store-specific parameters
  params:
    # sql driver
    driverName: postgres
    # sql connection string
    dsn: postgres://admin:admin@localhost/proxy?sslmode=disable
    # specific user column
    userColumn: username
    # specific queries
    queries:
      # specific add user sql query
      addUserExecQuery: "INSERT INTO users (\"username\", \"password\") VALUES (:user,:password)"
      # specific search user sql query
      getUserQuery: "SELECT password FROM users WHERE \"username\" = :user"
```