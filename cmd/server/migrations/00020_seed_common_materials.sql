-- +goose Up
-- Pre-seed common construction materials with size variants
-- Prices are approximate defaults; clients should update with their vendor pricing

-- ============================================================
-- FASTENERS - Nails
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', '16d Common Nails 5lb Box', 'box', 12.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', '16d Sinker Nails 5lb Box', 'box', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', '10d Common Nails 5lb Box', 'box', 11.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', '8d Common Nails 5lb Box', 'box', 11.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', '6d Common Nails 1lb Box', 'box', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Roofing Nails 1-1/4" 5lb Box', 'box', 13.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Roofing Nails 1-3/4" 5lb Box', 'box', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 4d 1lb Box', 'box', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 6d 1lb Box', 'box', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 8d 1lb Box', 'box', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Brad Nails 18ga 1" 1000ct', 'box', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Brad Nails 18ga 1-1/4" 1000ct', 'box', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Brad Nails 18ga 2" 1000ct', 'box', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Framing Nails 3-1/4" Full Round Head Strip', 'box', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Framing Nails 2-3/8" Full Round Head Strip', 'box', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Framing Nails 3-1/2" Coil', 'box', 50.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Joist Hanger Nails 1-1/2" 10d 5lb', 'box', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Ring Shank Siding Nails 2-1/2" 5lb Box', 'box', 15.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Galvanized Box Nails 8d 5lb Box', 'box', 14.00);

-- ============================================================
-- FASTENERS - Screws
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Drywall Screws 1-1/4" Coarse 1lb Box', 'box', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Drywall Screws 1-5/8" Coarse 1lb Box', 'box', 8.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Drywall Screws 2" Coarse 1lb Box', 'box', 9.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Drywall Screws 3" Coarse 1lb Box', 'box', 11.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Deck Screws #8 2" 1lb Box', 'box', 9.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Deck Screws #8 2-1/2" 1lb Box', 'box', 9.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Deck Screws #8 3" 1lb Box', 'box', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Deck Screws #10 3-1/2" 1lb Box', 'box', 11.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Structural Screws (GRK) 3-1/8" 100ct', 'box', 32.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Structural Screws (GRK) 4" 50ct', 'box', 28.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Structural Screws (GRK) 5-1/8" 50ct', 'box', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Structural Screws (SPAX) 3" 75ct', 'box', 28.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Structural Screws (SPAX) 4" 50ct', 'box', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Cabinet Screws 2-1/2" 100ct', 'box', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Pocket Hole Screws 1-1/4" 100ct', 'box', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Pocket Hole Screws 2-1/2" 100ct', 'box', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Concrete Screws (Tapcon) 1/4"x1-3/4" 75ct', 'box', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Concrete Screws (Tapcon) 1/4"x2-3/4" 75ct', 'box', 20.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Concrete Screws (Tapcon) 1/4"x3-1/4" 75ct', 'box', 22.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Cement Board Screws 1-5/8" 200ct', 'box', 14.00);

-- ============================================================
-- FASTENERS - Bolts & Anchors
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 3/8"x3"', 'ea', 0.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 3/8"x4"', 'ea', 0.90);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 3/8"x5"', 'ea', 1.10);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 3/8"x6"', 'ea', 1.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 1/2"x4"', 'ea', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 1/2"x6"', 'ea', 1.85);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 1/2"x8"', 'ea', 2.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Lag Bolt 1/2"x10"', 'ea', 2.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Carriage Bolt 3/8"x3"', 'ea', 0.65);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Carriage Bolt 3/8"x4"', 'ea', 0.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Carriage Bolt 3/8"x6"', 'ea', 1.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Carriage Bolt 1/2"x4"', 'ea', 1.20);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Carriage Bolt 1/2"x6"', 'ea', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Wedge Anchor 3/8"x3"', 'ea', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Wedge Anchor 1/2"x4-1/2"', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Wedge Anchor 5/8"x5"', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'J-Bolt 1/2"x8"', 'ea', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'J-Bolt 1/2"x10"', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'J-Bolt 5/8"x10"', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Hex Bolt 3/8"x3" w/ Nut & Washer', 'ea', 0.80);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Hex Bolt 1/2"x4" w/ Nut & Washer', 'ea', 1.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Fender Washers 3/8" 100ct', 'box', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Fender Washers 1/2" 100ct', 'box', 10.00);

