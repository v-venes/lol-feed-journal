CREATE TABLE
    IF NOT EXISTS players (
        puuid VARCHAR(255) PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        tag VARCHAR(16) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS players_matches (
        id VARCHAR(255) NOT NULL,
        puuid VARCHAR(255) NOT NULL,
        profile_icon INTEGER NOT NULL,
        summoner_name VARCHAR(255) NOT NULL,
        assists INTEGER NOT NULL,
        deaths INTEGER NOT NULL,
        kills INTEGER NOT NULL,
        kda FLOAT NOT NULL,
        champion_name VARCHAR(255) NOT NULL,
        lane VARCHAR(16) NOT NULL,
        role VARCHAR(16) NOT NULL,
        individual_position VARCHAR(16) NOT NULL,
        team_position VARCHAR(16) NOT NULL,
        win BOOLEAN NOT NULL,
        surrender BOOLEAN NOT NULL,
        remake BOOLEAN NOT NULL,
        match_date TIMESTAMP NOT NULL,
        item0 INTEGER NOT NULL,
        item1 INTEGER NOT NULL,
        item2 INTEGER NOT NULL,
        item3 INTEGER NOT NULL,
        item4 INTEGER NOT NULL,
        item5 INTEGER NOT NULL,
        item6 INTEGER NOT NULL,
        PRIMARY KEY (id, puuid),
    );