-- Update existing categories with images and descriptions

-- Update root categories with images and descriptions
UPDATE categories SET 
    description = '{"de": "Hochwertige Werkzeuge und Maschinen für professionelle Anwender", "en": "High-quality tools and machines for professional users"}'::jsonb,
    image = '/images/categories/werkzeuge.jpg'
WHERE code = 'werkzeuge';

UPDATE categories SET 
    description = '{"de": "Elektrische Komponenten, Kabel und Installationsmaterial", "en": "Electrical components, cables and installation materials"}'::jsonb,
    image = '/images/categories/elektro.jpg'
WHERE code = 'elektro';

UPDATE categories SET 
    description = '{"de": "Schrauben, Muttern, Dübel und Befestigungsmaterial", "en": "Screws, nuts, dowels and fastening materials"}'::jsonb,
    image = '/images/categories/befestigung.jpg'
WHERE code = 'befestigung';

-- Update subcategories (Werkzeuge)
UPDATE categories SET 
    description = '{"de": "Elektrische Bohrmaschinen und Bohrhämmer", "en": "Electric drills and hammer drills"}'::jsonb,
    image = '/images/categories/bohrmaschinen.jpg'
WHERE code = 'bohrmaschinen';

UPDATE categories SET 
    description = '{"de": "Winkelschleifer und Trennschleifer für vielfältige Anwendungen", "en": "Angle grinders and cut-off grinders for various applications"}'::jsonb,
    image = '/images/categories/schleifer.jpg'
WHERE code = 'schleifer';

UPDATE categories SET 
    description = '{"de": "Akkuschrauber und Schlagbohrschrauber", "en": "Cordless drills and impact drivers"}'::jsonb,
    image = '/images/categories/akkuschrauber.jpg'
WHERE code = 'akkuschrauber';

-- Update subcategories (Elektro)
UPDATE categories SET 
    description = '{"de": "Stromkabel, Datenkabel und Leitungen", "en": "Power cables, data cables and wires"}'::jsonb,
    image = '/images/categories/kabel.jpg'
WHERE code = 'kabel';

UPDATE categories SET 
    description = '{"de": "Schalter, Steckdosen und Dimmer", "en": "Switches, sockets and dimmers"}'::jsonb,
    image = '/images/categories/schalter.jpg'
WHERE code = 'schalter';

UPDATE categories SET 
    description = '{"de": "LED-Leuchtmittel, Halogen und Glühbirnen", "en": "LED bulbs, halogen and light bulbs"}'::jsonb,
    image = '/images/categories/leuchtmittel.jpg'
WHERE code = 'leuchtmittel';

-- Update subcategories (Befestigung)
UPDATE categories SET 
    description = '{"de": "Schrauben in verschiedenen Größen und Ausführungen", "en": "Screws in various sizes and designs"}'::jsonb,
    image = '/images/categories/schrauben.jpg'
WHERE code = 'schrauben';

UPDATE categories SET 
    description = '{"de": "Dübel für verschiedene Untergründe", "en": "Dowels for various surfaces"}'::jsonb,
    image = '/images/categories/duebel.jpg'
WHERE code = 'duebel';

UPDATE categories SET 
    description = '{"de": "Muttern, Unterlegscheiben und Sicherungen", "en": "Nuts, washers and lock washers"}'::jsonb,
    image = '/images/categories/muttern.jpg'
WHERE code = 'muttern';

COMMENT ON COLUMN categories.description IS 'Multi-language description updated with seed data';
COMMENT ON COLUMN categories.image IS 'Category image URLs updated with seed data';
