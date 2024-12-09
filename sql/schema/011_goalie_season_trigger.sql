-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_goalie_season_stats()
RETURNS TRIGGER AS $$
DECLARE 
    win INT := 0;
    loss INT := 0;
    otl INT := 0;
BEGIN   
    -- Convert decision text to integer
    CASE NEW.decision
        WHEN 'W' THEN
            win := 1;
        WHEN 'L' THEN
            loss := 1;
        WHEN 'O' THEN
            otl := 1;
        ELSE
    END CASE;

    -- Check if player already exists in goalie_season_stats for this season and team
    IF EXISTS (
        SELECT 1 
        FROM goalie_season_stats
        WHERE player_id = NEW.PLAYER_ID
        AND team = NEW.TEAM
        AND season = NEW.SEASON
    ) 
    THEN
        -- Update the existing record by adding the new game stats to the totals
        UPDATE goalie_season_stats
        SET 
            GAMES_PLAYED = GAMES_PLAYED + 1,
            WINS = WINS + win,
            LOSSES = LOSSES + loss,
            OTLOSSES = OTLOSSES + otl,
            GOALS_AGAINST = GOALS_AGAINST + NEW.GOALS_AGAINST,
            SHOTS_AGAINST = SHOTS_AGAINST + NEW.SHOTS_AGAINST,
            SAVES = SAVES + NEW.SAVES,
            SHUTOUTS = SHUTOUTS + NEW.SHUTOUT, -- cast boolean to int
            PEN_MINS = PEN_MINS + NEW.PEN_MINS,
            TIME_ON_ICE = TIME_ON_ICE + NEW.TIME_ON_ICE,
            UPDATED_AT = NOW()

        WHERE player_id = NEW.PLAYER_ID
        AND team = NEW.TEAM
        AND season = NEW.SEASON;
    ELSE
        -- If no record exists, add new one
        INSERT INTO goalie_season_stats (
            PLAYER_ID,
            PLAYER_NAME,
            TEAM,
            SEASON,
            GAMES_PLAYED,
            WINS, 
            LOSSES,
            OTLOSSES,
            GOALS_AGAINST,
            SHOTS_AGAINST,
            SAVES,
            SHUTOUTS,
            PEN_MINS,
            TIME_ON_ICE,
            UPDATED_AT
        ) VALUES (
            NEW.PLAYER_ID,
            NEW.PLAYER_NAME,
            NEW.TEAM,
            NEW.SEASON,
            1, -- GAMES_PLAYED
            win,
            loss,
            otl,
            NEW.GOALS_AGAINST,
            NEW.SHOTS_AGAINST,
            NEW.SAVES,
            NEW.SHUTOUT::INT,
            NEW.PEN_MINS,
            NEW.TIME_ON_ICE,
            NOW() -- UPDATED_AT
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER update_goalie_season_stats_trigger
AFTER INSERT ON goalie_game_stats
FOR EACH ROW
EXECUTE FUNCTION update_goalie_season_stats();

-- +goose Down
DROP TRIGGER update_goalie_season_stats_trigger ON goalie_game_stats;
DROP FUNCTION update_goalie_season_stats();
