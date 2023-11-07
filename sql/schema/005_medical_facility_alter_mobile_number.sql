-- +goose Up
ALTER TABLE medical_facility ALTER COLUMN mobile_number TYPE VARCHAR(20) USING mobile_number::VARCHAR(20);