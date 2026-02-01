-- Add downloads tracking table
CREATE TABLE IF NOT EXISTS downloads (
    id VARCHAR(36) PRIMARY KEY,
    book_id VARCHAR(36) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    downloaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT
);

-- Add index for faster queries
CREATE INDEX idx_downloads_book_id ON downloads(book_id);
CREATE INDEX idx_downloads_user_id ON downloads(user_id);
CREATE INDEX idx_downloads_downloaded_at ON downloads(downloaded_at);

-- Add download_count column to books table
ALTER TABLE books ADD COLUMN IF NOT EXISTS download_count INTEGER DEFAULT 0;

-- Create index on download_count for sorting
CREATE INDEX idx_books_download_count ON books(download_count DESC);

-- Function to update download count
CREATE OR REPLACE FUNCTION update_book_download_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE books
    SET download_count = (
        SELECT COUNT(*) FROM downloads WHERE book_id = NEW.book_id
    )
    WHERE id = NEW.book_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update download count
DROP TRIGGER IF EXISTS trigger_update_download_count ON downloads;
CREATE TRIGGER trigger_update_download_count
AFTER INSERT ON downloads
FOR EACH ROW
EXECUTE FUNCTION update_book_download_count();
