drop column if exists password_hash;
drop column if exists email_verified;
drop column if exists verification_token;
drop column if exists reset_token;
drop column if exists reset_token_expires;

drop index if exists idx_account_email_verified;
drop index if exists idx_account_verification_token;
drop index if exists idx_account_reset_token;