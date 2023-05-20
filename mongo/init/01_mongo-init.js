db = db.getSiblingDB('MusicStore');


db.createUser(
        {
            user: "erlendum",
            pwd: "parasha",
            roles: [
                {
                    role: "readWrite",
                    db: "MusicStore"
                }
            ]
        }
);

db.createCollection('instruments');
db.instruments.createIndex( {instrument_id : 1} , {unique : true} );
db.createCollection('discounts');
db.discounts.createIndex( {discount_id : 1} , {unique : true} );
db.createCollection('users', {
    validator: { $or:
        [
            {user_gender: {$in: ["Male", "Female"]}}
        ]
    }
});
db.users.createIndex( {user_id : 1} , {unique : true} );
db.createCollection('comparisonLists');
db.comparisonLists.createIndex( {comparisonList_id : 1} , {unique : true} );
db.createCollection('comparisonLists_instruments');
db.comparisonLists_instruments.createIndex( {comparisonLists_instruments_id : 1} , {unique : true} );
db.createCollection('orders', {
    validator: { $or:
        [
            {order_status: {$in: ["Created", "Delivering", "Delivered"]}}
        ]
    }
});
db.orders.createIndex( {order_id : 1} , {unique : true} );
db.createCollection('orderElements');
db.orderElements.createIndex( {order_element_id : 1} , {unique : true} );