-- ============================================================
-- FASTENERS - Staples & Pins
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'T50 Staples 3/8" 1000ct', 'box', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'T50 Staples 1/2" 1000ct', 'box', 7.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'T50 Staples 9/16" 1000ct', 'box', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Crown Staples 1/2" 16ga 10000ct', 'box', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Crown Staples 1-1/2" 16ga 10000ct', 'box', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Pin Nails 23ga 1" 2000ct', 'box', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Pin Nails 23ga 1-3/8" 2000ct', 'box', 9.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 16ga 1-1/4" 2500ct', 'box', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 16ga 2" 2500ct', 'box', 27.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 16ga 2-1/2" 2500ct', 'box', 29.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 15ga 1-1/4" 1000ct', 'box', 22.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 15ga 2" 1000ct', 'box', 24.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Fasteners', 'Finish Nails 15ga 2-1/2" 1000ct', 'box', 26.00);

-- ============================================================
-- HARDWARE - Simpson Strong-Tie Connectors
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'Hurricane Tie H2.5', 'ea', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'Hurricane Tie H10', 'ea', 2.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'Rafter Tie RT2', 'ea', 1.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'A34 Framing Angle', 'ea', 1.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'A35 Framing Angle', 'ea', 1.40);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'L50 Angle', 'ea', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'L70 Angle', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'L90 Angle', 'ea', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'TP35 Tie Plate 3x5', 'ea', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'TP37 Tie Plate 3x7', 'ea', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'MP24 Mending Plate 2x4', 'ea', 1.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'MP36 Mending Plate 3x6', 'ea', 1.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'LSTA 1-1/4"x18" Strap', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'LSTA 1-1/4"x24" Strap', 'ea', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'MSTA 1-3/8"x30" Strap', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'MSTA 1-3/8"x36" Strap', 'ea', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'HD2A Holdown', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'HD5A Holdown', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'HD7A Holdown', 'ea', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'HDU2 Holdown', 'ea', 20.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'HDU4 Holdown', 'ea', 28.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'ABU44 Adjustable Post Base 4x4', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'ABU66 Adjustable Post Base 6x6', 'ea', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'CB44 Column Base 4x4', 'ea', 15.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'CB46 Column Base 4x6', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'CB66 Column Base 6x6', 'ea', 22.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'Post Cap 4x4', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'Post Cap 6x6', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'Beam Seat 4x6', 'ea', 22.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Hardware', 'Beam Seat 4x8', 'ea', 25.00);

-- ============================================================
-- CONCRETE & MASONRY
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Ready-Mix Concrete per yard', 'yd', 175.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '60lb Concrete Mix (Quikrete)', 'bag', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '80lb Concrete Mix (Quikrete)', 'bag', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '50lb Fast-Setting Concrete', 'bag', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '60lb Mortar Mix Type S', 'bag', 8.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '80lb Mortar Mix Type N', 'bag', 9.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '#3 Rebar (3/8") 20''', 'ea', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '#4 Rebar (1/2") 20''', 'ea', 7.86);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '#5 Rebar (5/8") 20''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '#6 Rebar (3/4") 20''', 'ea', 16.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Rebar Tie Wire 3.5lb Roll', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Rebar Chairs 3" 50ct', 'bag', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Rebar Chairs 5" 25ct', 'bag', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '6x6 W1.4/W1.4 Welded Wire Mesh 5''x150''', 'ea', 125.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '6x6 W2.9/W2.9 Welded Wire Mesh 5''x150''', 'ea', 175.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '8" Sonotube per foot', 'lnft', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '10" Sonotube per foot', 'lnft', 4.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '14" Sonotube per foot', 'lnft', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '16" Sonotube per foot', 'lnft', 7.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '18" Sonotube per foot', 'lnft', 9.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', '24" Sonotube per foot', 'lnft', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Concrete Form Stakes 18"', 'ea', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Concrete Form Stakes 24"', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Concrete Form Stakes 36"', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Expansion Joint 1/2"x4"x10''', 'ea', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Expansion Joint 1/2"x6"x10''', 'ea', 7.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Vapor Barrier 6mil 20''x100''', 'ea', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Vapor Barrier 10mil 20''x100''', 'ea', 85.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Concrete Cure & Seal 5 gal', 'ea', 65.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Form Release Agent 5 gal', 'ea', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'CMU Block 8x8x16 Standard', 'ea', 2.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'CMU Block 8x8x16 Half', 'ea', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'CMU Block 12x8x16 Standard', 'ea', 3.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Concrete', 'Sill Plate Gasket 5-1/2"x50''', 'ea', 8.00);

