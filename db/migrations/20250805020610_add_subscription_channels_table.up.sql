CREATE TABLE IF NOT EXISTS subscription_channels (
    id SERIAL PRIMARY KEY,
    channel_name VARCHAR(255) NOT NULL,
    channel_url TEXT NOT NULL,
    channel_status VARCHAR(50) DEFAULT 'active' CHECK (channel_status IN ('active', 'suspended', 'cancelled', 'pending')),
    channel_type VARCHAR(50),
    channel_image_url TEXT NOT NULL,
    channel_description TEXT NOT NULL,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for better performance
CREATE INDEX idx_subscription_channels_channel_name ON subscription_channels(channel_name);
CREATE INDEX idx_subscription_channels_channel_url ON subscription_channels(channel_url);
CREATE INDEX idx_subscription_channels_channel_status ON subscription_channels(channel_status);
CREATE INDEX idx_subscription_channels_channel_type ON subscription_channels(channel_type);