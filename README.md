# exo-demo-shop

Tiny admin panel for the shop's order desk. It serves a single HTML page listing
every order with its customer, line items and total. Orders live in memory and
are seeded at start-up, so restarting resets everything.

Go standard library only, no dependencies.

## Run

```
go run .
```

Then open http://localhost:8080.

## Test

```
go test ./...
```

## Layout

- `main.go` — HTTP server and the orders page template
- `store.go` — order types, total calculation, in-memory store and seed data
