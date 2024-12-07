-- +migrate Up
CREATE TYPE user_role AS ENUM ('deelnemer', 'begeleider', 'vrijwilliger');
CREATE TYPE distance AS ENUM ('2.5km', '6km', '10km', '15km');
CREATE TYPE support_type AS ENUM ('ja', 'nee', 'anders');

CREATE TABLE IF NOT EXISTS registrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'deelnemer',
    distance distance NOT NULL DEFAULT '2.5km',
    needs_support support_type NOT NULL DEFAULT 'nee',
    support_details TEXT,
    terms_accepted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS registrations;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS distance;
DROP TYPE IF EXISTS support_type; 