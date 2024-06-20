# challenge-beta

![Diagram](doc/diagram.svg)

### Requirements
- Docker and Docker Compose must be installed on your machine

### Running the Application
1. Build the Docker containers:
   ```sh
   docker-compose build
   ```
2. Start the containers in detached mode:
   ```sh
   docker-compose up -d
   ```

Now you can see the Golang API running on port 8080, the RabbitMQ management interface running on port 15672, and PostgreSQL on port 5432.

All applications are exposing ports because it is a test environment.

### Volumes
- Both PostgreSQL and RabbitMQ containers have persistent volumes. These are located in the root app folder:
    - `./volumes/rabbitmq/data`
    - `./volumes/postgres/data`

### Application
- The app will be available at `localhost:8080`. The endpoints are:
    - **POST** `/pedidos`: Generates a new order, stores it in the PostgreSQL database with the status "PENDENTE," and publishes this order to a RabbitMQ queue.
    - **GET** `/pedidos/:pedidoId`: Retrieves a single order with the specified ID.
    - **GET** `/pedidos`: Lists all orders.