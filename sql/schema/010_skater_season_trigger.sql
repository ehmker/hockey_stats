-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_skater_season_stats() 
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the player already has an entry in skater_season_stats for this season and team
    IF EXISTS (
        SELECT 1
        FROM skater_season_stats
        WHERE player_id = NEW.player_id 
        AND team = NEW.team
        AND season = NEW.season
    ) THEN
        -- Update the existing record by adding the NEW game stats to the totals
        UPDATE skater_season_stats
        SET 
            games_played = games_played + 1,
            goals = goals + NEW.goals,
            assists = assists + NEW.assists,
            points = points + NEW.points,
            plus_minus = plus_minus + NEW.plus_minus,
            pen_mins = pen_mins + NEW.pen_mins,
            goals_ev = goals_ev + NEW.goals_ev,
            goals_pp = goals_pp + NEW.goals_pp,
            goals_sh = goals_sh + NEW.goals_sh,
            goals_gw = goals_gw + NEW.goals_gw,
            assists_ev = assists_ev + NEW.assists_ev,
            assists_pp = assists_pp + NEW.assists_pp,
            assists_sh = assists_sh + NEW.assists_sh,
            shots = shots + NEW.shots,
            shifts = shifts + NEW.shifts,
            time_on_ice = time_on_ice + NEW.time_on_ice,
            updated_at = now() -- update the timestamp to the current time
        WHERE 
            player_id = NEW.player_id 
            AND team = NEW.team
            AND season = NEW.season;
    ELSE
        -- If no existing record, insert a NEW one
        INSERT INTO skater_season_stats (
            player_id, 
            player_name,
            team, 
            season, 
            games_played,
            goals,
            assists, 
            points, 
            plus_minus,
            pen_mins,
            goals_ev, 
            goals_pp,
            goals_sh, 
            goals_gw,
            assists_ev,
            assists_pp,
            assists_sh,
            shots,
            shifts,
            time_on_ice,
            updated_at
        ) VALUES (
            NEW.player_id,
            NEW.player_name,
            NEW.team,
            NEW.season,
            1, -- games_played
            NEW.goals,
            NEW.assists,
            NEW.points,
            NEW.plus_minus,
            NEW.pen_mins,
            NEW.goals_ev,
            NEW.goals_pp,
            NEW.goals_sh,
            NEW.goals_gw,
            NEW.assists_ev,
            NEW.assists_pp,
            NEW.assists_sh,
            NEW.shots,
            NEW.shifts,
            NEW.time_on_ice,
            now() -- updated_at
        );
        
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER update_skater_stats_trigger
AFTER INSERT ON skater_game_stats
FOR EACH ROW
EXECUTE FUNCTION update_skater_season_stats();

-- +goose Down
DROP TRIGGER update_skater_stats_trigger ON skater_game_stats;
DROP FUNCTION update_skater_season_stats();