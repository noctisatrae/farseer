CREATE TABLE casts (
  -- common to all messages 
  id SERIAL PRIMARY KEY,
  fid BIGINT NOT NULL,
  
  -- TIME info
  -- when the message was received by the hub
  timestamp BIGINT NOT NULL,
  -- when was the message created in the DB?
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  -- when was it removed on Farcaster? 
  deleted_at TIMESTAMP,
  -- what was the last update in the DB?
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  -- cast specific stuff
  -- metadata of the cast
  -- I prefer the hash to be text so you can actually see the value in the DB & debug
  hash TEXT,
  parent_hash TEXT,
  parent_fid BIGINT,
  parent_url VARCHAR,

  -- content
  text TEXT,
  embeds JSONB,
  mentions BIGINT[],
  mentions_positions SMALLINT[]
)

CREATE TABLE links (
  -- common to all messages 
  id SERIAL PRIMARY KEY,

  -- TIME info
  -- when the message was received by the hub
  timestamp BIGINT NOT NULL,
  -- when was the message created in the DB?
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  -- when was it removed on Farcaster? 
  deleted_at TIMESTAMP,
  -- what was the last update in the DB?
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  -- specific
  fid BIGINT NOT NULL,
  target_fid BIGINT NOT NULL,
  hash TEXT,
  type TEXT
)

CREATE TABLE reactions (
  -- common to all messages 
  id SERIAL PRIMARY KEY,

  fid BIGINT NOT NULL,
  
  -- TIME info
  -- when the message was received by the hub
  timestamp BIGINT NOT NULL,
  -- when was the message created in the DB?
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  -- when was it removed on Farcaster? 
  deleted_at TIMESTAMP,
  -- what was the last update in the DB?
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  -- specific to reactions 
  reaction_type SMALLINT,
  hash TEXT,
  target_hash TEXT,
  target_fid BIGINT,
  target_url TEXT
)

CREATE TABLE verifications (
  -- common to all messages 
  id SERIAL PRIMARY KEY,

  fid BIGINT NOT NULL,
  
  -- TIME info
  -- when the message was received by the hub
  timestamp BIGINT NOT NULL,
  -- when was the message created in the DB?
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  hash TEXT,
  address TEXT,
  claim_signature TEXT,
  block_hash TEXT,
  verification_type SMALLINT, -- either 0 or 1
  chain_id SMALLINT,
  protocol SMALLINT
)