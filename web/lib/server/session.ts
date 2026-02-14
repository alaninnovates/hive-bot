import {getDb} from "./db";
import {encodeBase32, encodeHexLowerCase} from "@oslojs/encoding";
import {sha256} from "@oslojs/crypto/sha2";
import {cache} from "react";
import {cookies} from "next/headers";

import type {User} from "./user";

export async function validateSessionToken(token: string): Promise<SessionValidationResult> {
    const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
    const db = await getDb();
    const rows = await db.collection("sessions").aggregate([
        {$match: {id: sessionId}},
        {
            $lookup: {
                from: "users",
                localField: "discord_id",
                foreignField: "discord_id",
                as: "user"
            }
        },
        {$unwind: "$user"},
        {
            $project: {
                "session.id": "$id",
                "session.discord_id": "discord_id",
                "session.expires_at": "$expires_at",
                "user.discord_id": "$user.discord_id",
                "user.email": "$user.email",
                "user.username": "$user.username",
                "user.global_name": "$user.global_name",
                "user.profile_url": "$user.profile_url"
            }
        }
    ]).toArray();
    if (rows.length === 0) {
        return {session: null, user: null};
    }
    const row = rows[0];
    const session: Session = {
        id: row.session.id,
        discordId: row.session.discord_id,
        expiresAt: new Date(row.session.expires_at * 1000)
    };
    const user: User = {
        discordId: row.user.discord_id,
        email: row.user.email,
        username: row.user.username,
        globalName: row.user.global_name ?? undefined,
        profileUrl: row.user.profile_url ?? undefined
    };
    if (Date.now() >= session.expiresAt.getTime()) {
        await db.collection("sessions").deleteOne({id: session.id});
        return {session: null, user: null};
    }
    if (Date.now() >= session.expiresAt.getTime() - 1000 * 60 * 60 * 24 * 15) {
        session.expiresAt = new Date(Date.now() + 1000 * 60 * 60 * 24 * 30);
        await db.collection("sessions").updateOne({id: session.id}, {$set: {expires_at: Math.floor(session.expiresAt.getTime() / 1000)}});
    }
    return {session, user};
}

export const getCurrentSession = cache(async (): Promise<SessionValidationResult> => {
    const c = await cookies();
    const token = c.get("session")?.value ?? null;
    console.log("session", token);
    if (token === null) {
        return {session: null, user: null};
    }
    return validateSessionToken(token);
});

export async function invalidateSession(sessionId: string): Promise<void> {
    const db = await getDb();
    await db.collection("sessions").deleteOne({id: sessionId});
}

export async function invalidateUserSessions(discordId: number): Promise<void> {
    const db = await getDb();
    await db.collection("sessions").deleteMany({discord_id: discordId});
}

export async function setSessionTokenCookie(token: string, expiresAt: Date): Promise<void> {
    const c = await cookies();
    c.set("session", token, {
        httpOnly: true,
        path: "/",
        secure: process.env.NODE_ENV === "production",
        sameSite: "lax",
        expires: expiresAt
    });
}

export async function deleteSessionTokenCookie(): Promise<void> {
    const c = await cookies();
    c.set("session", "", {
        httpOnly: true,
        path: "/",
        secure: process.env.NODE_ENV === "production",
        sameSite: "lax",
        maxAge: 0
    });
}

export function generateSessionToken(): string {
    const tokenBytes = new Uint8Array(20);
    crypto.getRandomValues(tokenBytes);
    return encodeBase32(tokenBytes).toLowerCase();
}

export async function createSession(token: string, discordId: string): Promise<Session> {
    const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
    const session: Session = {
        id: sessionId,
        discordId,
        expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24 * 30)
    };
    const db = await getDb();
    await db.collection("sessions").insertOne({
        id: session.id,
        discord_id: session.discordId,
        expires_at: Math.floor(session.expiresAt.getTime() / 1000)
    });
    return session;
}

export interface Session {
    id: string;
    expiresAt: Date;
    discordId: string;
}

type SessionValidationResult = { session: Session; user: User } | { session: null; user: null };