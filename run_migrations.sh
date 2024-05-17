#!/bin/sh

set -e # Exit on error. The script will exit immediately if any command fails(i.e returns non-zero status)

echo "RUNNING MIGRATIONS......!"


# Load the environment variables from the app.env file into the current shell
# of the dcoker container or into the current container's environment. 
# This will make the environment variables defined  in the app.env 
# file available to the shell script.
# But this will overwrite the environment variables defined in the 
# docker-compose file or in the Dockerfile if they have the same name.
# source /app/app.env

# Run migrations
# /app/migrate -path /app/db/migrations -database "$DB_URI" -verbose up # Read the DB_URI from the environment variale


echo "STARTING THE APPLICATION SERVER......!"

# Take all parameters passed to the script and execute them as a command
exec "$@" # Execute the CMD passed to the script

# The CMD passed to the script will be the command to start the application server:
# CMD ["/app/main"]