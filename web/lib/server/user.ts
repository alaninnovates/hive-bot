import {getDb} from "./db";
import {Long} from "bson";

export async function createUser(user: User): Promise<User> {
    const db = await getDb();
    const result = await db.collection("users").insertOne({
        discord_id: user.discordId,
        email: user.email,
        username: user.username,
        profile_url: user.profileUrl ?? null,
        global_name: user.globalName ?? null
    });
    if (!result.acknowledged) {
        throw new Error("Unexpected error");
    }
    return user;
}

export async function getUserFromDiscordId(discordId: string): Promise<User | null> {
    const db = await getDb();
    const row = await db.collection("users").findOne({discord_id: discordId});
    if (row === null) {
        return null;
    }
    return {
        discordId: row.discord_id,
        email: row.email,
        username: row.username,
        profileUrl: row.profile_url ?? undefined,
        globalName: row.global_name ?? undefined
    }
}

export interface User {
    discordId: string;
    email: string;
    username: string;
    profileUrl?: string;
    globalName?: string;
}