-- ============================================================
-- DRYWALL & FINISHING
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '1/2" Drywall 4x8 Sheet', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '1/2" Drywall 4x10 Sheet', 'ea', 17.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '1/2" Drywall 4x12 Sheet', 'ea', 21.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '5/8" Drywall 4x8 Sheet', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '5/8" Drywall 4x10 Sheet', 'ea', 20.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '5/8" Drywall 4x12 Sheet', 'ea', 24.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '5/8" Type X Fire-Rated Drywall 4x8', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '5/8" Type X Fire-Rated Drywall 4x12', 'ea', 27.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '1/2" Moisture-Resistant Drywall 4x8', 'ea', 17.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '1/4" Drywall 4x8 Flexible', 'ea', 13.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'All-Purpose Joint Compound 4.5 gal', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Lightweight Joint Compound 4.5 gal', 'ea', 20.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Setting Compound 20min 18lb', 'bag', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Setting Compound 45min 18lb', 'bag', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Setting Compound 90min 18lb', 'bag', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Paper Drywall Tape 500''', 'ea', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Mesh Drywall Tape 300''', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Metal Corner Bead 8''', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Metal Corner Bead 10''', 'ea', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Paper-Faced Corner Bead 8''', 'ea', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Paper-Faced Corner Bead 10''', 'ea', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'L-Bead 10''', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'J-Bead 10''', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Bullnose Corner Bead 10''', 'ea', 5.50);

-- ============================================================
-- CEMENT BOARD & TILE BACKER
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '1/4" Cement Board 3x5', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', '1/2" Cement Board 3x5', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Kerdi Board 3/4" 24x72', 'ea', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'Kerdi Membrane per sqft', 'sqft', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Drywall', 'RedGard Waterproof Membrane 1 gal', 'ea', 35.00);

-- ============================================================
-- PLUMBING - PVC Pipe
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" PVC Sch 40 10''', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" PVC Sch 40 10''', 'ea', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1" PVC Sch 40 10''', 'ea', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1-1/2" PVC Sch 40 10''', 'ea', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '2" PVC Sch 40 10''', 'ea', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3" PVC Sch 40 10''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '4" PVC Sch 40 10''', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1-1/2" PVC DWV 10''', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '2" PVC DWV 10''', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3" PVC DWV 10''', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '4" PVC DWV 10''', 'ea', 18.00);

-- ============================================================
-- PLUMBING - PEX Tubing
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/8" PEX-A Tubing per foot', 'lnft', 0.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" PEX-A Tubing per foot', 'lnft', 0.65);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" PEX-A Tubing per foot', 'lnft', 1.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1" PEX-A Tubing per foot', 'lnft', 1.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" PEX-B Tubing per foot', 'lnft', 0.40);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" PEX-B Tubing per foot', 'lnft', 0.65);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1" PEX-B Tubing per foot', 'lnft', 1.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" PEX 90 Elbow', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" PEX 90 Elbow', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" PEX Tee', 'ea', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" PEX Tee', 'ea', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" PEX Ball Valve', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" PEX Ball Valve', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'PEX Manifold 6 Port', 'ea', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'PEX Manifold 8 Port', 'ea', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'PEX Manifold 12 Port', 'ea', 75.00);

-- ============================================================
-- PLUMBING - Copper Pipe
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" Copper Type M 10''', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" Copper Type M 10''', 'ea', 24.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1" Copper Type M 10''', 'ea', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1/2" Copper Type L 10''', 'ea', 20.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '3/4" Copper Type L 10''', 'ea', 32.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', '1" Copper Type L 10''', 'ea', 52.00);

