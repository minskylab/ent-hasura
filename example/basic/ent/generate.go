package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
//go:generate go run github.com/minskylab/ent-hasura/cmd/ent apply -e ../.env ./schema
