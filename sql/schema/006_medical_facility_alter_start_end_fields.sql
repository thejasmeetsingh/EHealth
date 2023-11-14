-- +goose Up
ALTER TABLE medical_facility_timings ALTER COLUMN start_datetime TYPE TIME;
ALTER TABLE medical_facility_timings ALTER COLUMN end_datetime TYPE TIME;