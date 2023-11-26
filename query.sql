-- name: GetOTP :one
SELECT otp, otp_expiration_time FROM users
WHERE phone_number = $1;

-- name: CheckPhoneExistence :one
SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = $1);


-- name: CreateUser :exec
INSERT INTO users (
    name, phone_number
) VALUES (
    $1, $2
);

-- name: UpdateUserOTP :exec
UPDATE users set otp = $2, otp_expiration_time = $3 WHERE phone_number = $1;