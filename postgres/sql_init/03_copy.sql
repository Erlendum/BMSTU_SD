copy store.Instruments FROM '/var/lib/postgresql/data/instruments.csv' DELIMITER ';' CSV HEADER;
copy store.Users FROM '/var/lib/postgresql/data/users.csv' DELIMITER ';' CSV HEADER;
copy store.Discounts FROM '/var/lib/postgresql/data/discounts.csv' DELIMITER ';' CSV HEADER;
copy store.ComparisonLists FROM '/var/lib/postgresql/data/comparisonLists.csv' DELIMITER ';' CSV HEADER;
copy store.ComparisonLists_instruments FROM '/var/lib/postgresql/data/comparisonLists_instruments.csv' DELIMITER ';' CSV HEADER;