-- ============================================================
-- PLUMBING - Fixtures & Misc
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'PVC Cement 16oz', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'PVC Primer 16oz', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'Pipe Clamps 1/2" (10pk)', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'Pipe Clamps 3/4" (10pk)', 'ea', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'Pipe Clamps 1" (10pk)', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'PTFE Thread Tape 1/2"x520"', 'ea', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'Water Heater 40 gal Gas', 'ea', 600.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'Water Heater 50 gal Gas', 'ea', 750.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Plumbing', 'Water Heater 50 gal Electric', 'ea', 550.00);

-- ============================================================
-- ELECTRICAL - Wire (Romex/NM-B)
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '14/2 NM-B Romex per foot', 'lnft', 0.35);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '14/2 NM-B Romex 250'' Roll', 'ea', 75.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '14/3 NM-B Romex per foot', 'lnft', 0.55);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '12/2 NM-B Romex per foot', 'lnft', 0.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '12/2 NM-B Romex 250'' Roll', 'ea', 110.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '12/3 NM-B Romex per foot', 'lnft', 0.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '10/2 NM-B Romex per foot', 'lnft', 0.85);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '10/3 NM-B Romex per foot', 'lnft', 1.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '8/3 NM-B Romex per foot', 'lnft', 2.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '6/3 NM-B Romex per foot', 'lnft', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '6/3 SER Cable per foot', 'lnft', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '4/3 SER Cable per foot', 'lnft', 6.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '2/3 SER Cable per foot', 'lnft', 9.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '14/2 UF-B Underground per foot', 'lnft', 0.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '12/2 UF-B Underground per foot', 'lnft', 0.70);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '10/2 UF-B Underground per foot', 'lnft', 1.10);

-- ============================================================
-- ELECTRICAL - Boxes & Devices
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Single Gang Old Work Box', 'ea', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Single Gang New Work Box', 'ea', 0.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Double Gang New Work Box', 'ea', 1.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Triple Gang New Work Box', 'ea', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '4" Octagon Box', 'ea', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '4" Round Pancake Box', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '4-11/16" Square Box', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Weatherproof Box Single Gang', 'ea', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '15A Duplex Outlet', 'ea', 1.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '20A Duplex Outlet', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '20A GFCI Outlet', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '15A Single Pole Switch', 'ea', 1.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '15A 3-Way Switch', 'ea', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '15A 4-Way Switch', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Dimmer Switch Single Pole', 'ea', 15.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Blank Wall Plate', 'ea', 0.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Single Gang Wall Plate', 'ea', 0.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Double Gang Wall Plate', 'ea', 1.25);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Wire Nuts Yellow (100ct)', 'box', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Wire Nuts Red (100ct)', 'box', 6.00);

-- ============================================================
-- ELECTRICAL - Panels & Breakers
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '100A Main Breaker Panel 20 Space', 'ea', 125.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '200A Main Breaker Panel 30 Space', 'ea', 175.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '200A Main Breaker Panel 40 Space', 'ea', 225.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Sub Panel 60A 12 Space', 'ea', 65.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', 'Sub Panel 100A 20 Space', 'ea', 95.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '15A Single Pole Breaker', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '20A Single Pole Breaker', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '30A Double Pole Breaker', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '40A Double Pole Breaker', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '50A Double Pole Breaker', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '15A AFCI Breaker', 'ea', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '20A AFCI Breaker', 'ea', 38.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '20A GFCI Breaker', 'ea', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '15A Dual Function (AFCI/GFCI) Breaker', 'ea', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '20A Dual Function (AFCI/GFCI) Breaker', 'ea', 48.00);

