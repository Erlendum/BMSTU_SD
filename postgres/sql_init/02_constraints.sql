alter table store.Users
    add constraint correct_user_gender CHECK (user_gender = 'Мужской' OR user_gender = 'Женский');

