CREATE TABLE subscription_details (
    id SERIAL PRIMARY KEY NOT NULL,
    subscription_channel_id INT NOT NULL,
    start_date DATE,
    due_date DATE,
    start_time TIME,
    due_time TIME,
    monthly_bill NUMERIC(10, 2),
    reminder_date DATE,
    reminder_time TIME,
    status VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_subscription_details_start_date ON subscription_details(start_date);
CREATE INDEX idx_subscription_details_due_date ON subscription_details(due_date);
CREATE INDEX idx_subscription_details_status ON subscription_details(status);
CREATE INDEX idx_subscription_details_reminder_date ON subscription_details(reminder_date);
CREATE INDEX idx_subscription_details_reminder_time ON subscription_details(reminder_time);

