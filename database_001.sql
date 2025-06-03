CREATE TABLE groups (
	ID uuid primary key not null,
	name varchar(50),
	created_at timestamptz default now()
);
CREATE INDEX idx_group_name ON groups(name);
CREATE TABLE users(

	ID uuid primary key not null,
	name varchar(50),
	surname varchar(250),
	phone_number varchar(20) not null unique,
	email varchar(250) not null unique,
	username varchar(250) not null,	
	password varchar(250) not null,
	is_verified bool,
	is_active bool default true,
	created_at timestamptz default now(),
	password_changed_at timestamptz default now()
);
CREATE INDEX idx_user_username ON users(username);
CREATE INDEX idx_user_email ON users(email);
CREATE TABLE group_members (
	ID uuid primary key not null,
	group_id uuid references groups(id) on delete restrict,
	user_id uuid references users(id) on delete restrict,
	joined_at timestamptz default now(),
	is_admin bool default false
);
CREATE INDEX idx_group_member_user ON group_members(user_id);
CREATE INDEX idx_group_member_group ON group_members(group_id);
CREATE TABLE meetings (
	ID uuid primary key not null,
	group_id uuid references groups(id) on delete restrict,
	created_at timestamptz default now(),
	start_time timestamptz,	
	created_by uuid references users(id) on delete restrict,
	title varchar(250),
	description text
);
CREATE INDEX idx_meeting_group ON meetings(group_id);
CREATE TABLE meeting_participants (
	ID uuid primary key not null,
	meeting_id uuid references meetings(id) on delete restrict,
	user_id uuid references users(id) on delete restrict
);
CREATE INDEX idx_meeting_participant_user ON meeting_participants(user_id);
CREATE INDEX idx_meeting_participant_meeting ON meeting_participants(meeting_id);
CREATE TABLE meeting_topics (
	ID uuid primary key not null,
	meeting_id uuid references meetings(id) on delete restrict,
	created_at timestamptz default now(),
	title varchar(250)
);
CREATE INDEX idx_meeting_topic_meeting ON meeting_topics(meeting_id);
CREATE TABLE meeting_topic_agreements (
	ID uuid primary key not null,
	meeting_topic_id uuid references meeting_topics(id) on delete restrict,
	created_at timestamptz default now(),
	created_by uuid references users(id) on delete restrict,
	title varchar(250)
);
CREATE INDEX idx_meeting_topic_agreement_topic ON meeting_topic_agreements(meeting_topic_id);
CREATE TABLE labels (
	ID uuid primary key not null,
	name varchar(50) unique,
	created_at timestamptz default now()
);
CREATE INDEX idx_label_name ON labels(name);
CREATE TABLE entity_labels (
	ID uuid primary key not null,
	entity_id uuid,
	entity_name varchar(250) not null,
	label_id uuid references labels(id) on delete restrict
);
CREATE UNIQUE INDEX idx_entity_label ON entity_labels(entity_id, label_id);
CREATE TABLE action_logs (
    id uuid primary key default gen_random_uuid(),
    user_id uuid,
    action_type varchar(100) not null,
    metadata jsonb,
	timezone varchar(100),
    performed_at timestamptz default now()
);
CREATE INDEX idx_action_user ON action_logs(user_id);
CREATE INDEX idx_action_type ON action_logs(action_type);
ALTER TABLE meeting_topic_agreements
ALTER COLUMN title TYPE text;