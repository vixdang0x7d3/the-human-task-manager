-- name: CreateMembership :one
INSERT INTO project_memberships(user_id, project_id, role)
VALUES (@user_id, @project_id, @role::membership_role)
RETURNING *;

-- name: AcceptMembership :one
UPDATE project_memberships m
SET role = 'member'::membership_role
WHERE m.user_id = @user_id
AND m.project_id = @project_id
RETURNING *;

-- name: DeleteMembership :one
DELETE FROM project_memberships
WHERE user_id = @user_id
AND project_id = @project_id
RETURNING *;

-- name: MembershipsByUserID :many
SELECT sqlc.embed(u), sqlc.embed(p), role, COUNT(*) OVER()
FROM project_memberships m 
JOIN users u ON m.user_id = u.id
JOIN projects p ON m.project_id = p.id
WHERE m.user_id = @user_id
LIMIT @nlimit
OFFSET @noffset;

-- name: MembershipsByProjectID :many
SELECT sqlc.embed(u), sqlc.embed(p), role, COUNT (*) OVER()
FROM project_memberships m 
JOIN users u ON m.user_id = u.id
JOIN projects p ON m.project_id = p.id
WHERE m.project_id = @project_id
LIMIT @nlimit 
OFFSET @noffset;

-- name: MembershipByIDs :one
SELECT * from project_memberships
WHERE project_id = @project_id
AND user_id = @user_id;

