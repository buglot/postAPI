# how to run
- 1. install golang
- 2. ```go mod tidy```
- 3. create file `.env`
    ```text
    DSN_MYSQL="user/root:pass@tcp(mysql:3306)/post?charset=utf8mb4&parseTime=True&loc=Local"
    JWT_SECRAT_KEY=secertkey
    ```
- 4. run `go run .`