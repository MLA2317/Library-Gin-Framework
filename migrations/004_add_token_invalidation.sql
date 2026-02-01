-- Add token invalidation timestamp to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS token_invalidated_at TIMESTAMP DEFAULT NULL;
