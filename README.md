# Gator Program

## Overview
Gator is a command-line program for managing and aggregating feeds. Users can register, log in, follow feeds, and browse aggregated posts. The program interacts with a PostgreSQL database, and tools such as `goose` (for database migrations) and `sqlc` (for generating type-safe database access code) are used.

---

## Features

### User Commands
- **`register <user>`**: Registers a new user in the system.
- **`login <user>`**: Logs in an existing user.

### Administrative Commands
- **`reset`**: Erases the entire database.

### General Commands
- **`users`**: Lists all registered users.
- **`agg`**: Aggregates posts from the feeds.
- **`feeds`**: Lists all available feeds.

### Feed Management
- **`addfeed <name of feed> <feed url>`**: Adds a new feed to the system.
- **`follow <feed>`**: Allows the currently logged-in user to follow a specific feed.
- **`following`**: Lists all the feeds the currently logged-in user is following.
- **`unfollow <feed>`**: Unfollows a specific feed.

### Browsing
- **`browse <limit(optional)>`**: Displays posts aggregated from the user's followed feeds. Defaults to showing 2 posts if no limit is specified.

---

## Setup and Installation

### Prerequisites
- **PostgreSQL**: Ensure PostgreSQL is installed and running.
- **Go**: Install Go on your system.
- **Goose**: Database migration tool.
- **SQLC**: For generating type-safe database code.

### Installation Steps
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd aggregator
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure the database:
   - Create a PostgreSQL database for the program.
   - Update the `DATABASE_URL` environment variable with your database credentials:
     ```bash
     export DATABASE_URL="postgres://<username>:<password>@localhost:5432/<dbname>"
     ```

4. Apply database migrations using `goose`:
   ```bash
   goose postgres "$DATABASE_URL" up
   ```

5. Generate database access code using `sqlc`:
   ```bash
   sqlc generate
   ```

6. Run the program:
   ```bash
   go run main.go
   ```

---

## Commands and Usage

### Register a New User
```bash
register <username>
```
- Example:
  ```bash
  register john
  ```

### Log In as a User
```bash
login <username>
```
- Example:
  ```bash
  login john
  ```

### Reset the Database
```bash
reset
```
- Example:
  ```bash
  reset
  ```

### List All Users
```bash
users
```
- Example:
  ```bash
  users
  ```

### Aggregate Posts
```bash
agg
```
- Example:
  ```bash
  agg
  ```

### List All Feeds
```bash
feeds
```
- Example:
  ```bash
  feeds
  ```

### Add a New Feed
```bash
addfeed <name> <feed_url>
```
- Example:
  ```bash
  addfeed TechNews https://example.com/rss
  ```

### Follow a Feed
```bash
follow <feed_name>
```
- Example:
  ```bash
  follow TechNews
  ```

### List Feeds Followed by the Current User
```bash
following
```
- Example:
  ```bash
  following
  ```

### Unfollow a Feed
```bash
unfollow <feed_name>
```
- Example:
  ```bash
  unfollow TechNews
  ```

### Browse Posts
```bash
browse <limit(optional)>
```
- Example:
  - Default limit:
    ```bash
    browse
    ```
  - Specify a limit:
    ```bash
    browse 5
    ```

---

## Database Schema

### Users Table
- **`id`**: UUID (Primary Key)
- **`username`**: String (Unique)
- **`created_at`**: Timestamp
- **`updated_at`**: Timestamp

### Feeds Table
- **`id`**: UUID (Primary Key)
- **`name`**: String (Unique)
- **`url`**: String
- **`created_at`**: Timestamp
- **`updated_at`**: Timestamp

### Feed Follows Table
- **`id`**: UUID (Primary Key)
- **`user_id`**: UUID (Foreign Key referencing `users.id`)
- **`feed_id`**: UUID (Foreign Key referencing `feeds.id`)
- **`created_at`**: Timestamp
- **`last_fetched_at`**: Timestamp (Nullable)

### Posts Table
- **`id`**: UUID (Primary Key)
- **`feed_id`**: UUID (Foreign Key referencing `feeds.id`)
- **`title`**: String
- **`content`**: Text
- **`published_at`**: Timestamp

---

## Tools Used

### Goose
Goose is used for managing database migrations. To apply migrations:
```bash
goose postgres "$DATABASE_URL" up
```

### SQLC
SQLC generates type-safe Go code from SQL queries. To generate the code:
```bash
sqlc generate
```

---

## Contributing
Feel free to submit issues or contribute by creating pull requests. Please ensure your code follows Go best practices and includes appropriate tests.

---

## License
This project is licensed under the MIT License.


