"use server";
import {getCurrentSession} from "@/lib/server/session";
import {getHivePosts} from "@/lib/database/posts";

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
    const posts = await getHivePosts();
    console.log(posts);

    return (
        <>
            {JSON.stringify(posts, null, 2)}
        </>
    )
}