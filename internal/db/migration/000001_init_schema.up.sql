CREATE TABLE "employees" (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER'
);

CREATE TABLE "projects" (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    created_by INT REFERENCES employees(id) ON DELETE SET NULL
);

CREATE TABLE "tasks" (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    due_to TIMESTAMPTZ NOT NULL,
    time_spent BIGINT DEFAULT 0,
    status VARCHAR(50) NOT NULL CHECK (status IN ('ASSIGNED', 'STARTED', 'SUSPENDED', 'COMPLETED')) NOT NULL,
    priority VARCHAR(50) CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')) NOT NULL,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    assigned_to INT REFERENCES employees(id) ON DELETE SET NULL
);

CREATE TABLE "histories" (
    id BIGSERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    changed_by INT REFERENCES employees(id) ON DELETE SET NULL,
    change_at TIMESTAMPTZ DEFAULT (now()) NOT NULL,
    old_status VARCHAR(50) NOT NULL,
    new_status VARCHAR(50) NOT NULL
);

CREATE TABLE project_invitations (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(64) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP,
    UNIQUE (project_id, email)
);


CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);

CREATE INDEX idx_employees_email ON employees(email);

CREATE INDEX idx_history_task_id ON histories(task_id);

CREATE INDEX idx_project_created_by ON projects(created_by);

CREATE INDEX idx_project_invitations_token ON project_invitations(token);
CREATE INDEX idx_project_invitations_email ON project_invitations(email);