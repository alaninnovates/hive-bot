import {getHivePosts} from "@/lib/database/posts";

export default async function Page(props: PageProps<'/posts/[id]'>) {
    const {id} = await props.params;
    const {post, hive, user} = (await getHivePosts({postId: id}))[0];
    return (
        <div className="container mx-auto py-12 px-4 grow">
            <h1 className="text-4xl font-bold">Post {id}</h1>
            <p className="text-gray-600">This is a placeholder for post {id}.</p>
        </div>
    );
}