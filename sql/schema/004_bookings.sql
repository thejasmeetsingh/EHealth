-- +goose Up

-- ('P', 'Pending')
-- ('A', 'Accepted')
-- ('R', 'Rejected')
CREATE TYPE booking_status AS ENUM ('P', 'A', 'R');

CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    medical_facility_id UUID NOT NULL REFERENCES medical_facility(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP NOT NULL,
    status booking_status NOT NULL
);

-- +goose Down
DROP TABLE bookings;