-- ============================================================
-- ELECTRICAL - Conduit
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '1/2" EMT Conduit 10''', 'ea', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '3/4" EMT Conduit 10''', 'ea', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '1" EMT Conduit 10''', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '1/2" PVC Conduit Sch 40 10''', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '3/4" PVC Conduit Sch 40 10''', 'ea', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '1" PVC Conduit Sch 40 10''', 'ea', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '1/2" Liquid-Tight Flex Conduit per foot', 'lnft', 1.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Electrical', '3/4" Liquid-Tight Flex Conduit per foot', 'lnft', 1.50);

-- ============================================================
-- ADHESIVES & SEALANTS
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Construction Adhesive (PL Premium) 10oz', 'ea', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Construction Adhesive (PL Premium) 28oz', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Construction Adhesive (Liquid Nails) 10oz', 'ea', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Construction Adhesive (Liquid Nails) 28oz', 'ea', 9.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Silicone Caulk White 10oz', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Silicone Caulk Clear 10oz', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Paintable Caulk (Alex Plus) 10oz', 'ea', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Polyurethane Caulk (Sika) 10oz', 'ea', 8.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Fire Caulk Intumescent 10oz', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Spray Foam Can (Great Stuff) 12oz', 'ea', 6.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Spray Foam Can (Great Stuff Pro) 24oz', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Wood Glue (Titebond III) 16oz', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Wood Glue (Titebond III) 1 gal', 'ea', 32.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Adhesives & Sealants', 'Epoxy 2-Part Quick Set', 'ea', 8.00);

-- ============================================================
-- PAINT & FINISHES
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Primer 1 gal', 'gal', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Primer 5 gal', 'ea', 95.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Exterior Primer 1 gal', 'gal', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Exterior Primer 5 gal', 'ea', 110.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Flat Paint 1 gal', 'gal', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Flat Paint 5 gal', 'ea', 140.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Eggshell Paint 1 gal', 'gal', 38.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Eggshell Paint 5 gal', 'ea', 155.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Semi-Gloss Paint 1 gal', 'gal', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Interior Semi-Gloss Paint 5 gal', 'ea', 165.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Exterior Flat Paint 1 gal', 'gal', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Exterior Satin Paint 1 gal', 'gal', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Exterior Satin Paint 5 gal', 'ea', 180.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Stain (Interior) 1 gal', 'gal', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Stain (Exterior) 1 gal', 'gal', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Stain (Exterior) 5 gal', 'ea', 165.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Polyurethane (Interior) 1 gal', 'gal', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Deck Stain 1 gal', 'gal', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Deck Stain 5 gal', 'ea', 165.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Painter''s Tape Blue 1.88"x60yd', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Drop Cloth 9x12 Canvas', 'ea', 15.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Roller Cover 9" (3pk)', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Paint Brush 2"', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Paint', 'Paint Brush 4"', 'ea', 12.00);

-- ============================================================
-- SIDING & EXTERIOR
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide 38 Series 8" Lap 16''', 'ea', 22.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide 38 Series 12" Lap 16''', 'ea', 28.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide Panel 4x8', 'ea', 42.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide Panel 4x9', 'ea', 48.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide Panel 4x10', 'ea', 52.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide Trim 1x4x16', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide Trim 1x6x16', 'ea', 24.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide Trim 1x8x16', 'ea', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'LP SmartSide Trim 1x12x16', 'ea', 42.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Hardie Plank 8.25" Lap 12''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Hardie Panel 4x8 Smooth', 'ea', 38.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Hardie Trim 5/4x4x12', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Hardie Trim 5/4x6x12', 'ea', 24.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Hardie Trim 5/4x8x12', 'ea', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Hardie Trim 5/4x12x12', 'ea', 42.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Vinyl Siding D4.5 per sqft', 'sqft', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Vinyl Soffit Vented per sqft', 'sqft', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Vinyl J-Channel 12.5''', 'ea', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Vinyl Starter Strip 10''', 'ea', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Vinyl Utility Trim 12.5''', 'ea', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Vinyl Outside Corner 10''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Vinyl Inside Corner 10''', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Aluminum Fascia 6"x12''', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Aluminum Fascia 8"x12''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Siding', 'Aluminum Drip Cap 10''', 'ea', 6.00);

