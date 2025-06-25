CREATE TABLE menu_items (
    ID UUID PRIMARY KEY,
    label VARCHAR(255) NOT NULL,
    icon VARCHAR(100),
    route VARCHAR(255),    
    parent_id UUID REFERENCES menu_items(id),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_separator BOOLEAN DEFAULT FALSE
);


CREATE INDEX idx_menu_items_parent_id ON menu_items(parent_id);
CREATE INDEX idx_menu_items_sort_order ON menu_items(sort_order);

CREATE TABLE profiles (
    ID UUID PRIMARY KEY,
    name varchar(100)
);

CREATE TABLE profile_menus(
    ID UUID PRIMARY KEY,
    menu_id UUID not null references menu_items(id),
    profile_id UUID not null references profiles(id)
);

CREATE INDEX idx_profile_menus_menu_id ON profile_menus(menu_id);
CREATE INDEX idx_profile_menus_profile_id ON profile_menus(profile_id);

ALTER TABLE users
ADD COLUMN profile_id UUID null;

CREATE INDEX idx_users_profile_id ON users(profile_id);