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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, ST_SetSRID(ST_MakePoint($11, $12), 4326), $13)
RETURNING *;

-- name: GetMedicalFacilityById :one
SELECT 
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
    ST_X(location) AS lat,
    ST_y(location) AS lng,
    user_id
FROM 
    medical_facility WHERE id=$1;

-- name: GetMedicalFacilityByUserId :one
SELECT 
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
    ST_X(location) AS lat,
    ST_y(location) AS lng,
    user_id
FROM 
    medical_facility WHERE user_id=$1;

-- name: UpdateMedicalFacility :one
UPDATE medical_facility SET
type=$1,
name=$2,
description=$3,
email=$4,
mobile_number=$5,
charges=$6,
address=$7,
location=ST_SetSRID(ST_MakePoint($8, $9), 4326)
WHERE id=$10 RETURNING *;

-- name: MedicalFacilityListing :many
SELECT
id,
type,
name,
charges,
address,
CAST(ST_DistanceSphere(location, ST_MakePoint($1, $2)) / 1000 AS FLOAT) AS distance
FROM medical_facility
ORDER BY distance;

-- name: MedicalFacilityDetail :one
SELECT
id,
type,
name,
description,
email,
mobile_number,
charges,
address,
ST_X(location) AS lat,
ST_y(location) AS lng,
CAST(ST_DistanceSphere(location, ST_MakePoint($1, $2)) / 1000 AS FLOAT) AS distance
FROM medical_facility
WHERE id=$3;