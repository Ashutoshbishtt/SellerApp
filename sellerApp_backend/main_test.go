package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStoreOrder(t *testing.T) {
	payload := []byte(`[
		{
			"id": "order-001",
			"status": "pending",
			"items": [
				{
					"id": "item-001",
					"description": "Item 1",
					"price": 10.99,
					"quantity": 2
				},
				{
					"id": "item-002",
					"description": "Item 2",
					"price": 5.99,
					"quantity": 3
				}
			],
			"total": 39.95,
			"currencyUnit": "USD"
		}
	]`)

	req, err := http.NewRequest("POST", "/storeOrder", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(storeOrder)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	expected := "Order(s) stored successfully"
	if rr.Body.String() != expected {
		t.Errorf("expected body to be %q; got %q", expected, rr.Body.String())
	}
}

func TestGetOrder(t *testing.T) {
	req, err := http.NewRequest("GET", "/getOrder?id=order-001", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getOrder)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	var response OrderPayload
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedID := "order-001"
	if response.ID != expectedID {
		t.Errorf("expected ID to be %q; got %q", expectedID, response.ID)
	}
	// Add more assertions for other fields if needed
}

func TestGetAllOrders(t *testing.T) {
	req, err := http.NewRequest("GET", "/getAllOrders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllOrders)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	var response []OrderPayload
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Add assertions to check the response data
}

func TestUpdateOrderStatus(t *testing.T) {
	payload := []byte(`{
		"id": "order-001",
		"status": "completed"
	}`)

	req, err := http.NewRequest("POST", "/updateOrderStatus", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateOrderStatus)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	expected := "Order status updated successfully"
	if rr.Body.String() != expected {
		t.Errorf("expected body to be %q; got %q", expected, rr.Body.String())
	}
}

func TestFindOrderByID(t *testing.T) {
	orders := []OrderPayload{
		{
			ID: "order-001",
		},
		{
			ID: "order-002",
		},
	}

	order := findOrderByID(orders, "order-001")
	if order == nil {
		t.Error("expected order with ID 'order-001', but not found")
	}

	order = findOrderByID(orders, "order-003")
	if order != nil {
		t.Error("expected nil order, but found one")
	}
}
