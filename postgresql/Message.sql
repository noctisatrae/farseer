CREATE TABLE all_messages (
    id SERIAL PRIMARY KEY,
    fid BIGINT NOT NULL,
    timestamp BIGINT NOT NULL,
    network INTEGER NOT NULL,
    type INTEGER NOT NULL,
    
    -- Field common to all casts
    hash TEXT,

    -- Fields specific to CastAdd
    cast_add_text TEXT,
    cast_add_parent_cast_id_fid BIGINT,
    cast_add_parent_cast_id_hash TEXT,
    cast_add_parent_url TEXT,
    cast_add_embeds_deprecated TEXT[],
    cast_add_mentions BIGINT[],
    cast_add_mentions_positions INTEGER[],
    cast_add_embeds JSONB,

    -- Fields specific to Reaction
    reaction_type VARCHAR(50),
    reaction_target_cast_id_fid BIGINT,
    reaction_target_cast_id_hash TEXT,
    reaction_target_url TEXT,

    -- Fields specific to Verification
    verification_address TEXT,
    verification_claim_signature TEXT,
    verification_block_hash TEXT,
    verification_type INTEGER,
    verification_chain_id INTEGER,
    verification_protocol VARCHAR(50),

    -- Fields specific to Link
    link_type VARCHAR(50),
    link_display_timestamp INTEGER,
    link_target_fid BIGINT,

    -- Fields specific to FrameAction
    frame_action_url TEXT,
    frame_action_button_index INTEGER,
    frame_action_cast_id_fid BIGINT,
    frame_action_cast_id_hash TEXT,
    frame_action_input_text TEXT,
    frame_action_state TEXT,
    frame_action_transaction_id TEXT,
    frame_action_address TEXT,

    removed_at TIMESTAMP
)