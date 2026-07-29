package main

import "testing"

func TestLineItemTotalNoDiscount(t *testing.T) {
	li := LineItem{SKU: "MUG-CR", Name: "Ceramic mug", Quantity: 3, UnitPriceCents: 900}
	if got, want := li.TotalCents(), int64(2700); got != want {
		t.Fatalf("TotalCents() = %d, want %d", got, want)
	}
}

func TestLineItemTotalSingleUnitDiscount(t *testing.T) {
	li := LineItem{SKU: "KTL-09", Name: "Kettle", Quantity: 1, UnitPriceCents: 7900, DiscountPct: 10}
	if got, want := li.TotalCents(), int64(7110); got != want {
		t.Fatalf("TotalCents() = %d, want %d", got, want)
	}
}

func TestLineItemTotalMultiUnitDiscount(t *testing.T) {
	li := LineItem{SKU: "ESP-1K", Name: "Espresso beans 1kg", Quantity: 4, UnitPriceCents: 2400, DiscountPct: 15}
	if got, want := li.TotalCents(), int64(8160); got != want {
		t.Fatalf("TotalCents() = %d, want %d", got, want)
	}
}

func TestLineItemTotalFullDiscount(t *testing.T) {
	li := LineItem{SKU: "PROMO", Name: "Sample bag", Quantity: 1, UnitPriceCents: 500, DiscountPct: 100}
	if got, want := li.TotalCents(), int64(0); got != want {
		t.Fatalf("TotalCents() = %d, want %d", got, want)
	}
}

func TestOrderTotalSumsLineItems(t *testing.T) {
	o := Order{
		ID:       9001,
		Customer: "Test Buyer",
		Items: []LineItem{
			{SKU: "CB-500", Quantity: 2, UnitPriceCents: 1250},
			{SKU: "FP-100", Quantity: 1, UnitPriceCents: 600},
		},
	}
	if got, want := o.TotalCents(), int64(3100); got != want {
		t.Fatalf("TotalCents() = %d, want %d", got, want)
	}
}

func TestOrderTotalEmpty(t *testing.T) {
	var o Order
	if got := o.TotalCents(); got != 0 {
		t.Fatalf("TotalCents() = %d, want 0", got)
	}
}

func TestOrderItemCount(t *testing.T) {
	o := Order{
		Items: []LineItem{
			{Quantity: 2, UnitPriceCents: 1250},
			{Quantity: 1, UnitPriceCents: 600},
		},
	}
	if got, want := o.ItemCount(), 3; got != want {
		t.Fatalf("ItemCount() = %d, want %d", got, want)
	}
}

func TestSeededOrderTotals(t *testing.T) {
	s := SeedStore()

	cases := []struct {
		id   int
		want int64
	}{
		{1041, 3100},
		{1042, 10010},
		{1043, 11610},
		{1044, 5400},
		{1045, 5440},
	}
	for _, tc := range cases {
		o, ok := s.Get(tc.id)
		if !ok {
			t.Fatalf("order %d missing from seed data", tc.id)
		}
		if got := o.TotalCents(); got != tc.want {
			t.Errorf("order %d TotalCents() = %d, want %d", tc.id, got, tc.want)
		}
	}
}

func TestStoreAllSortedByID(t *testing.T) {
	s := SeedStore()
	all := s.All()
	if len(all) == 0 {
		t.Fatal("All() returned no orders")
	}
	for i := 1; i < len(all); i++ {
		if all[i-1].ID >= all[i].ID {
			t.Fatalf("All() not sorted: %d before %d", all[i-1].ID, all[i].ID)
		}
	}
}

func TestStoreGetUnknownID(t *testing.T) {
	s := SeedStore()
	if _, ok := s.Get(1); ok {
		t.Fatal("Get(1) reported a hit")
	}
}

func TestStorePutOverwrites(t *testing.T) {
	s := NewStore()
	s.Put(Order{ID: 7, Customer: "First"})
	s.Put(Order{ID: 7, Customer: "Second"})

	if n := len(s.All()); n != 1 {
		t.Fatalf("All() returned %d orders, want 1", n)
	}
	o, _ := s.Get(7)
	if o.Customer != "Second" {
		t.Fatalf("Customer = %q, want %q", o.Customer, "Second")
	}
}

func TestFormatCents(t *testing.T) {
	cases := []struct {
		cents int64
		want  string
	}{
		{0, "$0.00"},
		{5, "$0.05"},
		{600, "$6.00"},
		{1250, "$12.50"},
		{11610, "$116.10"},
	}
	for _, tc := range cases {
		if got := FormatCents(tc.cents); got != tc.want {
			t.Errorf("FormatCents(%d) = %q, want %q", tc.cents, got, tc.want)
		}
	}
}
