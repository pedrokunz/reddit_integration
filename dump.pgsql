INSERT INTO customers (id, name, created_at, updated_at, integrations)
VALUES ('afacf9fa-7516-468f-b048-ac4c0562aa3f', 'Customer 1', now(), now(), '{"reddit"}');

INSERT INTO users (email, password, created_at, updated_at, customer_id, roles, attributes)
VALUES ('user1@test.com', '123', now(), now(), 'afacf9fa-7516-468f-b048-ac4c0562aa3f', '{"admin","member","reader"}', '{"execute.reddit.integration": "true", "reddit.post.sync": "true", "reddit.post.read_list": "true"}');