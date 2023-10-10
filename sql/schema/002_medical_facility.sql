-- +goose Up

-- ('ID', 'Individual Doctor')
-- ('C', 'Clinic')
-- ('H', 'Hospital')
CREATE TYPE facility_type AS ENUM ('ID', 'C', 'H');

CREATE TABLE medical_facility (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    type facility_type NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    email VARCHAR(50) NOT NULL,
    mobile_number VARCHAR(11) NOT NULL,
    charges DECIMAL(2) NOT NULL,
    address TEXT NOT NULL,
    location GEOMETRY(POINT, 4326) NOT NULL,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE medical_facility;