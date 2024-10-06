# log_service

### Setup Instructions

1. Copy the contents of `.env.sample` and create a new `.env` file. Update the environment variables in `.env` as per your requirements.
2. To bring up the containers, run the command:
   ```sh
   make up
   ```
3. To perform the database migration, run:
   ```sh
   make migrate_up
   ```
