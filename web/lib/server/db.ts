import { MongoClient, MongoClientOptions } from 'mongodb';

const uri = process.env.MONGODB_URI;

if (!uri) {
    throw new Error('Missing MongoDB URI');
}

const options: MongoClientOptions = {
};

let client: MongoClient;
let clientPromise: Promise<MongoClient> | null = null;

if (uri) {
    if (process.env.NODE_ENV === 'development') {
        // In development mode, use a global variable so that the value
        // is preserved across module reloads caused by HMR (Hot Module Replacement).
        const globalWithMongo = global as typeof globalThis & {
            _mongoClientPromise?: Promise<MongoClient>;
        };

        if (!globalWithMongo._mongoClientPromise) {
            client = new MongoClient(uri, options);
            globalWithMongo._mongoClientPromise = client.connect();
        }
        clientPromise = globalWithMongo._mongoClientPromise;
    } else {
        client = new MongoClient(uri, options);
        clientPromise = client.connect();
    }
}

const getDb = async () => {
    if (!clientPromise) {
        throw new Error('MongoDB client is not initialized');
    }
    const client = await clientPromise;
    return client.db('hive-bot');
}

export {clientPromise, getDb};