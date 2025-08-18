CREATE TABLE monthly_active_subscriptions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    subscription_channel_id INT NOT NULL,
    month_due_date DATE NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (account_id) REFERENCES account(id),
    FOREIGN KEY (subscription_channel_id) REFERENCES subscription_channels(id)
);