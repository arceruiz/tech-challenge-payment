
print('Start #################################################################');

db = db.getSiblingDB('payment');
db.createCollection('payments');

print('END #################################################################');