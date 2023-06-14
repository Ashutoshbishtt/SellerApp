## Project Overview

The project aims to create a service that supports order management with various functionalities such as adding orders, updating order status, pagination, filtering, and sorting. The application has been developed with the following key features:

- **Order Management**: The service allows users to add orders with different status values. It also supports updating the status of existing orders.

- **Pagination**: The application provides pagination functionality, allowing users to retrieve orders in smaller, manageable chunks.

- **Filtering**: Users can apply filters to retrieve specific orders based on various criteria such as order ID, status, total amount, currency unit, and more.

- **Sorting**: The service supports sorting of orders based on different fields, including order ID, status, total amount, currency unit, and item fields.

- **Tests**: The service includes comprehensive tests to ensure the functionality and integrity of the implemented features. These tests cover various scenarios and edge cases to validate the service's behavior.

- **Documentation**: Documentation has been provided to guide users on running the service and understanding its functionality. The documentation is available in both code comments and on the repository.

- **Containerization**: The application has been containerized using Docker, allowing for easy deployment and scalability.

- **Minimal External Dependencies**: The service has been developed with a focus on minimizing external dependencies. This approach enhances the portability and reduces potential conflicts or issues related to external libraries.

## Project Setup and Execution

### Prerequisites

- Node.js installed on your machine.
- MySQL Workbench installed for database management.
- Go installed for running the backend server.
- Docker installed for containerization.
- Postman installed for API testing.

### Frontend Setup

```shell
cd client
npm install
npm run build
