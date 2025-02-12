# Pave Bank challenge

This is an assignment for the Pave band interview process. The task is to build a fees API in encore that uses a temporal workflow started at the beginning of a fee period, and allows for the progressive accrual of fees.

# Setup

**Install Encore:**
- **macOS:** `brew install encoredev/tap/encore`
- **Linux:** `curl -L https://encore.dev/install.sh | bash`
- **Windows:** `iwr https://encore.dev/install.ps1 | iex`


**Running the Application**

To run the application locally, use the following command:
```bash
encore run
```

**Testing the Application**

To test the application, use the following command:
```bash
encore test ./...
```


# Bill Service API

This API provides endpoints to manage bills, including creating, querying, closing bills, and adding or removing items from a bill.

## Endpoints

### Create a Bill

**URL:** `/bill`  
**Method:** `POST`  
**Description:** Creates a new bill.  
**Request Body:**
```json
{
  "currency": "USD"
}
```
**Response:**
```json
{
  "id": "BILL-1234567890",
  "currency": "USD",
  "items": [],
  "status": "open"
}
```

### Query a Bill

**URL:** `/bill/:billID`  
**Method:** `GET`  
**Description:** Retrieves the details of a specific bill.  
**Response:**
```json
{
  "id": "BILL-1234567890",
  "currency": "USD",
  "items": [],
  "status": "open"
}
```

### Close a Bill

**URL:** `/bill/:billID/close`  
**Method:** `POST`  
**Description:** Closes a specific bill.  
**Response:** `200 OK`

### Add Item to a Bill

**URL:** `/bill/:billID/items`  
**Method:** `POST`  
**Description:** Adds an item to a specific bill.  
**Request Body:**
```json
{
  "itemID": "ITEM-123",
  "name": "Item Name",
  "price": 100.5
}
```
**Response:** `200 OK`

### Remove Item from a Bill

**URL:** `/bill/:billID/items/:itemID`  
**Method:** `DELETE`  
**Description:** Removes an item from a specific bill.  
**Response:** `200 OK`
