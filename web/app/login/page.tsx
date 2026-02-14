import { getCurrentSession } from "@/lib/server/session";
import { redirect } from "next/navigation";

export default async function Page() {
    const { user } = await getCurrentSession();
    if (user !== null) {
        return redirect("/");
    }
    return (
        <>
            <h1>Sign in</h1>
            <a href="/login/discord">Sign in with Discord</a>
        </>
    );
}