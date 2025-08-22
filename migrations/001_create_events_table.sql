-- 001_create_events_table.sql
-- Migration: Create events table
-- Created: 2025-08-21

-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create events table
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create index on start_time for better query performance
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);

-- Create index on created_at for ordering
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);

-- Create trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Drop trigger if exists, then create it
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data
INSERT INTO events (title, description, start_time, end_time) VALUES 
    ('Go Conference', 'A conference about Go programming language', 
     NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day 3 hours'),
    ('Docker Workshop', 'Practical workshop on Docker and containers', 
     NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days 4 hours'),
    ('PostgreSQL Meetup', 'Database developers meetup and networking', 
     NOW() + INTERVAL '3 days', NOW() + INTERVAL '3 days 2 hours'),
    ('DevOps Summit', 'Annual DevOps best practices and tools summit', 
     NOW() + INTERVAL '4 days', NOW() + INTERVAL '4 days 6 hours'),
    ('JavaScript Bootcamp', 'Intensive training on modern JavaScript frameworks', 
     NOW() + INTERVAL '5 days', NOW() + INTERVAL '5 days 8 hours')
ON CONFLICT (id) DO NOTHING;

-- Display results using SELECT (works with non-interactive mode)
SELECT 'Events table created successfully!' as message;
SELECT 'Sample data inserted.' as message;
SELECT 'Migration 001 completed successfully!' as status;
