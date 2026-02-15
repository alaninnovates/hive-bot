"use server";
import {getHivePosts} from "@/lib/database/posts";

export default async function Page() {
    const posts = await getHivePosts();
    console.log(posts);

    return (
        <>
            {JSON.stringify(posts, null, 2)}
        </>
    )
}