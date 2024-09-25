# Overview

Gator is a RSS aggregator.

## Prereqs

- Go
- Postgres

## Installation

```bash
go install ...
```

## Config

Create a `.gatorconfig.json` file in your home directory with the structure:

```json
{
    "db_url": "postgres//username:@localhost:5432/database?sslmode=disable"
}
```
Update the connection string accordingly

## Usage

Register a user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start aggregation:

```bash
gator agg 1m
```

Browse posts:

```bash
gator browse [limit]
```

Change login:

```bash
gator login <name>

List feeds:

```bash
gator feeds
```

List users:

```bash
gator users
```

Follow a feed that already exists in the database

```bash
gator follow <url>
```

Unfollow a feed that already exists in the database

```bash
gator unfollow <url>
```

