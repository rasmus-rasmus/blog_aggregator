-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.id, 
       posts.created_at, 
       posts.updated_at, 
       posts.title, 
       posts.url, 
       posts.description, 
       posts.published_at, 
       posts.feed_id 
FROM feed_follows 
JOIN posts
ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY published_at DESC
LIMIT $2;