import {generateState} from "arctic";
import {discord} from "@/lib/server/oauth";
import {cookies} from "next/headers";

export async function GET(): Promise<Response> {
    const state = generateState();
    const url = discord.createAuthorizationURL(state, null, ["identify", "email"]);

    const c = await cookies();
    c.set("discord_oauth_state", state, {
        path: "/",
        secure: process.env.NODE_ENV === "production",
        httpOnly: true,
        maxAge: 60 * 10,
        sameSite: "lax"
    });

    return new Response(null, {
        status: 302,
        headers: {
            Location: url.toString()
        }
    });
}