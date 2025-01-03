-- +goose Up
CREATE TABLE IF NOT EXISTS teams (
    full_name TEXT NOT NULL,
    short_name TEXT NOT NULL
);

INSERT INTO teams (
    full_name,
    short_name
)
VALUES
('Anaheim Ducks', 'ANA'),
('Boston Bruins', 'BOS'),
('Buffalo Sabres', 'BUF'),
('Calgary Flames', 'CGY'),
('Carolina Hurricanes', 'CAR'),
('Chicago Blackhawks', 'CHI'),
('Colorado Avalanche', 'COL'),
('Columbus Blue Jackets', 'CBJ'),
('Dallas Stars', 'DAL'),
('Detroit Red Wings', 'DET'),
('Edmonton Oilers', 'EDM'),
('Florida Panthers', 'FLA'),
('Los Angeles Kings', 'LAK'),
('Minnesota Wild', 'MIN'),
('Montreal Canadiens', 'MTL'),
('Nashville Predators', 'NSH'),
('New Jersey Devils', 'NJD'),
('New York Islanders', 'NYI'),
('New York Rangers', 'NYR'),
('Ottawa Senators', 'OTT'),
('Philadelphia Flyers', 'PHI'),
('Pittsburgh Penguins', 'PIT'),
('San Jose Sharks', 'SJS'),
('Seattle Kraken', 'SEA'),
('St. Louis Blues', 'STL'),
('Tampa Bay Lightning', 'TBL'),
('Toronto Maple Leafs', 'TOR'),
('Utah Hockey Club', 'UTA'),
('Vancouver Canucks', 'VAN'),
('Vegas Golden Knights', 'VEG'),
('Washington Capitals', 'WSH'),
('Winnipeg Jets', 'WPG');

-- +goose Down
DROP TABLE teams;