-- +goose Up

-- (0, 'Sunday')
-- (1, 'Monday')
-- (2, 'Tuesday')
-- (3, 'Wednesday')
-- (4, 'Thursday')
-- (5, 'Friday')
-- (6, 'Saturday')
CREATE TYPE weekday_type AS ENUM ('0', '1', '2', '3', '4', '5', '6');

CREATE TABLE medical_facility_timings (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    medical_facility_id UUID NOT NULL REFERENCES medical_facility(id) ON DELETE CASCADE,
    weekday weekday_type NOT NULL,
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP NOT NULL,
    CONSTRAINT UniqueWeekdays UNIQUE (medical_facility_id, weekday)
);

-- +goose Down
DROP TABLE medical_facility_timings;