-- ============================================================
-- DECKING
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '5/4x6x8 PT Deck Board', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '5/4x6x10 PT Deck Board', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '5/4x6x12 PT Deck Board', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '5/4x6x14 PT Deck Board', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '5/4x6x16 PT Deck Board', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '2x6x8 PT Deck Board', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '2x6x10 PT Deck Board', 'ea', 12.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '2x6x12 PT Deck Board', 'ea', 15.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', '2x6x16 PT Deck Board', 'ea', 20.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Composite Decking 1x6x12 (Trex Select)', 'ea', 28.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Composite Decking 1x6x16 (Trex Select)', 'ea', 38.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Composite Decking 1x6x20 (Trex Select)', 'ea', 48.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Composite Decking 1x6x12 (Trex Enhance)', 'ea', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Composite Decking 1x6x16 (Trex Enhance)', 'ea', 47.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Composite Decking 1x6x20 (Trex Enhance)', 'ea', 58.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Composite Fascia Board 1x12x12', 'ea', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Hidden Deck Fasteners (90ct)', 'box', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Deck Screws #10 2-1/2" 350ct', 'box', 28.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Deck Screws #10 3" 350ct', 'box', 32.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Decking', 'Deck Joist Tape 1-5/8"x50''', 'ea', 14.00);

-- ============================================================
-- DOORS
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Door Slab 2/0x6/8 Hollow Core', 'ea', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Door Slab 2/4x6/8 Hollow Core', 'ea', 48.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Door Slab 2/6x6/8 Hollow Core', 'ea', 50.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Door Slab 2/8x6/8 Hollow Core', 'ea', 52.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Door Slab 3/0x6/8 Hollow Core', 'ea', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Prehung Door 2/4x6/8', 'ea', 120.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Prehung Door 2/6x6/8', 'ea', 125.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Prehung Door 2/8x6/8', 'ea', 130.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Prehung Door 3/0x6/8', 'ea', 140.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Exterior Prehung Door 3/0x6/8 Fiberglass', 'ea', 350.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Exterior Prehung Door 3/0x6/8 Steel', 'ea', 250.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Sliding Patio Door 6/0x6/8', 'ea', 550.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'French Door Pair 5/0x6/8', 'ea', 800.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Bifold Door 2/0x6/8', 'ea', 60.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Bifold Door 2/6x6/8', 'ea', 65.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Bifold Door 3/0x6/8', 'ea', 70.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Door Hinges 3.5" (3pk)', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Door Knob Passage', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Interior Door Knob Privacy', 'ea', 15.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Exterior Door Deadbolt', 'ea', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Exterior Door Handle Set', 'ea', 45.00);

-- ============================================================
-- WINDOWS
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Double Hung Window 24x36', 'ea', 175.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Double Hung Window 30x48', 'ea', 200.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Double Hung Window 30x60', 'ea', 225.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Double Hung Window 36x48', 'ea', 225.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Double Hung Window 36x60', 'ea', 250.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Double Hung Window 36x72', 'ea', 275.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Casement Window 24x36', 'ea', 200.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Casement Window 24x48', 'ea', 225.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Casement Window 30x48', 'ea', 250.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Casement Window 30x60', 'ea', 275.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Picture Window 48x48', 'ea', 250.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Picture Window 48x60', 'ea', 300.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Picture Window 60x48', 'ea', 325.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Egress Window 48x36', 'ea', 275.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Basement Hopper Window 32x14', 'ea', 100.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Skylight 22x46', 'ea', 350.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Doors & Windows', 'Skylight 30x46', 'ea', 425.00);

