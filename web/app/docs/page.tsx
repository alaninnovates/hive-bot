import {
    Sidebar,
    SidebarContent,
    SidebarGroup, SidebarGroupContent, SidebarGroupLabel,
    SidebarMenu, SidebarMenuButton,
    SidebarMenuItem, SidebarProvider, SidebarTrigger, SidebarWrapper
} from "@/components/ui/sidebar";
import {CircleQuestionMark, Command, Home} from "lucide-react";
import Link from "next/link";
import {Button} from "@/components/ui/button";
import {commands} from "@/app/docs/data";
import {Kbd} from "@/components/ui/kbd";

const DocsSidebar = (props: React.HTMLAttributes<HTMLDivElement>) => {
    return (
        <Sidebar {...props}>
            <SidebarContent>
                <SidebarMenu>
                    <SidebarGroup>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild>
                                <Link href="#home">
                                    <Home className="mr-2 h-4 w-4"/>
                                    <span>Home</span>
                                </Link>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarGroup>
                    <SidebarGroup>
                        <SidebarGroupLabel>
                            <Command className="mr-2 h-4 w-4"/>
                            <span>Commands</span>
                        </SidebarGroupLabel>
                        <SidebarGroupContent>
                            <SidebarMenuItem>
                                {commands.map((command) => (
                                    <SidebarMenuButton asChild key={command.name}>
                                        <Link href={`#${command.name}`}>
                                            <span>{command.name}</span>
                                        </Link>
                                    </SidebarMenuButton>
                                ))}
                            </SidebarMenuItem>
                        </SidebarGroupContent>
                    </SidebarGroup>
                    <SidebarGroup>
                        <SidebarGroupLabel>
                            <CircleQuestionMark className="mr-2 h-4 w-4"/>
                            <span>FAQ</span>
                        </SidebarGroupLabel>
                        <SidebarGroupContent>
                            <SidebarMenuItem>
                                <SidebarMenuButton asChild>
                                    <Link href="#hive_slots">
                                        <span>How to: Hive Slots</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarGroupContent>
                    </SidebarGroup>
                </SidebarMenu>
            </SidebarContent>
        </Sidebar>
    )
}

export default async function Page() {
    return (
        <SidebarProvider className="container mx-auto px-4 grow">
            <div className="overflow-auto w-full min-h-full">
                <SidebarWrapper className="relative h-full min-h-auto">
                    <DocsSidebar className="absolute h-full"/>
                    <main className="p-4 pl-8 w-full overflow-y-auto space-y-8">
                        <div id="home" className="space-y-4">
                            <h1 className="text-4xl font-bold">Hive Builder Documentation</h1>
                            <p>
                                Hive Builder Bot is a tool that allows users to easily create and customize virtual
                                hives in
                                the game Bee Swarm Simulator within the Discord platform. The bot allows users to add
                                bees,
                                mutations, and beequips to their hives and also provides guides and resources to help
                                players improve their gameplay.
                            </p>
                            <p>
                                Hive Builder Bot was developed by alaninnovates#0123 and is specifically designed for
                                use in
                                Discord.
                            </p>
                            <div className="flex flex-wrap items-center justify-center gap-4">
                                <Button size="lg">
                                    Invite Hive Builder Bot
                                </Button>
                                <Button size="lg">
                                    Join The Meta Bee Discord
                                </Button>
                            </div>
                        </div>
                        <div id="commands" className="space-y-4">
                            <h2 className="text-2xl font-bold">Commands</h2>
                            <h3 className="text-xl font-semibold">Key</h3>
                            <p>
                                <Kbd>&lt; &gt;</Kbd> = Required
                                <br/>
                                <Kbd>[ ]</Kbd> = Optional
                            </p>
                            <div className="flex flex-wrap items-center justify-center gap-4">
                                {commands.map((command) => (
                                    <div key={command.name} id={command.name} className="w-full">
                                        <h4 className="text-lg font-semibold">{command.name}
                                            <Kbd>{command.command}</Kbd></h4>
                                        <p className="text-gray-600">{command.description}</p>
                                        {command.arguments.length > 0 && (
                                            <div className="mt-2">
                                                <h5 className="font-medium">Arguments:</h5>
                                                <ul className="list-disc list-inside">
                                                    {command.arguments.map((arg) => (
                                                        <li key={arg.name}>
                                                            <span
                                                                className="font-semibold">{arg.name}</span>: {arg.description}
                                                        </li>
                                                    ))}
                                                </ul>
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>
                        <div id="faq" className="space-y-4">
                            <h2 className="text-2xl font-bold">FAQ</h2>
                            <div id="hive_slots" className="space-y-2">
                                <h3 className="text-xl font-semibold">How to: Hive Slots</h3>
                                <p>
                                    When running the commands above, you may have run into a text argument called
                                    "slots". This argument tells the bot which hive slots you would like to modify. It
                                    is a very flexible argument and can be very powerful if you know how to use it.
                                </p>
                                <h4 className="font-medium">Specifying a single hive slot:</h4>
                                <p>
                                    To select any hive slot, type in any number between 1-50. These numbers represent
                                    which slot you want to place the bee in.
                                </p>
                                <h4 className="font-medium">Specifying multiple non-consecutive hive slots:</h4>
                                <p>
                                    To select any amount of hive slots, type in the slot numbers, with each slot having
                                    a comma after it. For example, if I wanted to select slots 18, 23, 28, and 33, I
                                    would type 18,23,28,33 for the slots argument. Make sure there are no spaces!
                                </p>
                                <h4 className="font-medium">Specifying a range of consecutive hive slots:</h4>
                                <p>
                                    To select a range of hive slots, type in the starting slot number and end slot
                                    number with a dash between them. For example, if I wanted to select all the slots in
                                    my hive, I would type 1-50 for the slots argument. Once again, make sure there are
                                    no spaces!
                                </p>
                                <p>
                                    You can combine these arguments to create a statement for many slots at once. For
                                    example, I could type in 1-20,23,34,40-45,51 and I would be able to select all of
                                    those hive slots!
                                </p>
                            </div>
                        </div>
                    </main>
                    <SidebarTrigger/>
                </SidebarWrapper>
            </div>
        </SidebarProvider>
    )
}