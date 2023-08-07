import { NextResponse } from "next/server";

export async function GET(req: Request, {params} : {params: {id: number}}) {
    return new Response(`Update id #${params.id}`);
}

export async function POST(req: Request, {params} : {params: {id: number}}) {
    const body = await req.json();
    if (!(body.title && body.todo_date && body.body && body.user_ids)) {
        return new Response("Bad input provided", {
            status: 400
        });
    }

    const updateTodo = await fetch(`${process.env.GATEWAY_API_ENDPOINT}/todos/${params.id}/update`, {
        method: "POST",
        body: JSON.stringify({
            title: body.title,
            todo_date: body.todo_date,
            body: body.body,
            user_ids: body.user_ids
        })
    });

    const todoData = await updateTodo.json();
    if (!updateTodo.ok) {
        return new Response(todoData.message, {
            status: updateTodo.status
        });
    }

    return NextResponse.json({todo: todoData.todo});
}