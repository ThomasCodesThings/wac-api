const mongoHost = process.env.DEPARTMENT_API_MONGODB_HOST
const mongoPort = process.env.DEPARTMENT_API_MONGODB_PORT

const mongoUser = process.env.DEPARTMENT_API_MONGODB_USERNAME
const mongoPassword = process.env.DEPARTMENT_API_MONGODB_PASSWORD

const database = "cernica-department" //process.env.DEPARTMENT_API_MONGODB_DATABASE
const collection = ["department", "operation"]

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

let successFullCount = 0;
// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    for (let coll of collection) {
        if (collections.includes(coll)) {
            print(`Collection '${coll}' already exists in database '${database}'`)
            successFullCount++;
        }
    }
}

if (successFullCount === collection.length) {
    print("All collections already exists")
    process.exit(0);
}

// initialize
// create database and collection
const db = connection.getDB(database)
for (let coll of collection) {
    db.createCollection(coll)
    print(`Collection '${coll}' created in database '${database}'`)
}

// exit with success
process.exit(0);