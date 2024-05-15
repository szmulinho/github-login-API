# Github-login API

This API allows you to log in to the system by your github account. It also checks your account if you having access to the "Szmul-med" repository.
Users who have it are automatically assigned the 'doctor' role, granting them
specialized privileges. If you aren't access is not granted and the role
assigned is 'user.'

## Endpoints

### Callback

Endpoint: /callback

Description: It handles the callback from GitHub's OAuth authentication process.

### Get Data

Method: GET

Description: Retrieves data of a specific github's user.

## Testing

Unit and integration tests are implemented for the API. Tests can be run individually or using Docker Compose.

### Run Tests Individually

To run the tests manually, use the following command:

```go test ./...```

### Run Tests with Docker Compose

A docker-compose.yml file is provided to facilitate testing in a containerized environment. To run the tests using Docker Compose, use the following command:

```docker-compose up``` or ```make doker-tests```

This command will build the Docker images, run the tests, and then stop the containers.

## Setup Instructions

### 1. Clone the repository:

```git clone https://github.com/szmulinho/github-login-API.git```
```cd github-login-api```

### 2. Install dependencies:

```go mod tidy```

### 3. Run the server:

```go run main.go```

## Contribution Guidelines

Fork the repository

Create a new branch ```git checkout -b feature/your-feature-name```

Commit your changes ```git commit -m 'Add some feature'```

Push to the branch ```git push origin feature/your-feature-name```

Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

For any questions or suggestions, feel free to open an issue or contact the repository owner.



