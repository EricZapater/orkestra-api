CREATE TABLE customers(
    ID uuid primary key not null,
    comercial_name varchar(150) not null,
    vat_number varchar(15) not null,
    phone_number varchar(25) not null
);

CREATE INDEX idx_customer_name ON customers(comercial_name);
CREATE INDEX idx_customer_vat_number ON customers(vat_number);
CREATE INDEX idx_customer_phone ON customers(phone_number);

CREATE TABLE customer_users(
    ID uuid primary key not null,
    customer_id uuid not null references customers(id) on delete restrict,
    user_id uuid not null references users(id) on delete restrict
);

CREATE INDEX idx_customer_user ON customer_users(customer_id, user_id);


CREATE TABLE projects(
    ID uuid primary key not null,
    description text,
    start_date timestamp,
    end_date timestamp,
    color varchar(25),
    customer_id uuid not null references customers(id) on delete restrict
);



ALTER TABLE projects ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
  to_tsvector('simple', coalesce(description, '') )
) STORED;

CREATE INDEX idx_projects_search_vector ON projects USING GIN (search_vector);

DROP VIEW public.search_index
CREATE OR REPLACE VIEW public.search_index
 AS
 SELECT 'Reunió'::text AS tipus,
    m.id,
    m.group_id,
    (COALESCE(m.title, ''::character varying)::text || ' '::text) || COALESCE(m.description, ''::text) AS text,
    'meeting/'::text || m.id AS link,
    m.search_vector
   FROM meetings m
UNION ALL
 SELECT 'Tema'::text AS tipus,
    t.id,
    m.group_id,
    t.title AS text,
    'meeting/'::text || t.meeting_id AS link,
    t.search_vector
   FROM meeting_topics t
     JOIN meetings m ON t.meeting_id = m.id
UNION ALL
 SELECT 'Acord'::text AS tipus,
    a.id,
    m.group_id,
    a.title AS text,
    'meeting/'::text || mt.meeting_id AS link,
    a.search_vector
   FROM meeting_topic_agreements a
     JOIN meeting_topics mt ON mt.id = a.meeting_topic_id
     JOIN meetings m ON mt.meeting_id = m.id
UNION ALL
  SELECT 'Projecte'::text as tipus,
   id,
   null as group_id,
   description as text,
   'project/'::text || id as link,
   search_vector
   FROM projects;

ALTER TABLE public.search_index
    OWNER TO ubuntu;