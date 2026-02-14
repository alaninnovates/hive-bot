import { generateSessionToken, createSession, setSessionTokenCookie } from "@/lib/server/session";
import { discord } from "@/lib/server/oauth";
import { cookies } from "next/headers";
import { createUser, getUserFromDiscordId } from "@/lib/server/user";
import { ObjectParser } from "@pilcrowjs/object-parser";

import type { OAuth2Tokens } from "arctic";

export async function GET(request: Request): Promise<Response> {
    const c = await cookies();
    const url = new URL(request.url);
    const code = url.searchParams.get("code");
    const state = url.searchParams.get("state");
    const storedState = c.get("discord_oauth_state")?.value ?? null;
    if (code === null || state === null || storedState === null) {
        return new Response("Please restart the process.", {
            status: 400
        });
    }
    if (state !== storedState) {
        return new Response("Please restart the process.", {
            status: 400
        });
    }

    let tokens: OAuth2Tokens;
    try {
        tokens = await discord.validateAuthorizationCode(code, null);
    } catch {
        // Invalid code or client credentials
        return new Response("Please restart the process.", {
            status: 400
        });
    }
    const discordAccessToken = tokens.accessToken();

    const userRequest = new Request("https://discord.com/api/users/@me");
    userRequest.headers.set("Authorization", `Bearer ${discordAccessToken}`);
    const userResponse = await fetch(userRequest);
    const userResult: unknown = await userResponse.json();
    const userParser = new ObjectParser(userResult);

    const discordUserId = parseInt(userParser.getString("id"));
    const username = userParser.getString("username");
    const globalName = userParser.get("global_name") ? userParser.getString("global_name") : undefined;
    const profileUrl = userParser.get("avatar") ?
            `https://cdn.discordapp.com/avatars/${discordUserId}/${userParser.getString("avatar")}.png`
        : `https://cdn.discordapp.com/embed/avatars/${(discordUserId >> 22) % 6}.png`;
    const email = userParser.getString("email");

    const existingUser = await getUserFromDiscordId(discordUserId);
    if (existingUser != null) {
        const sessionToken = generateSessionToken();
        const session = await createSession(sessionToken, existingUser.discordId);
        await setSessionTokenCookie(sessionToken, session.expiresAt);
        return new Response(null, {
            status: 302,
            headers: {
                Location: "/"
            }
        });
    }

    const userData = {
        discordId: discordUserId,
        email,
        username,
        globalName,
        profileUrl
    };
    const user = await createUser(userData);
    const sessionToken = generateSessionToken();
    const session = await createSession(sessionToken, user.discordId);
    await setSessionTokenCookie(sessionToken, session.expiresAt);
    return new Response(null, {
        status: 302,
        headers: {
            Location: "/"
        }
    });
}