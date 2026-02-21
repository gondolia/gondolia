-- Remove images and descriptions from seeded categories

UPDATE categories SET 
    description = '{}'::jsonb,
    image = NULL
WHERE code IN (
    'werkzeuge', 'elektro', 'befestigung',
    'bohrmaschinen', 'schleifer', 'akkuschrauber',
    'kabel', 'schalter', 'leuchtmittel',
    'schrauben', 'duebel', 'muttern'
);
