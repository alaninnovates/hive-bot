"use server";
import {getCurrentSession} from "@/lib/server/session";

export default async function Page() {
    const {user} = await getCurrentSession();
    if (user === null) {
        console.log(user);
        return (
            <>
                <h1>log in</h1>
                <a href="/login/discord">Sign in with Discord</a>
            </>
        );
    }
    return <>
        {JSON.stringify(user, null, 2)}
    </>
}