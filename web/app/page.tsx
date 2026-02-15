import {Button} from "@/components/ui/button";
import {ArrowUpRight} from "lucide-react";

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
                        <Button variant="outline">Browse Hives</Button>
                        <Button>Invite <ArrowUpRight className="ml-1 h-3 w-3"/></Button>
                    </div>
                </div>
                <div className="md:w-1/2 mt-8 md:mt-0">
                    <img src="/" alt="pic" className="w-full"/>
                </div>
            </div>
        </section>
    )
}
