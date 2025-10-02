-- Create summoner table to store Riot API data
CREATE TABLE IF NOT EXISTS summoners (
    id SERIAL PRIMARY KEY,
    puuid VARCHAR(255) UNIQUE NOT NULL,
    summoner_name VARCHAR(255) NOT NULL,
    summoner_level INTEGER NOT NULL,
    profile_icon_id INTEGER NOT NULL,
    region VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create match table to store match data
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    match_id VARCHAR(255) UNIQUE NOT NULL,
    game_mode VARCHAR(50),
    game_duration INTEGER,
    game_creation BIGINT,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_summoners_puuid ON summoners(puuid);
CREATE INDEX IF NOT EXISTS idx_summoners_region ON summoners(region);
CREATE INDEX IF NOT EXISTS idx_matches_match_id ON matches(match_id);
CREATE INDEX IF NOT EXISTS idx_matches_game_creation ON matches(game_creation);
