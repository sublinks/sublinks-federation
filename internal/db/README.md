# Database

Configure DSN in .env file in main directory of application

## Migrations

Create a files in /db/migrations named in the format:

- ##########_create_users_table.up.sql
- ##########_create_users_table.down.sql

Where ########## is an incrementing number. "up" is obviously the intended action file whereas "down" is the rollback.
