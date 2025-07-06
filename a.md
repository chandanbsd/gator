# Install Guide

## Dependencies

- Postgresql 17
- go 1.23.9

## Database migrations

- We use the goose database migration tool to create database tables

  - goose <database_driver> <connection_string> up

- For whatever reason if you need to wipe the data or reinstall then do
  - goose <database_drive> <connection_string> down

## Install

- go install github.com/chandanbsd/gator
