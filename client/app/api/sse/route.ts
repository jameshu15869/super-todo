import EventSource from "eventsource";

export const dynamic = 'force-dynamic';

export async function GET(req: Request) {
    let responseStream = new TransformStream();
    const writer = responseStream.writable.getWriter();
    const encoder = new TextEncoder();
    const eventSource = new EventSource(`${process.env.GATEWAY_API_ENDPOINT}/sse`);
    eventSource.onmessage = function(e) {
        writer.write(encoder.encode(`data: ${e.data}\n\n`));
    }

    return new Response(responseStream.readable, {
        headers: {
            "Content-Type" : "text/event-stream",
            "Connection" : "keep-alive",
            "Cache-Control" : "no-cache",
            "Content-Encoding" : "none"
        }
    });
}
