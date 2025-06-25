CREATE TABLE tasks (
    ID UUID PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    notes TEXT,
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    priority VARCHAR(2) NOT NULL DEFAULT 'D',
    project_id UUID REFERENCES projects(id),
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);

