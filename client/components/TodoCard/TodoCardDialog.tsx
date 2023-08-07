"use client";

import { Edit, Expand, MenuSquare, User2, X } from "lucide-react";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Separator } from "../ui/separator";
import { Avatar, AvatarFallback } from "../ui/avatar";
import React, { useState } from "react";
import TodoPeople from "./TodoPeople";
import { pb } from "@/types/pb";
import TodoEditForm from "./TodoEditForm";

interface TodoCardDialogSelfProps {
  todo: pb.Todo;
  users: pb.User[];
  allUsers: pb.User[];
}

type TodoCardDialogProps = TodoCardDialogSelfProps &
  React.HTMLAttributes<HTMLElement>;

export default function TodoCardDialog({
  className,
  todo,
  users,
  allUsers,
  ...props
}: TodoCardDialogProps) {
  const [open, setOpen] = useState(false);

  const closeForm = () => {
    setOpen(false);
  };
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild className={className}>
        {/* <small className="text-slate-200 group-hover/card:text-slate-400 text-xs font-medium ml-auto text-muted-foreground hover:underline hover:cursor-pointer"> */}
        {/* <small
          className="text-slate-200 group-hover/card:text-slate-400 
          text-xs font-medium ml-auto text-muted-foreground 
          group-hover/card:cursor-pointer 
          group-hover/card:underline"
        >
          Show more
        </small> */}

        {/* <Button variant="ghost">
          <small className="text-sm">Show more</small>
        </Button> */}

        <Button variant="ghost" className="px-2">
          <Edit />
        </Button>
      </DialogTrigger>
      <DialogContent className="mt-0">
        <DialogHeader>
          <DialogTitle>Edit Todo</DialogTitle>
        </DialogHeader>
        <TodoEditForm
          todo={todo}
          users={users}
          allUsers={allUsers}
          closeForm={closeForm}
        />
      </DialogContent>
    </Dialog>
  );
}
