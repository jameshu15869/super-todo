import { pb } from "@/types/pb";
import { NextResponse } from "next/server";

export async function GET(req: Request) {
    return NextResponse.json({msg : "add user"});
}

export async function POST(req: Request) {
    const body = await req.json();
    if (!(body.title && body.date && body.body && body.user_ids)) {
        return new Response("Bad input provided", {
            status: 400
        });
    }

    const putTodoRes = await fetch(`${process.env.GATEWAY_API_ENDPOINT}/todos/add`, {
        method: "POST",
        body: JSON.stringify({
            title: body.title,
            todo_date: body.date,
            body: body.body,
            user_ids: body.user_ids
        })
    });

    const todoData = await putTodoRes.json();
    if (!putTodoRes.ok) {
        return new Response(todoData.message, {
            status: putTodoRes.status
        });
    }

    const putCombinedRes = await fetch(`${process.env.GATEWAY_API_ENDPOINT}/combined/addarray`, {
        method: "POST",
        body: JSON.stringify({
            todo_id : todoData.todo.id,
            user_ids: body.user_ids
        })
    });
    const combinedData = await putCombinedRes.json();
    if (!putCombinedRes.ok) {
        return new Response(combinedData.message, {
            status: putCombinedRes.status
        });
    }

    return NextResponse.json({
        todo: todoData.todo,
        combined: combinedData.combined
    });
}