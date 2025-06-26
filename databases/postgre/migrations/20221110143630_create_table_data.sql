-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.cust (
    id bigserial NOT NULL,
    "name" varchar NULL,
    CONSTRAINT cust_pk PRIMARY KEY (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE public.cust;

-- +goose StatementEnd