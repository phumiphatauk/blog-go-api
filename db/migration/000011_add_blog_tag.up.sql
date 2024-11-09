CREATE TABLE blog_tag (
    id BIGSERIAL PRIMARY KEY,
    blog_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE blog_tag
ADD CONSTRAINT fk_blog_tag_blog FOREIGN KEY (blog_id) REFERENCES blog (id);

ALTER TABLE blog_tag
ADD CONSTRAINT fk_blog_tag_tag FOREIGN KEY (tag_id) REFERENCES tag (id);
