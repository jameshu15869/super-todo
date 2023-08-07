import { NextResponse } from "next/server";

export async function GET(req: Request, {params} : {params: {id: number}}) {
    return new Response(`Delete id #${params.id}`);
}

export async function POST(req: Request, {params} : {params: {id: number}}) {
    const updateTodo = await fetch(`${process.env.GATEWAY_API_ENDPOINT}/todos/${params.id}/delete`, {
        method: "POST",
    });

    const todoData = await updateTodo.json();
    if (!updateTodo.ok) {
        return new Response(todoData.message, {
            status: updateTodo.status
        });
    }

    return NextResponse.json({todo: todoData.todo});
}