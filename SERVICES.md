# Microservices Documentation

## Auth Service

### Get Session ID

- **Endpoint:** `GET /session_id`
- **Description:** Retrieves a new user session ID.
- **Cookie:** Session ID

### Update Session ID

- **Endpoint:** `UPDATE /session_id`
- **Description:** Updates the session ID to connect it to a username if the user is in the database and the password is correct.

### Login

- **Endpoint:** `POST /sign-up`
- **Description:** Creates a new user and corresponding entry in the PostgreSQL database.

### Get Coffee Shop ID

- **Endpoint:** `GET /coffee_shop_id`
- **Description:** Returns the coffee shop ID for authenticating the user when editing coffee shop information.
- **Parameters:** Session ID

### Get User ID

- **Endpoint:** `GET /user_id`
- **Description:** Returns the user ID for authenticating the user when leaving a check-in.
- **Parameters:** Session ID

## Check-in Service

### Post Check-in

- **Endpoint:** `POST /check_in`
- **Description:** Posts a check-in, logs into a logging service, and adds a check-in to Cassandra.
- **Parameters (Body):** User ID, Coffee Pack ID (optional), Coffee Shop ID (optional), Rating, Text (Optional)

### Get Check-in

- **Endpoint:** `GET /check_in`
- **Description:** Retrieves all/filtered check-ins from Cassandra.
- **Parameters:** User ID (optional), Coffee Pack ID (optional), Coffee Shop ID (optional)

## Coffee Pack Service

### Coffee Pack Model

- **Attributes:**
  - Name
  - ID
  - Description
  - Roastery
  - Image path
  - Weight
  - Flavour

### Add Coffee Pack

- **Endpoint:** `POST /packs`
- **Description:** Adds to the "packs" table in PostgreSQL.

### Get Coffee Packs

- **Endpoint:** `GET /packs`
- **Parameters:** Pack ID (path, optional)

### Get Coffee Pack by ID

- **Endpoint:** `GET /packs/:pack_id`

## Coffee Shops Service

### Coffee Shop Model

- **Attributes:**
  - Name
  - ID
  - Description
  - Address
  - Menu ID

### Menu Model

- **Attributes:**
  - Menu ID
  - Coffee Pack ID
  - Quantity
  - Price

### Add Coffee Shop

- **Endpoint:** `POST /shops`
- **Description:** Adds to the "shops" table in PostgreSQL.

### Get Coffee Shops

- **Endpoint:** `GET /shops`
- **Parameters:** Shop ID (path, optional)

### Get Coffee Shop by ID

- **Endpoint:** `GET /shops/:pack_id`
