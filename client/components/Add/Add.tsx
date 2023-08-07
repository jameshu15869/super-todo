import { Plus } from "lucide-react";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import React, { useState } from "react";
import { cn } from "@/lib/utils";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { PopoverArrow } from "@radix-ui/react-popover";
import AddUser from "./AddUser";
import AddTodo from "./AddTodo";

export default function Add({ className }: React.HTMLAttributes<HTMLElement>) {
  const [open, setOpen] = useState(false);
  const closePopover = () => {
    setOpen(false);
  };
  return (
    // <Dialog>
    //   <DialogTrigger>
    //     <Button className={cn("rounded-full", className)}>
    //       <Plus />
    //     </Button>
    //   </DialogTrigger>
    //   <DialogContent>
    //     <DialogTitle>Add User</DialogTitle>
    //     <DialogDescription>Add a user</DialogDescription>
    //   </DialogContent>
    // </Dialog>
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button className={cn("rounded-full", className)}>
          <Plus />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-36 p-0" sideOffset={3}>
        <div className="flex flex-col p-0">
          <AddUser closePopover={closePopover} />
          <AddTodo />
        </div>
        <PopoverArrow fill="white" />
      </PopoverContent>
    </Popover>
  );
}
