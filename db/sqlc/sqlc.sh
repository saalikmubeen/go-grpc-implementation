# SQLC (SQL Compiler) is a tool that generates Go code from SQL queries. 
# It is a tool that helps you write type-safe SQL queries in Go.
# Here's how it works:

# 1. You write SQL queries
# 2. You run sqlc to generate Go code that presents type-safe interfaces 
#    to those queries
# 3. You write application code that calls the methods sqlc generated

# You don’t have to write any boilerplate SQL querying code.

# Install 
brew install sqlc

sqlc -h

sqlc.init # creates the sqlc.yaml file

# generates the GO code using the settings defined sqlc.yaml file
# and the SQL queries written in the db/query directory
sqlc generate 