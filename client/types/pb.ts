export namespace pb {
    export interface ServerResponse {
        data: Todo;
    }

    export interface PbData {
        todos: Todo[] | undefined;
        users: User[] | undefined;
        combined: Combined[] | undefined;
    }

    export class Todo {
        id: number;
        title: string;
        todo_date: Date;
        body: string;
        constructor(id: number, title: string, date: Date, body: string) {
            this.id = id;
            this.title = title;
            this.todo_date = date;
            this.body = body;
        }
    }

    export class User {
        id: number;
        username: string;
        constructor(id: number, username: string) {
            this.id = id;
            this.username = username;
        }
    }

    export class Combined {
        id: number;
        user_id: number;
        todo_id: number;
        constructor(id: number, user_id: number, todo_id: number) {
            this.id = id;
            this.user_id = user_id;
            this.todo_id = todo_id;
        }
    }
}