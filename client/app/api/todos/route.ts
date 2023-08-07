import { pb } from "@/types/pb";
import { NextResponse } from "next/server";

export async function GET(req: Request) {
    const res = await fetch(`${process.env.GATEWAY_API_ENDPOINT}/todos`, { cache: 'no-store' });
    if(!res.ok) {
      console.log(res);
      throw new Error("Failed to fetch todos data");
    }
    const jsonData = await res.json();
    return NextResponse.json({data: jsonData.todos});
}