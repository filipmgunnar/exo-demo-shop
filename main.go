package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

const ordersPage = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>exo-demo-shop admin</title>
<style>
  body { font: 15px/1.5 -apple-system, Segoe UI, Roboto, sans-serif; margin: 2.5rem auto; max-width: 62rem; padding: 0 1rem; color: #1c1f23; }
  h1 { font-size: 1.4rem; margin-bottom: 0.25rem; }
  p.sub { color: #6b7280; margin-top: 0; }
  table { border-collapse: collapse; width: 100%; margin-top: 1.5rem; }
  th, td { text-align: left; padding: 0.6rem 0.75rem; border-bottom: 1px solid #e5e7eb; vertical-align: top; }
  th { font-size: 0.78rem; text-transform: uppercase; letter-spacing: 0.04em; color: #6b7280; }
  td.total { text-align: right; font-variant-numeric: tabular-nums; white-space: nowrap; }
  ul { margin: 0; padding-left: 1.1rem; }
  li { margin-bottom: 0.2rem; }
  .when { color: #6b7280; font-size: 0.82rem; }
  .off { color: #b45309; }
</style>
</head>
<body>
<h1>Orders</h1>
<p class="sub">{{len .Orders}} orders</p>
<table>
  <thead>
    <tr><th>Order</th><th>Customer</th><th>Line items</th><th style="text-align:right">Total</th></tr>
  </thead>
  <tbody>
  {{- range .Orders}}
    <tr>
      <td>#{{.ID}}<br><span class="when">{{.Placed.Format "Jan 2, 15:04"}}</span></td>
      <td>{{.Customer}}</td>
      <td>
        <ul>
        {{- range .Items}}
          <li>{{.Quantity}} &times; {{.Name}} @ {{money .UnitPriceCents}}{{if .DiscountPct}} <span class="off">({{.DiscountPct}}% off)</span>{{end}}</li>
        {{- end}}
        </ul>
      </td>
      <td class="total">{{money .TotalCents}}</td>
    </tr>
  {{- end}}
  </tbody>
</table>
</body>
</html>
`

// Parsed at start-up so a broken template kills the process instead of the
// first request that happens to hit it.
var ordersTmpl = template.Must(template.New("orders").
	Funcs(template.FuncMap{"money": FormatCents}).
	Parse(ordersPage))

type server struct {
	store *Store
}

func (s *server) handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := struct {
		Orders []Order
	}{Orders: s.store.All()}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ordersTmpl.Execute(w, data); err != nil {
		log.Printf("render orders: %v", err)
	}
}

func main() {
	srv := &server{store: SeedStore()}

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleOrders)

	httpSrv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
