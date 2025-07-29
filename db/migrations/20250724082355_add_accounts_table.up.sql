CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    google_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    given_name VARCHAR(255),
    family_name VARCHAR(255),
    picture_url TEXT,
    verified_email BOOLEAN DEFAULT false,
    
    -- Account Management
    tier VARCHAR(50) DEFAULT 'free' CHECK (tier IN ('free', 'basic', 'premium', 'enterprise')),
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'cancelled', 'pending')),
    
    -- Features & Settings
    features JSONB DEFAULT '{}',
    
    -- Usage Tracking
    subscription_count INTEGER DEFAULT 0,
    last_login_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for better performance
CREATE INDEX idx_account_email ON account(email);
CREATE INDEX idx_account_google_id ON account(google_id);
CREATE INDEX idx_account_tier ON account(tier);
CREATE INDEX idx_account_status ON account(status);