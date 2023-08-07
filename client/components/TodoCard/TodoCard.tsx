"use client";

import { supertodo } from "@/pb/super-todo";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Button } from "../ui/button";
import { Expand } from "lucide-react";
import React, { useLayoutEffect, useRef, useState } from "react";
import { cn } from "@/lib/utils";
import TodoCardDialog from "./TodoCardDialog";
import { Separator } from "../ui/separator";
import TodoPeople from "./TodoPeople";
import { pb } from "@/types/pb";
import TodoCardPeopleSummary from "./TodoCardPeopleSummary";

const ROW_SPAN_BREAKPOINT = 384; /* tailwind-96 = 384px = 12rem */
interface TodoCardProps {
  todo: pb.Todo;
  users: pb.User[];
  allUsers: pb.User[];
}

export default function TodoCard({
  className,
  children,
  todo,
  users,
  allUsers,
  ...props
}: React.HTMLAttributes<HTMLElement> & TodoCardProps) {
  const [expanded, setExpanded] = useState(false);
  const cardRef = useRef<HTMLDivElement>(null);
  const [rowSpan, setRowSpan] = useState(0);

  useLayoutEffect(() => {
    if (cardRef.current) {
      const { height } = cardRef.current.getBoundingClientRect();
      setRowSpan(Math.ceil(cardRef.current.scrollHeight / ROW_SPAN_BREAKPOINT));
    }
  }, []);

  return (
    <Card
      className={cn(
        `relative m-2 flex flex-col group/card`,
        rowSpan > 1 ? `row-span-${rowSpan}` : ``,
        className
      )}
      ref={cardRef}
    >
      {/* <Button variant="ghost" className="absolute right-0 mt-1">
        <Expand />
      </Button> */}
      <CardHeader className="relative">
        <CardTitle>{todo.title}</CardTitle>
        <CardDescription>
          {new Date(todo.todo_date).toLocaleDateString("en-us")}
        </CardDescription>
        <TodoCardDialog
          className="absolute right-1 top-0"
          users={users}
          todo={todo}
          allUsers={allUsers}
        />
      </CardHeader>
      <CardContent className="flex-grow">
        {/* <div className="max-h-20 text-ellipsis overflow-hidden whitespace-nowrap"> */}
        <div className="">
          {/* Lorem, ipsum dolor sit amet consectetur adipisicing elit. Obcaecati
          dolor corrupti illo modi voluptate, corporis et vero! Dolore facilis
          laboriosam error quia eveniet fuga, atque et eaque repellat veritatis
          obcaecati. Quod adipisci fuga deserunt quo, neque excepturi soluta
          ullam architecto sed deleniti omnis, sint cumque quia autem a
          exercitationem tenetur! */}
          {todo.body}
        </div>
      </CardContent>
      <CardFooter>
        <TodoCardPeopleSummary
          className="w-full justify-end gap-1"
          users={users}
        />
        {/* <small
          className="text-slate-200 group-hover/card:text-slate-400 text-xs font-medium ml-auto text-muted-foreground hover:underline hover:cursor-pointer"
          onClick={toggleExpanded}
        >
          {expanded ? "Show less" : "Show more"}
        </small> */}
        {/* <Button
          variant="ghost"
          className="absolute right-0"
          onClick={toggleExpanded}
        >
          <Expand />
        </Button> */}
      </CardFooter>
    </Card>
  );
}
