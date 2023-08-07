import { cn } from "@/lib/utils";
import React from "react";
import Add from "./Add/Add";

export default function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <div className="fixed top-0 z-10 w-full bg-white border-b p-3 flex flex-row justify-between">
      <nav>
        <h1 className="text-4xl font-extrabold tracking-tight">Super Todo</h1>
      </nav>
    </div>
  );
}
