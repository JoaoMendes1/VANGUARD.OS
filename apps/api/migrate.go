package main

import (
	"database/sql"
	"fmt"
	"log"
)

func RunMigrations(db *sql.DB) {
	query := `
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		alias VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		current_level INT DEFAULT 1,
		current_xp INT DEFAULT 0,
		credits INT DEFAULT 0,
		designation VARCHAR(100) DEFAULT 'INITIATE',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS hobbies (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(50) NOT NULL,
		level INT DEFAULT 1,
		current_xp INT DEFAULT 0,
		next_level_xp INT DEFAULT 1000,
		icon_name VARCHAR(50),
		UNIQUE(user_id, name)
	);

	CREATE TABLE IF NOT EXISTS protocols (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		streak_count INT DEFAULT 0,
		attribute_type VARCHAR(20) CHECK (attribute_type IN ('Kinetic', 'Neural', 'Core', 'Sync', 'Logic')),
		last_completed_at DATE,
		is_active BOOLEAN DEFAULT TRUE
	);

	CREATE TABLE IF NOT EXISTS operations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		hobby_id UUID REFERENCES hobbies(id) ON DELETE SET NULL,
		title VARCHAR(255) NOT NULL,
		priority VARCHAR(10) CHECK (priority IN ('Low', 'Medium', 'High')),
		xp_reward INT DEFAULT 0,
		credit_reward INT DEFAULT 0,
		is_completed BOOLEAN DEFAULT FALSE,
		deadline TIMESTAMP WITH TIME ZONE,
		completed_at TIMESTAMP WITH TIME ZONE
	);

	CREATE TABLE IF NOT EXISTS ledger_entries (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		amount INT NOT NULL,
		entry_type VARCHAR(10) CHECK (entry_type IN ('GAINED', 'DECAY', 'SPENT')),
		description TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sprints (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		success_rate DECIMAL(5,2) DEFAULT 0.00
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("❌ Failed to run database migrations: ", err)
	}

	fmt.Println("🏗️  [VANGUARD.OS] Database tables created/verified successfully!")
}