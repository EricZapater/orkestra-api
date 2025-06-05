ALTER TABLE meetings ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
  to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(description, ''))
) STORED;

CREATE INDEX idx_meetings_search_vector ON meetings USING GIN (search_vector);

ALTER TABLE meeting_topics ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
  to_tsvector('simple', coalesce(title, ''))
) STORED;

CREATE INDEX idx_meeting_topics_search_vector ON meeting_topics USING GIN (search_vector);

ALTER TABLE meeting_topic_agreements ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
  to_tsvector('simple', coalesce(title, ''))
) STORED;

CREATE INDEX idx_agreements_search_vector ON meeting_topic_agreements USING GIN (search_vector);


-- View: public.search_index

-- DROP VIEW public.search_index;

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
     JOIN meetings m ON mt.meeting_id = m.id;


ALTER TABLE public.search_index
    OWNER TO ubuntu;

