CREATE TABLE operators(
    ID UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(250) NOT NULL,
    cost decimal(10,2)
);

ALTER TABLE users 
ADD COLUMN is_customer BOOLEAN default(false);

CREATE TABLE operators_to_projects (
    id UUID PRIMARY KEY,
    operator_id UUID NOT NULL references operators(id) ON DELETE RESTRICT,
    project_id UUID NOT NULL references projects(id) ON DELETE RESTRICT,
    cost DECIMAL(10,2) NOT NULL,
    dedication_percent DECIMAL(5,2) NOT NULL CHECK (dedication_percent >= 0 AND dedication_percent <= 100),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
