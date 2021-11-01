# zdamegzaminzawodowy.pl - backend

This project contains the API and other core infrastructure items needed for all [zdamegzaminzawodowy.pl](https://zdamegzaminzawodowy.pl) apps.

## Development

### Prerequisites
1. Golang
2. PostgreSQL database

### Installation
**Required ENV variables (you can set them directly in your system or create the .env.local file):**
```
DB_USER=db_user
DB_NAME=db_name
DB_PORT=db_port
DB_HOST=db_host
DB_PASSWORD=db_pass
DB_POOL_SIZE=40
LOG_DB_QUERIES=true
ACCESS_SECRET=access_token_secret
FILE_STORAGE_PATH=path_to_the_folder_where_uploaded_files_will_be_stored
ENABLE_ACCESS_LOG=false
```

1. Clone this repo - ``git clone git@github.com:zdam-egzamin-zawodowy/backend.git``.
2. Set the required env variables.
3. Run the app - ``go run ./cmd/server/main.go``.

## License
Distributed under the MIT License. See ``LICENSE`` for more information.

## Contact
Dawid Wysokiński - [contact@dwysokinski.me](mailto:contact@dwysokinski.me)
