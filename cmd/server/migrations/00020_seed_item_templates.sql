-- +goose Up
-- Seed common construction item templates (global, available to all orgs)

-- Lumber - Studs
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Framing', '2x4x8 Stud', 'ea', 4.50),
('material', 'Framing', '2x4x10 Stud', 'ea', 6.00),
('material', 'Framing', '2x4x12 Stud', 'ea', 7.50),
('material', 'Framing', '2x4x16 Stud', 'ea', 10.00),
('material', 'Framing', '2x6x8 Stud', 'ea', 7.00),
('material', 'Framing', '2x6x10 Stud', 'ea', 9.00),
('material', 'Framing', '2x6x12 Stud', 'ea', 11.00),
('material', 'Framing', '2x6x16 Stud', 'ea', 14.00);

-- Lumber - Treated
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Framing', '2x4x8 PT', 'ea', 6.50),
('material', 'Framing', '2x6x8 PT', 'ea', 9.50),
('material', 'Framing', '2x6x12 PT', 'ea', 14.00),
('material', 'Framing', '4x4x8 PT Post', 'ea', 12.00),
('material', 'Framing', '4x4x10 PT Post', 'ea', 15.00),
('material', 'Framing', '6x6x8 PT Post', 'ea', 28.00);

-- Lumber - Beams/Headers
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Framing', '2x8x12', 'ea', 12.00),
('material', 'Framing', '2x8x16', 'ea', 16.00),
('material', 'Framing', '2x10x12', 'ea', 18.00),
('material', 'Framing', '2x10x16', 'ea', 24.00),
('material', 'Framing', '2x12x12', 'ea', 22.00),
('material', 'Framing', '2x12x16', 'ea', 30.00);

-- Sheathing
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Sheathing', '7/16 OSB 4x8', 'ea', 14.00),
('material', 'Sheathing', '1/2 CDX Plywood 4x8', 'ea', 32.00),
('material', 'Sheathing', '3/4 CDX Plywood 4x8', 'ea', 45.00),
('material', 'Sheathing', '1/2 Drywall 4x8', 'ea', 12.00),
('material', 'Sheathing', '5/8 Drywall 4x8', 'ea', 14.00),
('material', 'Sheathing', '1/2 Cement Board 3x5', 'ea', 18.00);

-- Insulation
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Insulation', 'R-13 Fiberglass Batt 15"', 'sf', 0.50),
('material', 'Insulation', 'R-19 Fiberglass Batt 15"', 'sf', 0.70),
('material', 'Insulation', 'R-30 Fiberglass Batt 15"', 'sf', 1.00),
('material', 'Insulation', 'R-38 Blown-in', 'sf', 1.25),
('material', 'Insulation', 'Foam Board 1" 4x8', 'ea', 18.00),
('material', 'Insulation', 'Foam Board 2" 4x8', 'ea', 32.00);

-- Fasteners
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Fasteners', '16d Framing Nails (5lb)', 'box', 12.00),
('material', 'Fasteners', '8d Common Nails (5lb)', 'box', 10.00),
('material', 'Fasteners', '3" Deck Screws (5lb)', 'box', 28.00),
('material', 'Fasteners', '1-5/8" Drywall Screws (5lb)', 'box', 18.00),
('material', 'Fasteners', 'Construction Adhesive (28oz)', 'ea', 6.00),
('material', 'Fasteners', 'Simpson Strong-Tie A34', 'ea', 1.50),
('material', 'Fasteners', 'Simpson Strong-Tie H2.5A', 'ea', 2.00),
('material', 'Fasteners', 'Joist Hanger 2x8', 'ea', 3.50),
('material', 'Fasteners', 'Joist Hanger 2x10', 'ea', 4.00);

-- Concrete
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Concrete', 'Concrete 80lb Bag', 'ea', 6.50),
('material', 'Concrete', 'Quikrete 60lb Bag', 'ea', 5.00),
('material', 'Concrete', 'Ready Mix Concrete', 'cy', 150.00),
('material', 'Concrete', 'Rebar #4 (20ft)', 'ea', 12.00),
('material', 'Concrete', 'Rebar #5 (20ft)', 'ea', 18.00),
('material', 'Concrete', 'Wire Mesh 6x6 5x5', 'ea', 45.00);

-- Roofing
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('material', 'Roofing', 'Architectural Shingles (bundle)', 'ea', 35.00),
('material', 'Roofing', '3-Tab Shingles (bundle)', 'ea', 28.00),
('material', 'Roofing', 'Felt Paper 15# (roll)', 'ea', 22.00),
('material', 'Roofing', 'Synthetic Underlayment (roll)', 'ea', 85.00),
('material', 'Roofing', 'Ice & Water Shield (roll)', 'ea', 95.00),
('material', 'Roofing', 'Ridge Vent (4ft)', 'ea', 12.00),
('material', 'Roofing', 'Drip Edge 10ft', 'ea', 8.00);

-- Labor
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('labor', 'Framing', 'Framing Labor', 'hr', 65.00),
('labor', 'Framing', 'Wall Framing', 'lf', 8.00),
('labor', 'Framing', 'Roof Framing', 'sf', 4.50),
('labor', 'Drywall', 'Drywall Hang', 'sf', 1.50),
('labor', 'Drywall', 'Drywall Finish', 'sf', 1.75),
('labor', 'Roofing', 'Roofing Install', 'sq', 85.00),
('labor', 'Concrete', 'Concrete Pour & Finish', 'sf', 8.00),
('labor', 'General', 'General Labor', 'hr', 45.00),
('labor', 'General', 'Skilled Labor', 'hr', 65.00),
('labor', 'General', 'Demolition', 'hr', 50.00);

-- Equipment
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES
('equipment', 'Rental', 'Skid Steer (day)', 'day', 350.00),
('equipment', 'Rental', 'Excavator Mini (day)', 'day', 400.00),
('equipment', 'Rental', 'Scissor Lift (day)', 'day', 200.00),
('equipment', 'Rental', 'Boom Lift (day)', 'day', 350.00),
('equipment', 'Rental', 'Concrete Mixer (day)', 'day', 75.00),
('equipment', 'Rental', 'Compactor Plate (day)', 'day', 65.00),
('equipment', 'Rental', 'Generator 5000W (day)', 'day', 85.00),
('equipment', 'Delivery', 'Material Delivery', 'ea', 75.00),
('equipment', 'Disposal', 'Dumpster 10yd', 'ea', 450.00),
('equipment', 'Disposal', 'Dumpster 20yd', 'ea', 550.00);

-- +goose Down
DELETE FROM item_templates WHERE org_id IS NULL;
