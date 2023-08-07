import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import React, { useContext } from "react";
import { pb } from "@/types/pb";
import { useCombined } from "@/networking/useCombined";
import { DataContext } from "@/context/DataContext";

interface UserDetailProps {
  user: pb.User;
}

export default function UserDetail({
  user,
  children,
}: UserDetailProps & React.HTMLAttributes<HTMLElement>) {
  const { todos, combined } = useContext(DataContext);

  const filterUserTodos = () => {
    if (combined && todos) {
      const todoIdsToFind = combined
        .filter((combined) => combined.user_id === user.id)
        .map((combined) => combined.todo_id);
      const foundTodos = todos.filter((todo) =>
        todoIdsToFind.includes(todo.id)
      );
      if (todoIdsToFind.length == 0) {
        return <p className="text-md text-muted-foreground">No todos found</p>;
      }
      return (
        <>
          {foundTodos.map((todo) => (
            <div
              key={todo.id}
              className="whitespace-nowrap overflow-hidden text-ellipsis"
            >
              <p className="text-md text-ellipsis overflow-hidden whitespace-nowrap">
                <strong>{todo.title}</strong>: {todo.body}
              </p>
            </div>
          ))}
        </>
      );
    }
    throw new Error("Filtering user todos errored");
  };
  if (!todos || !combined) {
    return <div>Loading...</div>;
  }
  return (
    <Dialog>
      <DialogTrigger>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="underline">
            {user.username + "'s "} Todos
          </DialogTitle>
        </DialogHeader>
        {filterUserTodos()}
      </DialogContent>
    </Dialog>
  );
}
