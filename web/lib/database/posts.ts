import {getDb} from "@/lib/server/db";
import {User} from "@/lib/server/user";
import {ObjectId, WithId} from "mongodb";
import {Long} from "bson";

export const getHivePosts = async ({postId}: { postId?: string; } = {}): Promise<{
    post: WithId<Post>;
    hive: Hive;
    user: User;
}[]> => {
    const db = await getDb();
    const posts = await db.collection("posts").aggregate([
        ...(postId ? [{$match: {_id: new ObjectId(postId)}}] : []),
        {
            $lookup: {
                from: "hives",
                localField: "hive_id",
                foreignField: "_id",
                as: "hive"
            },
        },
        {$unwind: "$hive"},
        {
            $lookup: {
                from: "users",
                let: {userId: "$hive.user_id"},
                pipeline: [
                    {
                        $match: {
                            $expr: {
                                $eq: [
                                    {$toLong: "$discord_id"},
                                    "$$userId"
                                ]
                            }
                        }
                    }
                ],
                as: "user"
            }
        },
        {$unwind: "$user"},
        {
            $project: {
                "post._id": "$_id",
                "post.title": "$title",
                "post.content": "$content",
                "post.createdAt": "$created_at",
                "post.hiveId": "$hive_id",
                "hive.name": "$hive.name",
                "hive.userId": "$hive.user_id",
                "hive.bees": "$hive.bees",
                "user.discordId": "$user.discord_id",
                "user.email": "$user.email",
                "user.username": "$user.username",
                "user.globalName": "$user.global_name",
                "user.profileUrl": "$user.profile_url",
            }
        }
    ]).toArray();
    return posts as {
        post: WithId<Post>;
        hive: Hive;
        user: User;
    }[];
};

interface Post {
    title: string;
    content: string;
    createdAt: Date;
    hiveId: ObjectId;
}

interface Hive {
    name: string;
    userId: Long;
    bees: {
        [key: number]: {
            name: string;
            level: number;
            gifted: boolean;
            beequip: string;
            mutation: string;
        }
    };
}
