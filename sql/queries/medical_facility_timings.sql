-- name: AddMedicalFacilityTimings :one
INSERT INTO medical_facility_timings (
    id,
    created_at,
    modified_at,
    medical_facility_id,
    weekday,
    start_datetime,
    end_datetime
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetMedicalFacilityTimingDetails :many
SELECT * FROM medical_facility_timings WHERE medical_facility_id=$1;

-- name: GetMedicalFacilityTimingById :one
SELECT * FROM medical_facility_timings WHERE id=$1;

-- name: UpdateMedicalFacilityTimings :one
UPDATE medical_facility_timings SET
weekday=$1,
start_datetime=$2,
end_datetime=$3
WHERE id=$4 
RETURNING *;