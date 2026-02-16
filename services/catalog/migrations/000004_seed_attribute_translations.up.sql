-- Seed common attribute translations (German)
-- Note: Replace with actual tenant_id from your system
-- This assumes a tenant_id exists, adjust as needed

-- Common technical attributes
INSERT INTO attribute_translations (tenant_id, attribute_key, locale, display_name, unit) VALUES
-- Dimensions
('00000000-0000-0000-0000-000000000001', 'thickness_mm', 'de', 'Dicke', 'mm'),
('00000000-0000-0000-0000-000000000001', 'width_mm', 'de', 'Breite', 'mm'),
('00000000-0000-0000-0000-000000000001', 'length_mm', 'de', 'Länge', 'mm'),
('00000000-0000-0000-0000-000000000001', 'height_mm', 'de', 'Höhe', 'mm'),
('00000000-0000-0000-0000-000000000001', 'diameter_mm', 'de', 'Durchmesser', 'mm'),

-- Weight
('00000000-0000-0000-0000-000000000001', 'weight_kg', 'de', 'Gewicht', 'kg'),
('00000000-0000-0000-0000-000000000001', 'weight_g', 'de', 'Gewicht', 'g'),

-- Electrical
('00000000-0000-0000-0000-000000000001', 'voltage', 'de', 'Spannung', 'V'),
('00000000-0000-0000-0000-000000000001', 'voltage_v', 'de', 'Spannung', 'V'),
('00000000-0000-0000-0000-000000000001', 'current_a', 'de', 'Stromstärke', 'A'),
('00000000-0000-0000-0000-000000000001', 'power_w', 'de', 'Leistung', 'W'),
('00000000-0000-0000-0000-000000000001', 'frequency_hz', 'de', 'Frequenz', 'Hz'),

-- Pressure
('00000000-0000-0000-0000-000000000001', 'max_pressure', 'de', 'Max. Druck', 'bar'),
('00000000-0000-0000-0000-000000000001', 'pressure_bar', 'de', 'Druck', 'bar'),

-- Temperature
('00000000-0000-0000-0000-000000000001', 'max_temperature', 'de', 'Max. Temperatur', '°C'),
('00000000-0000-0000-0000-000000000001', 'min_temperature', 'de', 'Min. Temperatur', '°C'),
('00000000-0000-0000-0000-000000000001', 'temperature_c', 'de', 'Temperatur', '°C'),

-- Material
('00000000-0000-0000-0000-000000000001', 'material', 'de', 'Material', NULL),
('00000000-0000-0000-0000-000000000001', 'material_type', 'de', 'Materialtyp', NULL),
('00000000-0000-0000-0000-000000000001', 'finish', 'de', 'Oberfläche', NULL),

-- Color
('00000000-0000-0000-0000-000000000001', 'color', 'de', 'Farbe', NULL),
('00000000-0000-0000-0000-000000000001', 'colour', 'de', 'Farbe', NULL),

-- Capacity
('00000000-0000-0000-0000-000000000001', 'capacity_l', 'de', 'Kapazität', 'l'),
('00000000-0000-0000-0000-000000000001', 'capacity_ml', 'de', 'Kapazität', 'ml'),

-- Common properties
('00000000-0000-0000-0000-000000000001', 'brand', 'de', 'Marke', NULL),
('00000000-0000-0000-0000-000000000001', 'model', 'de', 'Modell', NULL),
('00000000-0000-0000-0000-000000000001', 'manufacturer', 'de', 'Hersteller', NULL)

ON CONFLICT (tenant_id, attribute_key, locale) DO NOTHING;

-- English translations
INSERT INTO attribute_translations (tenant_id, attribute_key, locale, display_name, unit) VALUES
-- Dimensions
('00000000-0000-0000-0000-000000000001', 'thickness_mm', 'en', 'Thickness', 'mm'),
('00000000-0000-0000-0000-000000000001', 'width_mm', 'en', 'Width', 'mm'),
('00000000-0000-0000-0000-000000000001', 'length_mm', 'en', 'Length', 'mm'),
('00000000-0000-0000-0000-000000000001', 'height_mm', 'en', 'Height', 'mm'),
('00000000-0000-0000-0000-000000000001', 'diameter_mm', 'en', 'Diameter', 'mm'),

-- Weight
('00000000-0000-0000-0000-000000000001', 'weight_kg', 'en', 'Weight', 'kg'),
('00000000-0000-0000-0000-000000000001', 'weight_g', 'en', 'Weight', 'g'),

-- Electrical
('00000000-0000-0000-0000-000000000001', 'voltage', 'en', 'Voltage', 'V'),
('00000000-0000-0000-0000-000000000001', 'voltage_v', 'en', 'Voltage', 'V'),
('00000000-0000-0000-0000-000000000001', 'current_a', 'en', 'Current', 'A'),
('00000000-0000-0000-0000-000000000001', 'power_w', 'en', 'Power', 'W'),
('00000000-0000-0000-0000-000000000001', 'frequency_hz', 'en', 'Frequency', 'Hz'),

-- Pressure
('00000000-0000-0000-0000-000000000001', 'max_pressure', 'en', 'Max. Pressure', 'bar'),
('00000000-0000-0000-0000-000000000001', 'pressure_bar', 'en', 'Pressure', 'bar'),

-- Temperature
('00000000-0000-0000-0000-000000000001', 'max_temperature', 'en', 'Max. Temperature', '°C'),
('00000000-0000-0000-0000-000000000001', 'min_temperature', 'en', 'Min. Temperature', '°C'),
('00000000-0000-0000-0000-000000000001', 'temperature_c', 'en', 'Temperature', '°C'),

-- Material
('00000000-0000-0000-0000-000000000001', 'material', 'en', 'Material', NULL),
('00000000-0000-0000-0000-000000000001', 'material_type', 'en', 'Material Type', NULL),
('00000000-0000-0000-0000-000000000001', 'finish', 'en', 'Finish', NULL),

-- Color
('00000000-0000-0000-0000-000000000001', 'color', 'en', 'Color', NULL),
('00000000-0000-0000-0000-000000000001', 'colour', 'en', 'Colour', NULL),

-- Capacity
('00000000-0000-0000-0000-000000000001', 'capacity_l', 'en', 'Capacity', 'l'),
('00000000-0000-0000-0000-000000000001', 'capacity_ml', 'en', 'Capacity', 'ml'),

-- Common properties
('00000000-0000-0000-0000-000000000001', 'brand', 'en', 'Brand', NULL),
('00000000-0000-0000-0000-000000000001', 'model', 'en', 'Model', NULL),
('00000000-0000-0000-0000-000000000001', 'manufacturer', 'en', 'Manufacturer', NULL)

ON CONFLICT (tenant_id, attribute_key, locale) DO NOTHING;
