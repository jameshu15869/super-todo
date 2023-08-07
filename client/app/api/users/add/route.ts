import { pb } from "@/types/pb";
import { NextResponse } from "next/server";

export async function GET(req: Request) {
    return NextResponse.json({msg : "hi"});
}

export async function POST(req: Request) {
    const body = await req.json();
    if (!body.username) {
        return new Response("No username provided", {
            status: 500
        });
    }

    const res = await fetch(`${process.env.GATEWAY_API_ENDPOINT}/users/add`, {
        method: "POST",
        body: JSON.stringify({
            username: body.username
        })
    });

    const jsonData = await res.json();
    if (!jsonData.ok) {
        return new Response(jsonData.message, {
            status: res.status
        });
    }
    return NextResponse.json({ body });
}