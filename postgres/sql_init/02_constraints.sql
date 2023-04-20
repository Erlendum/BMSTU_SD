alter table store.Users
    add constraint correct_user_gender CHECK (user_gender = 'Male' OR user_gender = 'Female');

