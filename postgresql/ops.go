package main

const (
	// A query that allows to see if a cast exist in the DB
	CastCheck = `
	SELECT created_at FROM casts WHERE hash = $1 
	`

	// A query that allows to add a cast into the DB
	CastAdd = `
	INSERT INTO casts (
		fid,

		timestamp,

		hash,
		parent_hash,
		parent_fid,
		parent_url,

		text,
		embeds,
		mentions,
		mentions_positions
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	CastAddRemoved = `
	INSERT INTO casts (
		fid, 

		timestamp,
		deleted_at,

		hash
	) VALUES ($1, $2, $3, $4)
	`

	// A query that updates a cast when it's removed
	UpdateCastOnRemove = `
	UPDATE casts
	SET
		deleted_at = $1
		updated_at = CURRENT_TIMESTAMP
	WHERE
		hash = $2
	`

	// A query to add a link into the database
	LinkAdd = `
	INSERT INTO links (
		timestamp,

		fid,
		target_fid,
		hash,
		type
	) VALUES ($1, $2, $3, $4, $5)
	`

	LinkRemove = `
	UPDATE links
	SET
		deleted_at = $1
		updated_at = CURRENT_TIMESTAMP
	WHERE
		target_fid = $2
	`

	ReactionAdd = `
	INSERT INTO reactions (
		fid,

		timestamp,
		
		reaction_type,
		hash,
		target_hash,
		target_fid,
		target_url
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	VerificationAdd = `
	INSERT INTO verifications (
		fid,

		timestamp,

		hash,
		address,
		claim_signature,
		block_hash,
		verification_type,
		chain_id,
		protocol
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
)
