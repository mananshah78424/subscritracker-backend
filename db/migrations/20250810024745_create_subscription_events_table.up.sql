CREATE TABLE subscription_events (
    id SERIAL PRIMARY KEY,
    subscription_details_id INT NOT NULL REFERENCES subscription_details(id),
    account_id INT NOT NULL REFERENCES account(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_subscription_events_subscription_details_id ON subscription_events(subscription_details_id);
CREATE INDEX idx_subscription_events_account_id ON subscription_events(account_id);
CREATE INDEX idx_subscription_events_created_at ON subscription_events(created_at);
CREATE INDEX idx_subscription_events_updated_at ON subscription_events(updated_at);



