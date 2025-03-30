# SWIFT Code API

## Overview
This project provides a RESTful API to manage **SWIFT codes** for banks. It allows users to:
- Retrieve bank details by SWIFT code.
- Fetch all SWIFT codes for a specific country.
- Add new SWIFT codes.
- Delete existing SWIFT codes.


## Technologies Used
- **Golang** (Gorilla Mux, GORM)
- **PostgreSQL** (Primary database)
- **SQLite** (In-memory DB for testing)
- **Docker & Docker Compose**

---

## Setup & Installation

### **1️⃣ Clone the Repository**
```sh
git clone https://github.com/Unakepro/SwiftParser.git
cd SwiftParser
```

### **2️⃣ Setup Environment Variables**
Create a `.env` file in folder and configure the database:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=swift_db
```
Note: for docker you have to use DB_HOST=db

### **3️⃣ Start Project using Docker**

Production mode
```sh
docker-compose up --build app
```

Test mode
```sh
docker-compose up --build test
```


## API Endpoints

### **1️⃣ Get SWIFT Code Details**
**GET** `/swift-code/{swiftCode}`
#### **Response:**
```json
{
  "Address": "123 Wall Street",
  "BankName": "Bank USA",
  "ISO2Code": "US",
  "CountryName": "United States",
  "IsHeadquarter": true,
  "SwiftCode": "TESTUS1",
   "Branches": [
    {
      "Address": "456 Elm Street",
      "BankName": "Bank USA Branch",
      "ISO2Code": "US",
      "IsHeadquarter": false,
      "SwiftCode": "TESTUS2",
    }
  ]
}
```

### **2️⃣ Get All SWIFT Codes for a Country**
**GET** `/swift-codes/{countryISO2}`

### **3️⃣ Add a New SWIFT Code**
**POST** `/swift-code`
#### **Request:**
```json
{
  "SwiftCode": "NEWCODE",
  "BankName": "New Bank",
  "Address": "123 Bank St",
  "ISO2Code": "MC",
  "IsHeadquarter": true
}
```

### **4️⃣ Delete a SWIFT Code**
**DELETE** `/swift-code/{swiftCode}`

---


## License
This project is licensed under the **MIT License**.

