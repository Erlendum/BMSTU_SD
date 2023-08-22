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

create table store.Orders
(
    order_id serial primary key,
    order_time date,
    order_price int,
    order_status text default 'Created',
    user_id int not null,
    foreign key (user_id) references store.Users(user_id) on delete cascade
);

create table store.Order_elements
(
    order_element_id serial primary key,
    order_element_amount int,
    order_element_price int,
    instrument_id int not null,
    order_id int not null,
    FOREIGN KEY (instrument_id) references store.Instruments(instrument_id) on delete cascade,
    foreign key (order_id) references store.Orders(order_id) on delete cascade
);


alter table store.Users
    add constraint correct_user_gender CHECK (user_gender = 'Male' OR user_gender = 'Female');

alter table store.Orders
    add constraint correct_order_status check (order_status = 'Created' or order_status = 'Delivering' or order_status = 'Delivered');


insert into store.instruments (instrument_name, instrument_price, instrument_material, instrument_type, instrument_brand, instrument_img)
values ('I1', 2500, 'Дуб', 'Акустические гитары', 'no', 'no');

insert into store.users (user_login, user_password, user_fio, user_date_birth, user_gender, user_is_admin)
values ('admin', '$2a$10$1OOImwqdj8VCsC10WmVZZOqQqZ3roHRhFq69jZSwOdUQOZhzczv4S', 'admin', '2003-01-22', 'Male', true);

insert into store.comparisonLists(user_id, comparisonList_total_price, comparisonList_amount)
values (1, 0, 0);