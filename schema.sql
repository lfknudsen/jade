CREATE TABLE IF NOT EXISTS folder (
                                      uuid VARCHAR PRIMARY KEY,
                                      name VARCHAR,
                                      parent VARCHAR REFERENCES folder ON DELETE CASCADE NULL
);

CREATE TABLE IF NOT EXISTS file (
                                    uuid VARCHAR PRIMARY KEY,
                                    name VARCHAR,
                                    parent VARCHAR REFERENCES folder ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS child (
                                     parent VARCHAR REFERENCES folder ON DELETE CASCADE,
                                     child VARCHAR REFERENCES file ON DELETE CASCADE,
                                     CONSTRAINT rel PRIMARY KEY (parent, child)
);

CREATE TABLE IF NOT EXISTS subdirectory (
                                            parent VARCHAR REFERENCES folder ON DELETE CASCADE,
                                            child VARCHAR REFERENCES folder ON DELETE CASCADE,
                                            CONSTRAINT rel PRIMARY KEY (parent, child)
);