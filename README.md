# Cloning the Repository
To get started, follow these steps to clone the repository and set up the necessary environment variables:

1. Clone the repository to your local machine:

```bash
	git clone <repository_url>
```

2. Create a .env file in the root directory of the project.

3. Add the following required variables to the .env file:

- POSTGRES_URL: This variable should contain the URL that points to your PostgreSQL server. Make sure to replace <your_postgres_server_url> with the appropriate database connection URL.

- SESSION_TOKEN_KEY: Specify the name that you want your session cookies to have. Choose a unique and meaningful name for this variable.

- SERVER_URL: In a development environment, set this variable to "localhost." For production, it should point to the index domain of the server running your application.

Thank you for using Chick-fil-A Suite! We hope this web application simplifies and enhances your business systems for Chick-fil-A locations. Happy serving! 🐔🍟