import {Discord} from "arctic";

const clientId = process.env.DISCORD_CLIENT_ID;
const clientSecret = process.env.DISCORD_CLIENT_SECRET;
const redirectURI = process.env.DISCORD_REDIRECT_URI;

if (!clientId || !clientSecret || !redirectURI) {
    throw new Error("Missing Discord OAuth configuration");
}

export const discord = new Discord(clientId, clientSecret, redirectURI);
