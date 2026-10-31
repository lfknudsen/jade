package main

const dropTables = `
DROP TABLE IF EXISTS subdirectory;
DROP TABLE IF EXISTS child;
DROP TABLE IF EXISTS file;
DROP TABLE IF EXISTS folder;
`

const createTables = `
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
	CONSTRAINT parent_child_relation PRIMARY KEY (parent, child)
);

CREATE TABLE IF NOT EXISTS subdirectory (
	parent VARCHAR REFERENCES folder ON DELETE CASCADE,
	child VARCHAR REFERENCES folder ON DELETE CASCADE,
	CONSTRAINT parent_subdirectory_relation PRIMARY KEY (parent, child)
);`

const insertFolder = `
INSERT INTO folder VALUES ($1, $2, $3);
`

const insertRootFolder = `
INSERT INTO folder (uuid, name) VALUES ($1, $2);
`
const insertFile = `
INSERT INTO file VALUES ($1, $2, $3);
`

const insertFolderRelation = `
INSERT INTO subdirectory VALUES ($1, $2);
`

const insertChildRelation = `
INSERT INTO child VALUES ($1, $2);
`
