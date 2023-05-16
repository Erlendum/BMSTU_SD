alter table store.Users
    add constraint correct_user_gender CHECK (user_gender = 'Male' OR user_gender = 'Female');

alter table store.Orders
    add constraint correct_order_status check (order_status = 'Created' or order_status = 'Delivering' or order_status = 'Delivered');
