-- projects
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- authors
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL UNIQUE
);

-- issues
CREATE TABLE issues (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    project_key VARCHAR(255) NOT NULL REFERENCES projects(key) ON DELETE CASCADE ON UPDATE CASCADE,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,
    closed TIMESTAMP,
    summary VARCHAR(255) NOT NULL,
    description TEXT,
    issue_type VARCHAR(255) NOT NULL,
    priority VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    time_spent INT,
    creator_id INT NOT NULL REFERENCES authors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    assignee_id INT REFERENCES authors(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- status_changes
CREATE TABLE status_changes (
    id SERIAL PRIMARY KEY,
    issue_id VARCHAR(255) NOT NULL REFERENCES issues(key) ON DELETE CASCADE ON UPDATE CASCADE,
    author_id INT REFERENCES authors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created TIMESTAMP NOT NULL,
    from_status VARCHAR(255),
    to_status VARCHAR(255),
    UNIQUE (issue_id, created)
);

-- analytics
CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    project_key VARCHAR(255) REFERENCES projects(key) ON DELETE CASCADE ON UPDATE CASCADE,
    task_number INT NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_key, task_number)
);

-- Индексы для оптимизации
CREATE INDEX idx_issues_project_key ON issues(project_key);
CREATE INDEX idx_issues_status ON issues(status);
CREATE INDEX idx_issues_created ON issues(created);
CREATE INDEX idx_status_changes_issue_id ON status_changes(issue_id);
CREATE INDEX idx_status_changes_created ON status_changes(created);

-- Триггеры для логирования изменений
CREATE OR REPLACE FUNCTION log_issue_changes()
    RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO status_changes(issue_id, author_id, created, from_status, to_status)
        VALUES(NEW.key, NEW.creator_id, NEW.created, NULL, NEW.status);
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        IF OLD.status != NEW.status THEN
            INSERT INTO status_changes(issue_id, author_id, created, from_status, to_status)
            VALUES(OLD.key, NEW.creator_id, NOW(), OLD.status, NEW.status);
        END IF;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER issue_change_trigger
    AFTER INSERT OR UPDATE ON issues
    FOR EACH ROW EXECUTE FUNCTION log_issue_changes();