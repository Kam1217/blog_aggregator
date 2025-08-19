# Gator CLI
A command-line RSS feed aggregator that allows you to follow and browse posts from multiple RSS feeds.

## Prerequisites
Before running this program, you'll need to have the following installed on your system:

### PostgreSQL
- **macOS**: Install via Homebrew: `brew install postgresql`
- **Ubuntu/Debian**: `sudo apt-get install postgresql postgresql-contrib`
- **Windows**: Download from [postgresql.org](https://www.postgresql.org/download/windows/)

### Go
- Download and install Go from [golang.org](https://golang.org/dl/)
- Verify installation: `go version`

## Installation

Install the Gator CLI using Go's built-in package manager:

```bash
go install github.com/Kam1217/blog_aggregator
```
## Commands 

### User Management
- register: Create a new user account within the gator system
- login: Log in an existing user. Many commands require a user to be logged in
- users: List all currently registered users
- reset: Reset the application's state, such as user data or the database

### Feed Management
- addfeed (Requires login): Add a new RSS feed to the system for tracking
- feeds: List all RSS feeds currently configured in the system
- follow (Requires login): Subscribe to a specific RSS feed to see its posts
- following (Requires login): Display all RSS feeds that the current user is following
- unfollow (Requires login): Stop following a specific RSS feed

### Content
- agg: Trigger the aggregation process every x amount of time, which fetches and processes new posts from all configured RSS feeds
- browse (Requires login): Browse through the posts collected from the feeds the current user follows

## Usage Example

### Create a new user
go run . register <username>

### Log in
go run . login <username>

### Add an RSS feed
go run . addfeed <feed_name> <feed_url>

### Follow a feed
go run . follow <feed_name>

### Aggregate new posts
go run . agg <time>

### Browse your posts
go run . browse <optional - how many you posts you wish to see>
