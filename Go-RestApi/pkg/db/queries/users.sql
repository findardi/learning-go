-- name: UserGetByID :one
select * from users where id = $1 limit 1;

-- name: UserGetByUsername :one
select * from users where username = $1 limit 1;

-- name: UserGetByEmail :one
select * from users where email = $1 limit 1;

-- name: UserInsert :one
insert into users
    ( username, email, password, created_at, updated_at )
values
    ( $1, $2, $3, $4, $5 )
returning username;

-- name: UserUpdate :one
update users set
    username = $1, email = $2, created_at = $3, updated_at = $4
where id = $5
returning *;

-- name: UserDelete :one
delete from users where id = $1 returning id;
