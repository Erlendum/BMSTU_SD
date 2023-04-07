drop schema if exists store cascade;

create SCHEMA store;

create table store.Instruments
(
	instrument_id serial primary key,
	instrument_name text not null,
	instrument_price int,
	instrument_material text not null,
	instrument_type text not null,
	instrument_brand text not null,
	instrument_img text not null
);	

create table store.Users
(
	user_id serial primary key,
	user_login text not null,
	user_password text not null,
	user_fio text not null,
	user_date_birth date,
	user_gender text not null,
	user_is_admin boolean
);	

create table store.Discounts
(
	discount_id serial primary key,
	instrument_id int not null,
	user_id int not null,
	foreign key (instrument_id) references store.Instruments(instrument_id) on delete cascade,
	foreign key (user_id) references store.Users(user_id) on delete cascade,
	discount_amount int,
	discount_type text not null,
	discount_date_begin date,
	discount_date_end date
);

create table store.ComparisonLists
(
	comparisonList_id serial primary key,
	user_id int not null,
	foreign key (user_id) references store.Users(user_id) on delete cascade,
	comparisonList_total_price int,
	comparisonList_amount int
);

create table store.comparisonLists_instruments (
    comparisonLists_instruments_id serial primary key,
    comparisonList_id int not null,
    instrument_id int not null,
    FOREIGN KEY (comparisonList_id) references store.comparisonLists(comparisonList_id) on delete cascade,
    FOREIGN KEY (instrument_id) references store.Instruments(instrument_id) on delete cascade
);


