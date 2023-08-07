"use client";

import TodoCard from "./TodoCard/TodoCard";
import Add from "./Add/Add";
import React, { useEffect } from "react";
import { cn } from "@/lib/utils";
import { Toaster } from "./ui/toaster";
import { useSSE } from "@/networking/useSSE";
import { DataContext } from "@/context/DataContext";
import Spinner from "./Spinner/Spinner";

export default function MainContainer({
  className,
}: React.HTMLAttributes<HTMLElement>) {
  const {
    todos: dynamicTodos,
    combined: dynamicCombined,
    users: dynamicUsers,
    isLoading: mainLoad,
  } = useSSE();
  
  if (mainLoad) {
    return (
      <div className="w-screen h-screen flex flex-col align-center justify-center">
        <div className="flex flex-row justify-center gap-2 -ml-7">
          <p className="text-lg">Loading</p>
          <Spinner isLoading={true} />
        </div>
      </div>
    );
  }

  const getUsersFromCombined = (todo_id: number) => {
    if (dynamicUsers && dynamicCombined) {
      const userIds = dynamicCombined
        .filter((combine) => combine.todo_id == todo_id)
        .map((combine) => combine.user_id);
      const users = dynamicUsers.filter((user) => userIds.includes(user.id));
      return users;
    }

    return [];
  };

  return (
    <div className={cn(className)}>
      <DataContext.Provider
        value={{
          todos: dynamicTodos,
          combined: dynamicCombined,
          users: dynamicUsers,
        }}
      >
        <Add className="fixed z-10 bottom-4 left-4 w-16 h-16" />
        <div className="grid h-full p-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 auto-rows-[24rem]">
          {dynamicTodos!.map((todo) => (
            <TodoCard
              key={todo.id}
              todo={todo}
              users={getUsersFromCombined(todo.id)}
              allUsers={dynamicUsers!}
            />
          ))}
        </div>
      </DataContext.Provider>
      <Toaster />
    </div>
  );
}
