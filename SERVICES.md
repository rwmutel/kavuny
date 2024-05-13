# Microservices Documentation

## Auth Service

Responsible for storing account information and binding it to a session identified by `session_id` stored in the 
cookies.
Account information is saved in PostgresSQL database.
Passwords are saved as a hashed value after adding the salt to an original password.
Sessions are stored in a distributed map of Hazelcast cluster.

### Endpoints

- `GET /session_id`
  - **Description:** Setups a session if it was not set up before, otherwise renews it. Sets a cookie `session_id`.
  - **Parameters:**
    - `session_id` (cookie): session id to renew (optional)
- `POST /log_in`
  - **Description:** Binds an account to a session_id. Session is automatically set up and renewed if necessary.
  - **Parameters:**
    - `session_id` (cookie): session id to bind with (optional)
    - `login` (query): login of the account
    - `password` (query): password of the account
- `POST /sign_up`
  - **Description:** Creates a new account and binds it to a session_id.
    Session is automatically set up and renewed if necessary.
  - **Parameters:**
    - `session_id` (cookie): session id to bind with (optional)
    - `login` (query): login of the account
    - `password` (query): password of the account
    - `user_type` (query): type of the account (`user` or `shop`)
- `GET /id`
  - **Description:** Retrieve the id from the account bound to the session id.
    Session is automatically set up and renewed if necessary.
  - **Parameters:**
    - `session_id` (cookie): session id (optional)

## Coffee Pack Service

Responsible for storing and returning the data about coffee packs.
Data is stored in a PostgresSQL database.

### Coffee Pack Model

  - `name`: string
  - `roastery`: string
  - `description`: string (optional)
  - `image_path`: string
  - `country`: string
  - `weight`: array of integers
  - `flavour`: array of strings

### Endpoints  

- `GET /packs`
  - **Description:** Returns the list of the coffee packs.
  - **Parameters:**
    - `ids` (query): ids of the packs to return (optional)
- `GET /packs/{id}`
  - **Description:** Returns the pack for specified `id`
  - **Parameters:**
    - `id` (path): id of the coffee pack
- `POST /packs`
  - **Description:** Adds to the "packs" table in PostgreSQL.
  - **Parameters:**
    - `coffee_pack` (body/json): coffee pack object
    - `session_id` (cookie): session id (optional)

## Coffee Shops Service

Responsible for storing and returning the data about coffee shops alongside their menus.
Data is stored in a PostgresSQL database.

### Coffee Shop Model

  - `id` - integer
  - `name` - string
  - `description` - string
  - `image_path` string
  - `address_text` string
  - `address_latitude` number
  - `address_longitude` number

### Menu Item Model

  - `coffee_pack_id` - integer
  - `price` - number
  - `quantity` - integer

### Endpoints

- `GET /coffee-shops`
  - **Description:** Returns the list of all coffee-shops.
- `GET /coffee-shops/{id}`
  - **Description:** Returns the coffee-shops for specified `id`
  - **Parameters:**
    - `id` (path): id of the shop
- `PUT /coffee-shops/{id}`
  - **Description:** Updates the coffee-shops for specified `id`
  - **Parameters:**
    - `id` (path): id of the shop
    - `coffee_shop` (body/json): coffee shop object
    - `session_id` (cookie): session id (optional)
- `GET /coffee-shops/{id}/menu`
  - **Description:** Returns the list of menu items for specified coffee shop `id`
  - **Parameters:**
    - `id` (path): id of the shop
- `POST /coffee-shops/{id}/menu`
  - **Description:** Adds the menu item for specified coffee shop `id`
  - **Parameters:**
    - `id` (path): id of the shop
    - `item` (body/json): menu item object
    - `session_id` (cookie): session id (optional)
- `DELETE /coffee-shops/{id}/menu`
  - **Description:** Deletes the menu item for specified coffee shop `id`
  - **Parameters:**
    - `id` (path): id of the shop
    - `item_id` (query): menu item id to delete
    - `session_id` (cookie): session id (optional)

## Check-in service

Responsible for storing and returning the data about check-ins.
Data is stored in a Cassandra keyspace.

### Check-in Model

- `coffee_shop_id`: integer (optional)
- `coffee_pack_id`: integer (optional)
- `check_in_time`: datetime
- `rating`: integer
- `check_in_text`: string (optional)

At least one of the `coffee_shop_id` and `coffee_pack_id` have to be specified.

### Endpoints

- `GET /check_ins`
  - **Description:** Returns the list of all check-ins for specified parameters.
  - **Parameters:**
    - `coffee_shop_id` (query): id of the shop to get check_ins for (optional)
    - `coffee_pack_id` (query): id of the pack to get check_ins for (optional)
    - `user_id` (query): id of the user to get check_ins for (optional)
- `POST /check_ins`
  - **Description:** Creates a check-in for logged-in user.
  - **Parameters:**
    - `check_in` (body/json): check-in to create
    - `session_id` (cookie): session id (optional)
