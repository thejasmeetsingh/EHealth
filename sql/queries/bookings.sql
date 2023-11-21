-- name: CreateBooking :one
INSERT INTO bookings (id, created_at, modified_at, medical_facility_id, user_id, start_datetime, end_datetime, status) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateBookingStatus :one
UPDATE bookings SET status=$1 WHERE id=$2
RETURNING *;

-- name: GetBooking :one
SELECT * FROM bookings WHERE id=$1;

-- name: GetMedicalFacilityBookings :many
SELECT * FROM bookings WHERE medical_facility_id=$1 AND status=$2;

-- name: GetUserBookings :many
SELECT * FROM bookings WHERE user_id=$1 AND status=$2;

-- name: OverlappingUserBookings :many
SELECT id FROM bookings WHERE (start_datetime, end_datetime) OVERLAPS ($1, $2) AND status='A' AND user_id=$3;

-- name: OverlappingMedicalFacilityBookings :many
SELECT id FROM bookings WHERE (start_datetime, end_datetime) OVERLAPS ($1, $2) AND status='A' AND medical_facility_id=$3;