export const commands = [
    {
        name: "Create Hive",
        command: "/hive create",
        description: "This is the first command you run before starting to build your hive. This initializes your profile, so you can start adding bees. Additionally, if you want to reset your whole hive, run this command to start with a clean slate.",
        arguments: [],
    },
    {
        name: "Add Bee",
        command: "/hive add <name> <slots> <level> [gifted]",
        description: "Add bee(s) to your hive.",
        arguments: [
            {
                name: "name",
                description: "The name of the bee to add. You can select from the list or start typing to filter out names.",
            },
            {
                name: "slots",
                description: "Refer to How to – Slots",
            },
            {
                name: "level",
                description: "Type in a number from 0-25, which is the level of your bee. Level 0 means that no level will be displayed.",
            },
            {
                name: "gifted",
                description: "Whether or not the bee should be gifted",
            },
        ],
    },
    {
        name: "Remove Bee",
        command: "/hive remove <slots>",
        description: "Remove bee(s) from your hive.",
        arguments: [
            {
                name: "slots",
                description: "Refer to \"How to Slots\" at the end of documentation.",
            },
        ],
    },
    {
        name: "Set Mutation",
        command: "/hive setmutation <slots> <name>",
        description: "Set the mutation of a bee within the slots you specified.",
        arguments: [
            {
                name: "slots",
                description: "Refer to \"How to Slots\" at the end of documentation.",
            },
            {
                name: "name",
                description: "The name of the mutation to set. You can select from the list or start typing to filter out names.",
            },
        ],
    },
    {
        name: "Set Beequip",
        command: "/hive setbeequip <slots> <name>",
        description: "Set the beequip of a bee within the slots you specified.",
        arguments: [
            {
                name: "slots",
                description: "Refer to \"How to Slots\" at the end of documentation.",
            },
            {
                name: "name",
                description: "The name of the beequip to set. You can select from the list or start typing to filter out names.",
            },
        ],
    },
    {
        name: "Gift All",
        command: "/hive giftall",
        description: "Gift all of the bees in your hive.",
        arguments: [],
    },
    {
        name: "Set Level",
        command: "/hive setlevel <level>",
        description: "Set the level of ALL of the bees in your hive. Note: This is irreversible.",
        arguments: [
            {
                name: "level",
                description: "Type in a number from 0-25, which is the level to set all of your bees to. Level 0 means that no level will be displayed.",
            },
        ],
    },
    {
        name: "View Hive",
        command: "/hive view [show_hive_numbers] [slots_on_top]",
        description: "View your hive.",
        arguments: [
            {
                name: "show_hive_numbers",
                description: "Show the numbers on each hive slot. These numbers are to help you figure out what slot you want to add a bee to, but you might not want them when showing the hive to others. Set this to False to disable the numbers.",
            },
            {
                name: "slots_on_top",
                description: "When show_hive_numbers is enabled, this parameter allows you to draw hive numbers above the faces of any bee in your hive. This allows you to easily replace a bee, even if your hive is full!",
            },
        ],
    },
    {
        name: "Save Hive",
        command: "/hive save <name>",
        description: "Save your hive, so you can always load it up again. If you want to overwrite a previous save, just use the same name as the old save. Be sure to ALWAYS SAVE YOUR HIVE! This way, you won’t lose any data.",
        arguments: [
            {
                name: "name",
                description: "The name of your save",
            },
        ],
    },
    {
        name: "Hive Info",
        command: "/hive info",
        description: "Get a summary of your hive and its contents.",
        arguments: [],
    },
    {
        name: "List Saves",
        command: "/hive saves list",
        description: "List all of your hives.",
        arguments: [],
    },
    {
        name: "Load Hive",
        command: "/hive saves load <id>",
        description: "Load a previously saved hive. Note: THIS WILL OVERWRITE YOUR CURRENT HIVE! If you have data in your current working hive that you want to save, be sure to save it first!",
        arguments: [
            {
                name: "id",
                description: "The ID of the saved hive. You can obtain this from the /hive saves list command.",
            },
        ],
    },
    {
        name: "Delete Hive",
        command: "/hive saves delete <id>",
        description: "Delete a previously saved hive. Note: THIS IS IRREVERSIBLE!",
        arguments: [
            {
                name: "id",
                description: "The ID of the saved hive. You can obtain this from the /hive saves list command.",
            },
        ],
    },
    {
        name: "Games",
        command: "/game <game>",
        description: "All the subcommands within this command are games that you can play for fun.",
        arguments: [
            {
                name: "game",
                description: "The name of the game to play. You can select from the list or start typing to filter out names.",
            },
        ],
    }
];