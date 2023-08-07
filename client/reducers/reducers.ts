import { pb } from "@/types/pb";

interface IAction {
    type: string;
}

interface TodoAction extends IAction {
    todo: pb.Todo
}

type Action = TodoAction;

export function todoReducer(todos: pb.Todo[], action: Action): pb.Todo[] | undefined {
    switch (action.type) {
        case 'add-todo':
            const newTodo = action.todo;
            return [
                ...todos,
                newTodo
            ];
    }

    return undefined;
}