import {
    NavigationMenu,
    NavigationMenuItem, NavigationMenuLink,
    NavigationMenuList,
} from "@/components/ui/navigation-menu";
import {getCurrentSession} from "@/lib/server/session";
import {Button} from "@/components/ui/button";
import {ArrowUpRight} from "lucide-react";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";

export const Navbar = async () => {
    const {user} = await getCurrentSession();
    return (
        <div className="flex items-center justify-between px-6 py-3">
            <a href="/" className="flex shrink-0 items-center gap-2">
                <img
                    alt="Logo" className="h-12 w-auto dark:invert"
                    src="https://meta-bee.com/wp-content/uploads/2024/06/metabee-log-new-png.png"/>
            </a>
            <NavigationMenu>
                <div className="relative">
                    <NavigationMenuList>
                        <NavigationMenuItem>
                            <NavigationMenuLink href="/">Home</NavigationMenuLink>
                        </NavigationMenuItem>
                        <NavigationMenuItem>
                            <NavigationMenuLink href="/">Hive Builds</NavigationMenuLink>
                        </NavigationMenuItem>
                        <NavigationMenuItem>
                            <NavigationMenuLink href="https://meta-bee.com" target="_blank" rel="noopener noreferrer">Meta
                                Bee<ArrowUpRight className="ml-1 h-3 w-3"/></NavigationMenuLink>
                        </NavigationMenuItem>
                    </NavigationMenuList>
                </div>
                <div className="absolute top-full left-0 isolate z-50 flex justify-center"></div>
            </NavigationMenu>
            <div className="flex items-center gap-2.5">
                {user !== null ? (
                    <a href="/profile">
                        <Button variant="outline" className="flex items-center gap-2">
                            <Avatar>
                                <AvatarImage
                                    src={user.profileUrl}
                                    alt="Profile Picture"
                                />
                                <AvatarFallback>{user.globalName ? user.globalName[0] : user.username[0]}</AvatarFallback>
                            </Avatar>
                            {user.globalName ?? user.username}
                        </Button>
                    </a>
                ) : (
                    <a href="/login">
                        <Button>Login</Button>
                    </a>
                )}
            </div>
        </div>
    )
};