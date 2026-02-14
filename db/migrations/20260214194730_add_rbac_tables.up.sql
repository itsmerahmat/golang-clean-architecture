-- Create roles table
CREATE TABLE roles (
    id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB;

-- Create permissions table
CREATE TABLE permissions (
    id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY unique_permission (resource, action)
) ENGINE = InnoDB;

-- Create user_roles junction table
CREATE TABLE user_roles (
    user_id VARCHAR(100) NOT NULL,
    role_id VARCHAR(100) NOT NULL,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE = InnoDB;

-- Create role_permissions junction table
CREATE TABLE role_permissions (
    role_id VARCHAR(100) NOT NULL,
    permission_id VARCHAR(100) NOT NULL,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE = InnoDB;

-- Insert default roles
INSERT INTO roles (id, name, description, created_at, updated_at)
VALUES 
    ('role-admin', 'admin', 'Administrator with full access', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('role-user', 'user', 'Regular user with limited access', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000);

-- Insert default permissions
INSERT INTO permissions (id, name, resource, action, description, created_at, updated_at)
VALUES
    ('perm-user-read', 'user:read', 'user', 'read', 'Read user information', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-user-write', 'user:write', 'user', 'write', 'Create and update users', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-user-delete', 'user:delete', 'user', 'delete', 'Delete users', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-contact-read', 'contact:read', 'contact', 'read', 'Read contacts', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-contact-write', 'contact:write', 'contact', 'write', 'Create and update contacts', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-contact-delete', 'contact:delete', 'contact', 'delete', 'Delete contacts', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-address-read', 'address:read', 'address', 'read', 'Read addresses', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-address-write', 'address:write', 'address', 'write', 'Create and update addresses', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-address-delete', 'address:delete', 'address', 'delete', 'Delete addresses', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-role-read', 'role:read', 'role', 'read', 'Read roles', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-role-write', 'role:write', 'role', 'write', 'Create and update roles', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
    ('perm-role-delete', 'role:delete', 'role', 'delete', 'Delete roles', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000);

-- Assign all permissions to admin role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT 'role-admin', id, UNIX_TIMESTAMP() * 1000
FROM permissions;

-- Assign user permissions to user role (read/write for own resources)
INSERT INTO role_permissions (role_id, permission_id, created_at)
VALUES
    ('role-user', 'perm-user-read', UNIX_TIMESTAMP() * 1000),
    ('role-user', 'perm-contact-read', UNIX_TIMESTAMP() * 1000),
    ('role-user', 'perm-contact-write', UNIX_TIMESTAMP() * 1000),
    ('role-user', 'perm-contact-delete', UNIX_TIMESTAMP() * 1000),
    ('role-user', 'perm-address-read', UNIX_TIMESTAMP() * 1000),
    ('role-user', 'perm-address-write', UNIX_TIMESTAMP() * 1000),
    ('role-user', 'perm-address-delete', UNIX_TIMESTAMP() * 1000);
