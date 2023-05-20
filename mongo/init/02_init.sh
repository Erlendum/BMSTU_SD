#!/bin/sh 

mongoimport -d MusicStore -c instruments --type csv --file /data/data/instruments.csv --headerline
mongoimport -d MusicStore -c discounts --type csv --file /data/data/discounts.csv --headerline
mongoimport -d MusicStore -c users --type csv --file /data/data/users.csv --headerline
mongoimport -d MusicStore -c comparisonLists --type csv --file /data/data/comparisonLists.csv --headerline
mongoimport -d MusicStore -c comparisonLists_instruments --type csv --file /data/data/comparisonLists_instruments.csv --headerline