CREATE TABLE actors (
	id CHAR(255) NOT NULL,
	actor_type CHAR(255) NOT NULL,
	public_key text NOT NULL,
	private_key text NOT NULL,
	PRIMARY KEY (id)
);
