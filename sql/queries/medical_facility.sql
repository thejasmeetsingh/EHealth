-- name: CreateMedicalFacility :one
INSERT INTO medical_facility (
    id,
    created_at,
    modified_at,
    type,
    name,
    description,
    email,
    mobile_number,
    charges,
    address,
    location,
    user_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetMedicalFacilityById :one
SELECT * FROM medical_facility WHERE id=$1;

-- name: GetMedicalFacilityByUserId :one
SELECT * FROM medical_facility WHERE user_id=$1;

-- name: UpdateMedicalFacility :one
UPDATE medical_facility SET
type=$1,
name=$2,
description=$3,
email=$4,
mobile_number=$5,
charges=$6,
address=$7,
location=$8
WHERE id=$9 RETURNING *;