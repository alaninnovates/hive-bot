import pymongo
from pymongo import MongoClient

conn = MongoClient('mongodb+srv://admin:R09n16mwx0iMXwLV@cluster0.n7wrnue.mongodb.net/?retryWrites=true&w=majority')
db = conn.get_database("hive-bot")
hives = db.get_collection("hives")

h = hives.find({})

for hive in h:
    # turn hive["bees"] from an object of objects to an object of arrays
    bees = hive["bees"]
    for key in bees:
        bees[key] = [bees[key]]
    # update doc
    hives.update_one({"_id": hive["_id"]}, {"$set": {"bees": bees}})
    print(hive["_id"], "updated")