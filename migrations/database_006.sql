CREATE TABLE cost_items(
    ID UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id),
    amount numeric(12,2) not null,
    short_description varchar(50),
    notes text,
    date date
);

CREATE INDEX idx_cost_items_project_id ON cost_items(project_id);

ALTER TABLE projects
ADD COLUMN amount numeric(12,2) not null default(0.0);

ALTER TABLE projects
ADD COLUMN estimated_cost numeric(12,2) not null default(0.0);