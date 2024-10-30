#!/bin/bash

# Start the MySQL server
service mysql start

# Create the users table and populate it
mysql -u root -e "CREATE DATABASE spooky"
mysql -u root users < ./users.sql