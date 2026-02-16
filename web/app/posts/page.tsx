import {getHivePosts} from "@/lib/database/posts";
import {ScrollArea, ScrollBar} from "@/components/ui/scroll-area";

const officialBuilds = [
    {
        link: 'https://meta-bee.com/ultimate-blue-hive-meta/',
        name: 'The Ultimate Blue Hive Guide 2026',
        author: 'Acid',
        image: 'https://meta-bee.com/wp-content/uploads/2023/08/blue-hive.png',
        date: new Date('Feb 1, 2026'),
    },
    {
        link: 'https://meta-bee.com/red-hive-guide/',
        name: 'The Ultimate Red Hive Guide 2026',
        author: 'Acid',
        image: 'https://meta-bee.com/wp-content/uploads/2023/08/red-hive.png',
        date: new Date('Feb 1, 2026'),
    },
    {
        link: 'https://meta-bee.com/the-ultimate-white-hive-guide-for-endgame/',
        name: 'The Ultimate Endgame White Hive Guide! 2026',
        author: 'Acid',
        image: 'https://meta-bee.com/wp-content/uploads/2024/02/white-hive.png',
        date: new Date('Feb 1, 2026'),
    }
]

export default async function Page() {
    const hives = await getHivePosts();
    return (
        <div className="container mx-auto py-12 px-4 grow">
            <div className="mb-8">
                <h1 className="text-4xl font-bold">Hive Posts</h1>
                <p className="text-gray-600">Browse hive posts created by the community.</p>
            </div>
            <div className="mb-8">
                <h2 className="text-2xl font-bold mb-4">Official Meta Bee Builds</h2>
                <ScrollArea className="w-full">
                    <div className="flex space-x-4">
                        {officialBuilds.map((hive) => (
                            <a
                                key={hive.link}
                                href={hive.link}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="max-w-[300px] bg-white rounded-md overflow-hidden border border-gray-200"
                            >
                                <img
                                    src={hive.image}
                                    alt={hive.name}
                                    className="w-full h-48 object-cover"
                                />
                                <div className="p-4">
                                    <h3 className="text-lg font-semibold">{hive.name}</h3>
                                    <p className="text-gray-500 text-sm">{hive.author}</p>
                                    <p className="text-gray-400 text-xs mt-1">
                                        {hive.date.toLocaleDateString()}
                                    </p>
                                </div>
                            </a>
                        ))}
                    </div>
                    <ScrollBar orientation="horizontal" />
                </ScrollArea>
            </div>
            <div>
                <h2 className="text-2xl font-bold mb-4">Community Builds</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {hives.map(({post, hive, user}) => (
                        <a
                            key={post._id.toString()}
                            href={`/posts/${post._id}`}
                            className="bg-white rounded-md overflow-hidden border border-gray-200"
                        >
                            <img
                                src={post.imageUrl}
                                alt={post.title}
                                className="w-full h-48 object-cover"
                            />
                            <div className="p-4">
                                <h3 className="text-lg font-semibold">{post.title}</h3>
                                <p className="text-gray-500 text-sm">By {user.globalName}</p>
                                <p className="text-gray-400 text-xs mt-1">
                                    {new Date(post.createdAt).toLocaleDateString()}
                                </p>
                            </div>
                        </a>
                    ))}
                </div>
            </div>
        </div>
    )
}