CREATE TABLE reports (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    format VARCHAR(20) NOT NULL CHECK (
    format IN ('csv', 'json', 'xlsx', 'pdf')
    ),
    file_path TEXT NOT NULL,
    file_name TEXT NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_reports_user_id
    ON reports(user_id);

CREATE INDEX idx_reports_created_at
    ON reports(created_at);