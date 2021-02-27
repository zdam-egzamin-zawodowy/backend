#!/bin/sh
cd ./internal/graphql
go run github.com/99designs/gqlgen
go mod tidy
cd ../..