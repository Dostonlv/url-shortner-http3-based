
CREATE DATABASE urlshortener;

-- Migration script for creating URLs table
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    clicks INTEGER NOT NULL DEFAULT 0
);

-- Create index for faster lookups
CREATE INDEX idx_short_code ON urls(short_code);