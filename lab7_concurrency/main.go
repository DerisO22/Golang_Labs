package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	id     int
	status string
}

var (
	orderIDCounter int
	orderIDMu sync.Mutex
)

const MAX = 20

func main() {
	orders := make([]*Order, MAX)
	var wg sync.WaitGroup

	wg.Add(MAX)

	for i := 0; i < MAX; i++ {
		go func() {
			defer wg.Done()
			orders[i] = generateOrder()
		}()
	}

	wg.Wait()

	reportOrderStatus(orders)

	fmt.Println("\n--- Processing Orders ---")

	wg.Add(4)

	func() {
		defer wg.Done()
		go processOrders(orders)
	}()

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			for _, order := range orders {
				updateOrderStatus(order)
			}
		}()
	}

	wg.Wait()
	
	reportOrderStatus(orders)

	fmt.Println("All Done")
}

func generateOrder() *Order {
	var orderID int

	orderIDMu.Lock()
	orderID = orderIDCounter

	time.Sleep((time.Duration(rand.Intn(100))) * time.Microsecond)

	orderIDCounter++
	orderIDMu.Unlock()

	return &Order{
		id: orderID, status: "Received",
	}
}

func generateOrderStatus() string {
	status := []string{"Received", "Processing", "Served"}[rand.Intn(3)]

	return status
}

func processOrders(orders []*Order) {
	for _, order := range orders {
		time.Sleep((time.Duration(rand.Intn(500))) * time.Microsecond)
		fmt.Printf("Processing order %d\n", order.id)
	}
}

func reportOrderStatus(orders []*Order) {
	time.Sleep(1 * time.Microsecond)
	fmt.Println("\n--- Order Status ---")

	for _, order := range orders {
		fmt.Printf("Order %d: %s\n", order.id, order.status)
	}
}

func updateOrderStatus(order *Order) {
	time.Sleep((time.Duration(rand.Intn(300))) * time.Microsecond)

	order.status = generateOrderStatus()

	fmt.Printf("Updating order %d status to %s\n", order.id, order.status)
}
