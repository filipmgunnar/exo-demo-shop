package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Amounts are integer cents everywhere; float dollars drifted by a cent or two
// on the weekly payout reconciliation.
type LineItem struct {
	SKU            string
	Name           string
	Quantity       int
	UnitPriceCents int64
	DiscountPct    int
}

type Order struct {
	ID       int
	Customer string
	Placed   time.Time
	Items    []LineItem
}

func (li LineItem) TotalCents() int64 {
	gross := li.UnitPriceCents * int64(li.Quantity)
	if li.DiscountPct == 0 {
		return gross
	}
	return gross * int64(100-li.DiscountPct) / 100
}

func (o Order) TotalCents() int64 {
	var total int64
	for _, li := range o.Items {
		total += li.TotalCents()
	}
	return total
}

func (o Order) ItemCount() int {
	n := 0
	for _, li := range o.Items {
		n += li.Quantity
	}
	return n
}

func FormatCents(cents int64) string {
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}

type Store struct {
	mu     sync.RWMutex
	orders map[int]Order
}

func NewStore() *Store {
	return &Store{orders: make(map[int]Order)}
}

func (s *Store) Put(o Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[o.ID] = o
}

func (s *Store) Get(id int) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orders[id]
	return o, ok
}

func (s *Store) All() []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Order, 0, len(s.orders))
	for _, o := range s.orders {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func SeedStore() *Store {
	s := NewStore()
	for _, o := range seedOrders() {
		s.Put(o)
	}
	return s
}

func seedOrders() []Order {
	return []Order{
		{
			ID:       1041,
			Customer: "Priya Raman",
			Placed:   time.Date(2026, 7, 21, 9, 14, 0, 0, time.UTC),
			Items: []LineItem{
				{SKU: "CB-500", Name: "Cold brew concentrate 500ml", Quantity: 2, UnitPriceCents: 1250},
				{SKU: "FP-100", Name: "Filter papers, 100 ct", Quantity: 1, UnitPriceCents: 600},
			},
		},
		{
			ID:       1042,
			Customer: "Marcus Boyle",
			Placed:   time.Date(2026, 7, 21, 15, 2, 0, 0, time.UTC),
			Items: []LineItem{
				{SKU: "ESP-1K", Name: "Espresso beans 1kg", Quantity: 4, UnitPriceCents: 2400, DiscountPct: 15},
				{SKU: "PTC-06", Name: "Milk frothing pitcher 600ml", Quantity: 1, UnitPriceCents: 1850},
			},
		},
		{
			ID:       1043,
			Customer: "Elena Fischer",
			Placed:   time.Date(2026, 7, 22, 11, 40, 0, 0, time.UTC),
			Items: []LineItem{
				{SKU: "KTL-09", Name: "Pour-over kettle 0.9l", Quantity: 1, UnitPriceCents: 7900, DiscountPct: 10},
				{SKU: "SCL-02", Name: "Brew scale", Quantity: 1, UnitPriceCents: 4500},
			},
		},
		{
			ID:       1044,
			Customer: "Tomas Ferreira",
			Placed:   time.Date(2026, 7, 23, 8, 5, 0, 0, time.UTC),
			Items: []LineItem{
				{SKU: "MUG-CR", Name: "Ceramic mug 300ml", Quantity: 6, UnitPriceCents: 900},
			},
		},
		{
			ID:       1045,
			Customer: "Dana Whitfield",
			Placed:   time.Date(2026, 7, 24, 17, 31, 0, 0, time.UTC),
			Items: []LineItem{
				{SKU: "DEC-250", Name: "Decaf beans 250g", Quantity: 3, UnitPriceCents: 1100, DiscountPct: 20},
				{SKU: "CLN-30", Name: "Cleaning tablets, 30 ct", Quantity: 2, UnitPriceCents: 1400},
			},
		},
	}
}
