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
SELECT * FROM bookings WHERE medical_facility_id=$1 AND status=$2 LIMIT $3 OFFSET $4;

-- name: GetUserBookings :many
SELECT * FROM bookings WHERE user_id=$1 AND status=$2 LIMIT $3 OFFSET $4;

-- name: OverlappingPendingBookings :many
SELECT b.id, b.start_datetime, b.end_datetime, mf.name, mf.address, u.email FROM bookings b 
JOIN medical_facility mf ON b.medical_facility_id=mf.id
JOIN users u ON b.user_id=u.id
WHERE (b.start_datetime, b.end_datetime) OVERLAPS ($1, $2) AND status='P' AND b.medical_facility_id=$3 AND b.id!=$4;

-- name: OverlappingAcceptedBookingCount :one
SELECT COUNT(id) FROM bookings WHERE (start_datetime, end_datetime) OVERLAPS ($1, $2) AND status='A' AND medical_facility_id=$3;