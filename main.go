package main

import (
	"errors"
	"fmt"
	"log"
)

func main() {

	productRepo := NewInMemoryProductRepository()
	orderRepo := NewInMemoryOrderRepository()
	orderService := NewOrderService(productRepo, orderRepo)


	setupInitialProducts(productRepo)

	fmt.Println("--- Initial Stock ---")
	printCurrentStock(productRepo)
	fmt.Println("-----------------------")


	fmt.Println("\n--- 1. Scenario: Creating a Valid Order ---")
	validOrderRequest := CreateOrderRequest{
		Customer: "Ana",
		Items: []CreateOrderItemRequest{
			{ProductID: "P001", Quantity: 1},
			{ProductID: "P002", Quantity: 2},
		},
	}
	createdOrder, err := orderService.CreateOrder(validOrderRequest)
	if err != nil {
		log.Fatalf("Unexpected error while creating order: %v", err)
	}
	fmt.Printf("Order created successfully! ID: %s, Customer: %s, Status: %s, Total: R$%.2f\n",
		createdOrder.ID, createdOrder.Customer, createdOrder.Status, createdOrder.Total())

	fmt.Println("\n--- Stock after Ana's Order ---")
	printCurrentStock(productRepo)
	fmt.Println("------------------------------------")

	fmt.Println("\n--- 2. Scenario: Paying the Order ---")
	paidOrder, err := orderService.PayOrder(createdOrder.ID)
	if err != nil {
		log.Fatalf("Unexpected error while paying order: %v", err)
	}
	fmt.Printf("Order %s paid successfully! New status: %s\n", paidOrder.ID, paidOrder.Status)



	fmt.Println("\n--- 3. Scenario: Attempting to Cancel a Paid Order ---")
	_, err = orderService.CancelOrder(createdOrder.ID)
	if err != nil {
		fmt.Printf("Error (expected) while canceling paid order: %v\n", err)
		if !errors.Is(err, ErrInvalidStatusChange) {
			log.Printf("WARNING: The returned error was not the expected 'ErrInvalidStatusChange'.")
		}
	} else {
		log.Fatal("Unexpected error: A paid order was canceled, which should not be possible.")
	}

	fmt.Println("\n--- 4. Scenario: Attempting to Create an Order with Insufficient Stock ---")
	insufficientStockRequest := CreateOrderRequest{
		Customer: "Carlos",
		Items: []CreateOrderItemRequest{
			{ProductID: "P001", Quantity: 10},
		},
	}
	_, err = orderService.CreateOrder(insufficientStockRequest)
	if err != nil {
		fmt.Printf("Error (expected) while creating order with insufficient stock: %v\n", err)
		if !errors.Is(err, ErrInsufficientStock) {
			log.Printf("WARNING: The returned error was not the expected 'ErrInsufficientStock'.")
		}
	} else {
		log.Fatal("Unexpected error: An order with insufficient stock was created.")
	}

	fmt.Println("\n--- 5. Scenario: Attempting to Create an Order with Invalid Customer ---")
	invalidCustomerRequest := CreateOrderRequest{
		Customer: "",
		Items: []CreateOrderItemRequest{
			{ProductID: "P003", Quantity: 1},
		},
	}
	_, err = orderService.CreateOrder(invalidCustomerRequest)
	if err != nil {
		fmt.Printf("Error (expected) while creating order with invalid customer: %v\n", err)
		if !errors.Is(err, ErrInvalidCustomer) {
			log.Printf("WARNING: The returned error was not the expected 'ErrInvalidCustomer'.")
		}
	} else {
		log.Fatal("Unexpected error: An order with invalid customer was created.")
	}

	fmt.Println("\n--- 6. Scenario: Attempting to Create an Order with Multiple Errors ---")
	multiErrorRequest := CreateOrderRequest{
		Customer: "Jorge",
		Items: []CreateOrderItemRequest{
			{ProductID: "P999", Quantity: 1},
			{ProductID: "P001", Quantity: 99},
			{ProductID: "P002", Quantity: -1},
		},
	}
	_, err = orderService.CreateOrder(multiErrorRequest)
	if err != nil {
		fmt.Printf("Error (expected) with multiple validation issues:\n%v\n", err)
		if !errors.Is(err, ErrProductNotFound) || !errors.Is(err, ErrInsufficientStock) || !errors.Is(err, ErrInvalidQuantity) {
			log.Printf("WARNING: The joined error does not contain all expected error types.")
		}
	} else {
		log.Fatal("Unexpected error: An order with multiple validation errors was created.")
	}


	fmt.Println("\n--- Final Stock (should not have changed in error scenarios) ---")
	printCurrentStock(productRepo)
	fmt.Println("----------------------------------------------------------------")


	fmt.Println("\n--- 6. Scenario: Creating More Orders for Filtering ---")

	highValueOrder, _ := orderService.CreateOrder(CreateOrderRequest{
		Customer: "Beatriz",
		Items: []CreateOrderItemRequest{
			{ProductID: "P001", Quantity: 2},
		},
	})

	orderToCancel, _ := orderService.CreateOrder(CreateOrderRequest{
		Customer: "Daniel",
		Items:    []CreateOrderItemRequest{{ProductID: "P003", Quantity: 1}},
	})
	orderService.CancelOrder(orderToCancel.ID)

	fmt.Println("Additional orders created for filtering tests.")
	fmt.Printf("High-value order for Beatriz created. Subtotal: R$%.2f, Discount: R$%.2f, Final Total: R$%.2f\n",
		highValueOrder.Subtotal(),
		highValueOrder.Discount(),
		highValueOrder.Total())

	fmt.Println("\n--- Filtering Paid Orders ---")
	paidOrders, _ := orderService.ListOrders(func(order *Order) bool {
		return order.Status == StatusPaid
	})
	fmt.Printf("Found %d paid order(s):\n", len(paidOrders))
	for _, order := range paidOrders {
		fmt.Printf("- ID: %s, Customer: %s, Status: %s\n", order.ID, order.Customer, order.Status)
	}

	fmt.Println("\n--- Filtering Pending Orders ---")
	pendingOrders, _ := orderService.ListOrders(func(order *Order) bool {
		return order.Status == StatusPending
	})
	fmt.Printf("Found %d pending order(s):\n", len(pendingOrders))
	for _, order := range pendingOrders {
		fmt.Printf("- ID: %s, Customer: %s, Status: %s\n", order.ID, order.Customer, order.Status)
	}

	fmt.Println("\n--- Filtering Orders with Total Above R$5000 ---")
	highValueOrders, _ := orderService.ListOrders(func(order *Order) bool {
		return order.Total() > 5000.00
	})
	fmt.Printf("Found %d high-value order(s):\n", len(highValueOrders))
	for _, order := range highValueOrders {
		fmt.Printf("- ID: %s, Customer: %s, Total: R$%.2f\n", order.ID, order.Customer, order.Total())
	}
}

func setupInitialProducts(repo ProductRepository) {
	repo.Save(&Product{ID: "P001", Name: "Notebook", Price: 3500.00, Stock: 5})
	repo.Save(&Product{ID: "P002", Name: "Mouse", Price: 80.00, Stock: 10})
	repo.Save(&Product{ID: "P003", Name: "Teclado", Price: 180.00, Stock: 8})
}

func printCurrentStock(repo ProductRepository) {
	products, err := repo.List()
	if err != nil {
		log.Fatalf("Falha ao listar produtos: %v", err)
	}
	for _, p := range products {
		fmt.Printf("- Produto: %s (%s), Estoque: %d\n", p.ID, p.Name, p.Stock)
	}
}