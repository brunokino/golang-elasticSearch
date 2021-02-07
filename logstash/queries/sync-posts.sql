SELECT
       l.id, l.operation, l.post_id, p.id, p.title, p.body
FROM post_logs l
LEFT JOIN posts p
    ON p.id = l.post_id
WHERE l.id > 0 AND l.created_at < NOW()
ORDER BY l.id;