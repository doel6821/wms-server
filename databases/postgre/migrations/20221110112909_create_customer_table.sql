-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.customers (
	id bigserial NOT NULL,
	"name" varchar NULL,
	phone varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT customers_pk PRIMARY KEY (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE public.customers;

-- +goose StatementEnd