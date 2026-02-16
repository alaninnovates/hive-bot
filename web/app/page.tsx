import {Button} from "@/components/ui/button";
import {ArrowUpRight} from "lucide-react";
import Link from "next/link";
import {
    DiscordActionRow,
    DiscordAttachments,
    DiscordButton,
    DiscordCommand,
    DiscordEmbed,
    DiscordMessage,
    DiscordMessages,
} from "@skyra/discord-components-react";

export default async function Page() {
    return (
        <section className="container mx-auto py-12 px-4 grow flex items-center">
            <div className="flex flex-col md:flex-row items-center">
                <div className="md:w-1/2">
                    <h1 className="text-4xl font-bold mb-4">Hive Builder</h1>
                    <p className="text-lg text-gray-600 mb-6">
                        The best Bee Swarm Simulator hive builder! Create
                        guide-quality hive images and
                        share your hives with the community.
                    </p>
                    <div className="flex flex-col md:flex-row items-center gap-4">
                        <Link href="/posts">
                            <Button variant="outline">Browse Hives</Button>
                        </Link>
                        <Link href="/docs">
                            <Button>Invite <ArrowUpRight className="ml-1 h-3 w-3"/></Button>
                        </Link>
                    </div>
                </div>
                <div className="md:w-1/2 mt-8 md:mt-0">
                    <div className="w-full rounded-md p-4 overflow-auto bg-[#36393E]">
                        <DiscordMessages className="w-full h-full">
                            <DiscordMessage profile="Hive Builder"
                                            author="Hive Builder"
                                            avatar="https://cdn.discordapp.com/app-icons/1051308449172049970/2200fff798c868556892fae4981d6acb.png"
                            >
                                <DiscordCommand slot="reply"
                                                profile="alaninnovates"
                                                command="/hive view"
                                                author="alaninnovates"
                                />
                                <DiscordEmbed slot="embeds" embedTitle="alaninnovates's Hive"
                                              image="https://meta-bee.com/wp-content/uploads/2023/08/diamond-30.3.png"
                                              color="#3B3B40"/>
                                <DiscordAttachments slot="components">
                                    <DiscordActionRow>
                                        <DiscordButton type="primary">Add Bee</DiscordButton>
                                        <DiscordButton type="primary">Gift All</DiscordButton>
                                        <DiscordButton type="primary">Set Level All</DiscordButton>
                                        <DiscordButton type="success">Hive Info</DiscordButton>
                                        <DiscordButton type="secondary">Rerender Hive</DiscordButton>
                                    </DiscordActionRow>
                                </DiscordAttachments>
                            </DiscordMessage>
                        </DiscordMessages>
                    </div>
                </div>
            </div>
        </section>
    )
}
