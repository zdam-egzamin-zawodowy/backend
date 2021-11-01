# zdamegzaminzawodowy.pl - backend

This project contains the API and other core infrastructure items needed for all [zdamegzaminzawodowy.pl](https://zdamegzaminzawodowy.pl) apps.

## Development

### Prerequisites
1. Golang
2. Docker + docker-compose
3. make

### Installation

1. ``git clone git@github.com:zdam-egzamin-zawodowy/backend.git``.
2. ``cd backend && cp .env.example .env.local``
3. ``make docker-compose-up``
4. Run the app - ``go run ./cmd/server/main.go``.

## License
Distributed under the MIT License. See ``LICENSE`` for more information.

## Contact
Dawid Wysokiński - [contact@dwysokinski.me](mailto:contact@dwysokinski.me)
