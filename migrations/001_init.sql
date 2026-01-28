-- Library Management System - Complete Migration
-- This will drop existing tables and recreate them
-- Run: psql -U postgres -d bookgolang -f migrations/001_full_migration.sql

\echo 'Starting migration...'
\echo ''

-- ===== STEP 1: DROP EXISTING TABLES =====
\echo 'Step 1: Dropping existing tables...'

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users CASCADE;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories CASCADE;
DROP TRIGGER IF EXISTS update_books_updated_at ON books CASCADE;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments CASCADE;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Drop tables (order matters: child tables first)
DROP TABLE IF EXISTS likes CASCADE;
DROP TABLE IF EXISTS comments CASCADE;
DROP TABLE IF EXISTS saved_books CASCADE;
DROP TABLE IF EXISTS books CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS users CASCADE;

\echo '✓ Existing tables dropped'
\echo ''

-- ===== STEP 2: CREATE TABLES =====
\echo 'Step 2: Creating tables...'

-- Create users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'member')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

\echo '✓ Users table created'

-- Create categories table
CREATE TABLE categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

\echo '✓ Categories table created'

-- Create books table
CREATE TABLE books (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    pdf_file VARCHAR(500) NOT NULL,
    category_id VARCHAR(36) NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    owner_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    like_count INTEGER DEFAULT 0,
    dislike_count INTEGER DEFAULT 0,
    save_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

\echo '✓ Books table created'

-- Create saved_books table
CREATE TABLE saved_books (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id VARCHAR(36) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, book_id)
);

\echo '✓ Saved books table created'

-- Create comments table
CREATE TABLE comments (
    id VARCHAR(36) PRIMARY KEY,
    book_id VARCHAR(36) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    parent_id VARCHAR(36) REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

\echo '✓ Comments table created'

-- Create likes table
CREATE TABLE likes (
    id VARCHAR(36) PRIMARY KEY,
    book_id VARCHAR(36) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_like BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, book_id)
);

\echo '✓ Likes table created'
\echo ''

-- ===== STEP 3: CREATE INDEXES =====
\echo 'Step 3: Creating indexes...'

CREATE INDEX idx_books_category ON books(category_id);
CREATE INDEX idx_books_owner ON books(owner_id);
CREATE INDEX idx_books_save_count ON books(save_count DESC);
CREATE INDEX idx_saved_books_user ON saved_books(user_id);
CREATE INDEX idx_saved_books_book ON saved_books(book_id);
CREATE INDEX idx_comments_book ON comments(book_id);
CREATE INDEX idx_comments_parent ON comments(parent_id);
CREATE INDEX idx_likes_book ON likes(book_id);
CREATE INDEX idx_likes_user ON likes(user_id);

\echo '✓ Indexes created'
\echo ''

-- ===== STEP 4: CREATE TRIGGERS =====
\echo 'Step 4: Creating triggers...'

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

\echo '✓ Function created'

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at 
    BEFORE UPDATE ON categories
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_books_updated_at 
    BEFORE UPDATE ON books
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_comments_updated_at 
    BEFORE UPDATE ON comments
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

\echo '✓ Triggers created'
\echo ''

-- ===== STEP 5: VERIFY TABLES =====
\echo 'Step 5: Verifying tables...'

SELECT 
    table_name,
    (SELECT COUNT(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
FROM information_schema.tables t
WHERE table_schema = 'public' 
    AND table_type = 'BASE TABLE'
ORDER BY table_name;

\echo ''
\echo '========================================='
\echo 'Migration completed successfully! ✓'
\echo '========================================='
\echo ''
\echo 'Tables created:'
\echo '  - users'
\echo '  - categories'
\echo '  - books'
\echo '  - saved_books'
\echo '  - comments'
\echo '  - likes'
\echo ''
\echo 'You can now start the server:'
\echo '  go run cmd/main.go'
\echo ''