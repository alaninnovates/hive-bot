import {getHivePosts} from "@/lib/database/posts";

export default async function Page(props: PageProps<'/posts/[id]'>) {
    const {id} = await props.params;
    const {post, hive, user} = (await getHivePosts({postId: id}))[0];
    return (
        <div className="container mx-auto py-12 px-4 grow">
            <a href="/posts" className="text-sm text-gray-500 hover:text-gray-700 mb-4 inline-block">
                &larr; Back to Posts
            </a>
            <div className="mt-4 mb-6">
                <div className="mt-4 mb-2">
                    <h1 className="text-4xl font-bold">{post.title}</h1>
                    <p className="text-gray-500 text-sm">By {user.globalName}</p>
                    <p className="text-gray-400 text-xs mt-1">
                        {new Date(post.createdAt).toLocaleDateString()}
                    </p>
                </div>
                <div className="prose max-w-none">
                    <p>{post.content}</p>
                </div>
            </div>
            <img
                src={post.imageUrl}
                alt={post.title}
                className="max-w-sm w-full object-cover"
            />
            <div className="mt-8">
                <h3 className="text-xl font-semibold mb-2">Hive Details</h3>
            </div>
        </div>
    );
}