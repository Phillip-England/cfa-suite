# About

Chick-fil-A Suite is a web application for Chick-fil-A restaurants. It contains features which are generally useful amongst a wide variaty of Chick-fil-A's.

Thank you for using Chick-fil-A Suite! We hope this web application simplifies and enhances your business systems for Chick-fil-A locations. Happy serving! 🐔🍟

# Cloning the Repository
To get started, follow these steps to clone the repository and set up the necessary environment variables:

1. Clone the repository to your local machine:

```bash
git clone https://github.com/phillip-england/cfa-suite
```

2. Create a .env file in the root directory of the project.

```bash
touch .env
```

3. Add the following required variables to the .env file:

- POSTGRES_URL: This variable should contain the URL that points to your PostgreSQL server. Make sure to replace <your_postgres_server_url> with the appropriate database connection URL.

- SESSION_TOKEN_KEY: Specify the name that you want your session cookies to have. Choose a unique and meaningful name for this variable.

- SERVER_URL: In a development environment, set this variable to "localhost." For production, it should point to the index domain of the server running your application.

- APP_EMAIL: The email address you wish account recovery to be sent through.

- APP_EMAIL_PASSWORD: The password of the email address you with account recovery to be sent through.

- SMTP_SERVER: The SMTP server address for your email provider.

- PORT: The port number you wish the application to be served on.

- GO_ENV: set to either "dev" or "prod" depending on the environment the application is ran in.

4. Install the tailwind dependancy

```bash
npm install
```

5. Run tailwind during active developement:

```bash
npx tailwindcss -i "./static/input.css" -o "./static/output.css" --watch
```