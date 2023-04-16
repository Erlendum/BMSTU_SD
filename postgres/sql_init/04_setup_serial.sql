BEGIN;
LOCK TABLE store.instruments IN EXCLUSIVE MODE;
SELECT setval('store.instruments_instrument_id_seq', COALESCE((SELECT MAX(instrument_id)+1 FROM store.instruments), 1), false);
COMMIT;

BEGIN;
LOCK TABLE store.users IN EXCLUSIVE MODE;
SELECT setval('store.users_user_id_seq', COALESCE((SELECT MAX(user_id)+1 FROM store.users), 1), false);
COMMIT;

BEGIN;
LOCK TABLE store.discounts IN EXCLUSIVE MODE;
SELECT setval('store.discounts_discount_id_seq', COALESCE((SELECT MAX(discount_id)+1 FROM store.discounts), 1), false);
COMMIT;

BEGIN;
LOCK TABLE store.comparisonLists IN EXCLUSIVE MODE;
SELECT setval('store.comparisonLists_comparisonList_id_seq', COALESCE((SELECT MAX(comparisonList_id)+1 FROM store.comparisonLists), 1), false);
COMMIT;