-- ============================================================
-- TRIM & MOLDING
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Baseboard 3-1/4" MDF 16''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Baseboard 5-1/4" MDF 16''', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Baseboard 3-1/4" Pine 16''', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Baseboard 5-1/4" Pine 16''', 'ea', 24.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Casing 2-1/4" MDF 7''', 'ea', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Casing 3-1/4" MDF 7''', 'ea', 5.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Casing 2-1/4" Pine 7''', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Casing 3-1/4" Pine 7''', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Crown Molding 3-5/8" MDF 16''', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Crown Molding 5-1/4" MDF 16''', 'ea', 24.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Crown Molding 3-5/8" Pine 8''', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Quarter Round 3/4" MDF 8''', 'ea', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Shoe Molding 1/2"x3/4" MDF 16''', 'ea', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Chair Rail 2-1/2" 8''', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Window Stool 3/4"x3-1/2" 8''', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Window Apron 2-1/4" 8''', 'ea', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Rosette Block 3-1/4"', 'ea', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Trim & Molding', 'Plinth Block 3-1/4"', 'ea', 3.50);

-- ============================================================
-- FLOORING
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Luxury Vinyl Plank (LVP) per sqft', 'sqft', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Luxury Vinyl Tile (LVT) per sqft', 'sqft', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Laminate Flooring per sqft', 'sqft', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Engineered Hardwood per sqft', 'sqft', 5.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Solid Hardwood Oak per sqft', 'sqft', 6.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Ceramic Tile 12x12 per sqft', 'sqft', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Ceramic Tile 12x24 per sqft', 'sqft', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Porcelain Tile 12x24 per sqft', 'sqft', 3.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Porcelain Tile 24x24 per sqft', 'sqft', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Carpet per sqft (mid-grade)', 'sqft', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Carpet Pad per sqft', 'sqft', 0.75);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Floor Underlayment 100sqft Roll', 'ea', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Tile Thinset Mortar 50lb', 'bag', 22.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Tile Grout Sanded 25lb', 'bag', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Tile Grout Unsanded 10lb', 'bag', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Transition Strip T-Molding 36"', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Flooring', 'Transition Strip Reducer 36"', 'ea', 10.00);

-- ============================================================
-- HVAC
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '6" Round Duct 5''', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '7" Round Duct 5''', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '8" Round Duct 5''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '10" Round Duct 5''', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '12" Round Duct 5''', 'ea', 20.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '6" Insulated Flex Duct 25''', 'ea', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '8" Insulated Flex Duct 25''', 'ea', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '10" Insulated Flex Duct 25''', 'ea', 50.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '12" Insulated Flex Duct 25''', 'ea', 60.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '4x10 Floor Register', 'ea', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '4x12 Floor Register', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '6x10 Floor Register', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '6x12 Floor Register', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '8x14 Return Air Grille', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '14x20 Return Air Grille', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', '20x25 Return Air Grille', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'Duct Tape 2"x60yd', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'Foil Duct Tape 2.5"x60yd', 'ea', 14.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'Duct Mastic 1 gal', 'ea', 15.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'HVAC Hanger Strap 1"x100''', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'Bath Exhaust Fan 80 CFM', 'ea', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'Bath Exhaust Fan 110 CFM', 'ea', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'Range Hood Vent 30"', 'ea', 150.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'HVAC', 'Dryer Vent Kit 4"x8''', 'ea', 15.00);

-- ============================================================
-- LUMBER - Additional PT & Common Sizes (extending existing)
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x4x8 PT', 'ea', 6.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x4x10 PT', 'ea', 8.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x4x12 PT', 'ea', 10.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x4x16 PT', 'ea', 13.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x6x8 PT', 'ea', 9.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x6x10 PT', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x6x12 PT', 'ea', 14.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x8x12 PT', 'ea', 22.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x10x12 PT', 'ea', 28.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '2x12x12 PT', 'ea', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '1x4x8 Pine', 'ea', 4.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '1x6x8 Pine', 'ea', 6.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '1x8x8 Pine', 'ea', 8.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '1x10x8 Pine', 'ea', 12.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '1x12x8 Pine', 'ea', 16.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '3/4" OSB 4x8', 'ea', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '7/16" OSB 4x8', 'ea', 18.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '3/4" Plywood Sanded 4x8', 'ea', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '1/2" Plywood Sanded 4x8', 'ea', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('material', 'Lumber', '1/4" Plywood Luan 4x8', 'ea', 18.00);

