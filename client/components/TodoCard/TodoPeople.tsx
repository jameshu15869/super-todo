import { useState } from "react";
import { Button } from "../ui/button";
import DialogPerson from "./TodoPerson";
import TodoPerson from "./TodoPerson";
import { pb } from "@/types/pb";
import { cn } from "@/lib/utils";

interface TodoPeopleProps {
  users: pb.User[];
  canEdit: boolean;
  handleDelete?: (user_id: number) => void;
}

export default function TodoPeople({
  users,
  canEdit,
  handleDelete,
  className,
}: TodoPeopleProps & React.HTMLAttributes<HTMLElement>) {
  return (
    <div className={className}>
      <h2
        className={cn(
          "text-sm w-full",
          users.length == 0 ? "text-destructive" : ""
        )}
      >
        Assigned Users
      </h2>
      {users.length == 0 && (
        <p className="text-sm text-muted-foreground">
          Please select at least one user.
        </p>
      )}
      <div className="mt-1 grid gap-1">
        {users.map((user, index) => (
          <TodoPerson
            key={index}
            user={user}
            canEdit={canEdit}
            handleDelete={handleDelete}
          />
        ))}
      </div>
    </div>
  );
}
