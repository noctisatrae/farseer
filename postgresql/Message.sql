CREATE TABLE all_messages (
    id SERIAL PRIMARY KEY,
    fid BIGINT NOT NULL,
    timestamp BIGINT NOT NULL,
    network VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    
    -- Fields specific to CastAdd
    cast_add_text TEXT,
    cast_add_parent_cast_id_fid BIGINT,
    cast_add_parent_cast_id_hash BYTEA,
    cast_add_parent_url TEXT,
    cast_add_embeds_deprecated TEXT[],
    cast_add_mentions BIGINT[],
    cast_add_mentions_positions INTEGER[],
    cast_add_embeds JSONB,

    -- Fields specific to CastRemove
    cast_remove_target_hash BYTEA,

    -- Fields specific to Reaction
    reaction_type VARCHAR(50),
    reaction_target_cast_id_fid BIGINT,
    reaction_target_cast_id_hash BYTEA,
    reaction_target_url TEXT,

    -- Fields specific to Verification
    verification_address BYTEA,
    verification_claim_signature BYTEA,
    verification_block_hash BYTEA,
    verification_type INTEGER,
    verification_chain_id INTEGER,
    verification_protocol VARCHAR(50),

    -- Fields specific to Link
    link_type VARCHAR(50),
    link_display_timestamp INTEGER,
    link_target_fid BIGINT,

    -- Fields specific to FrameAction
    frame_action_url BYTEA,
    frame_action_button_index INTEGER,
    frame_action_cast_id_fid BIGINT,
    frame_action_cast_id_hash BYTEA,
    frame_action_input_text BYTEA,
    frame_action_state BYTEA,
    frame_action_transaction_id BYTEA,
    frame_action_address BYTEA,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)