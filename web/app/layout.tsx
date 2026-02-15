import type {Metadata} from "next";
import {Geist, Geist_Mono, Inter} from "next/font/google";
import "./globals.css";
import {Navbar} from "@/components/navbar";

const inter = Inter({subsets: ['latin'], variable: '--font-sans'});

const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "Hive Builder",
    description: "Best Bee Swarm Simulator hive builder!",
};

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" className={inter.variable}>
        <body
            className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        >
            <div className="flex flex-col min-h-screen">
                <Navbar/>
                {children}
            </div>
        </body>
        </html>
    );
}