-- ============================================================
-- GENERIC LABOR RATES
-- ============================================================
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'General Laborer per hour', 'hr', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Skilled Carpenter per hour', 'hr', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Journeyman Electrician per hour', 'hr', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Journeyman Plumber per hour', 'hr', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'HVAC Technician per hour', 'hr', 55.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Painter per hour', 'hr', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Drywall Hanger per hour', 'hr', 35.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Drywall Finisher per hour', 'hr', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Concrete Worker per hour', 'hr', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Tile Setter per hour', 'hr', 45.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Roofer per hour', 'hr', 40.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Framing per sqft', 'sqft', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Drywall Hang & Finish per sqft', 'sqft', 2.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Painting per sqft', 'sqft', 1.50);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Roofing Install per sqft', 'sqft', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Siding Install per sqft', 'sqft', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Tile Install per sqft', 'sqft', 7.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Flooring Install per sqft (LVP/Laminate)', 'sqft', 2.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Flooring Install per sqft (Hardwood)', 'sqft', 4.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Concrete Flatwork per sqft', 'sqft', 8.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Concrete Foundation per lnft', 'lnft', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Demolition per hour', 'hr', 30.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Cleanup per hour', 'hr', 25.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Trim Carpentry per lnft', 'lnft', 3.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Door Hang & Trim per door', 'ea', 150.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Window Install per window', 'ea', 175.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Electrical Rough-In per outlet/switch', 'ea', 85.00);
INSERT INTO item_templates (type, category, name, default_unit, default_price) VALUES ('labor', 'Labor', 'Plumbing Rough-In per fixture', 'ea', 250.00);

-- +goose Down
-- Remove only the items added by this migration
-- Using category names unique to this migration for clean rollback
DELETE FROM item_templates WHERE category IN (
    'Fasteners', 'Hardware', 'Concrete', 'Drywall',
    'Plumbing', 'Electrical', 'Adhesives & Sealants',
    'Paint', 'Siding', 'Decking', 'Doors & Windows',
    'Trim & Molding', 'Flooring', 'HVAC'
);
-- Remove PT lumber and additional sizes added to existing Lumber category
DELETE FROM item_templates WHERE category = 'Lumber' AND name IN (
    '2x4x8 PT', '2x4x10 PT', '2x4x12 PT', '2x4x16 PT',
    '2x6x8 PT', '2x6x10 PT', '2x6x12 PT',
    '2x8x12 PT', '2x10x12 PT', '2x12x12 PT',
    '1x4x8 Pine', '1x6x8 Pine', '1x8x8 Pine', '1x10x8 Pine', '1x12x8 Pine',
    '3/4" OSB 4x8', '7/16" OSB 4x8',
    '3/4" Plywood Sanded 4x8', '1/2" Plywood Sanded 4x8', '1/4" Plywood Luan 4x8'
);
-- Remove generic labor rates (keep original named employees)
DELETE FROM item_templates WHERE category = 'Labor' AND name IN (
    'General Laborer per hour', 'Skilled Carpenter per hour',
    'Journeyman Electrician per hour', 'Journeyman Plumber per hour',
    'HVAC Technician per hour', 'Painter per hour',
    'Drywall Hanger per hour', 'Drywall Finisher per hour',
    'Concrete Worker per hour', 'Tile Setter per hour', 'Roofer per hour',
    'Framing per sqft', 'Drywall Hang & Finish per sqft',
    'Painting per sqft', 'Roofing Install per sqft',
    'Siding Install per sqft', 'Tile Install per sqft',
    'Flooring Install per sqft (LVP/Laminate)', 'Flooring Install per sqft (Hardwood)',
    'Concrete Flatwork per sqft', 'Concrete Foundation per lnft',
    'Demolition per hour', 'Cleanup per hour', 'Trim Carpentry per lnft',
    'Door Hang & Trim per door', 'Window Install per window',
    'Electrical Rough-In per outlet/switch', 'Plumbing Rough-In per fixture'
);
