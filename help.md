
## GraphQL
```graphql
mutation createOrder {
  createOrder(input: {id: "eeeee", Price: 50, Tax: 19}) {
    id
    Price
    Tax
    FinalPrice
  }
}

query queryOrders {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}

```