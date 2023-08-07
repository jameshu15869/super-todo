import { z } from "zod";
import TodoEditForm from "../TodoCard/TodoEditForm";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { useContext, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { pb } from "@/types/pb";
import { Form, FormControl, FormField, FormItem, FormLabel } from "../ui/form";
import { Input } from "../ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { format } from "date-fns";
import { CalendarIcon, Check, ChevronsUpDown } from "lucide-react";
import { Calendar } from "../ui/calendar";
import TodoPeople from "../TodoCard/TodoPeople";
import {
  Command,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "../ui/command";
import { useUsers } from "@/networking/useUsers";
import { cn } from "@/lib/utils";
import { DataContext } from "@/context/DataContext";

const formSchema = z.object({
  title: z.string().min(1),
  date: z.date(),
  body: z.string().min(1),
  user_ids: z.array(z.number()).min(1),
});

interface AddTodoFormProps {
  allUsers: pb.User[];
}

function AddTodoForm({ allUsers }: AddTodoFormProps) {
  const [modifiedUsers, setModifiedUsers] = useState<pb.User[]>([]);
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState<number>();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: "",
      date: new Date(Date.now()),
      body: "",
      user_ids: [],
    },
  });

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    // console.log(values);
  };

  const handleAddUser = (user: pb.User) => {
    const copy = modifiedUsers.slice();
    copy.push(user);
    form.setValue(
      "user_ids",
      copy.map((user) => user.id)
    );
    setModifiedUsers(copy);
  };

  const handleDeleteUser = (user_id: number) => {
    const filtered = modifiedUsers.filter((user) => user_id != user.id);
    setModifiedUsers(filtered);
    form.setValue(
      "user_ids",
      filtered.map((user) => user.id)
    );
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col space-y-8"
      >
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Title</FormLabel>
              <FormControl>
                <Input placeholder="Title" {...field} />
              </FormControl>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="date"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Date of birth</FormLabel>
              <Popover>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant={"outline"}
                      className={cn(
                        "w-[240px] pl-3 text-left font-normal",
                        !field.value && "text-muted-foreground"
                      )}
                    >
                      {field.value ? (
                        format(new Date(field.value), "PPP")
                      ) : (
                        <span>Pick a date</span>
                      )}
                      <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value}
                    // @ts-ignore
                    onSelect={field.onChange}
                    disabled={(date) => date < new Date("1900-01-01")}
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="body"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Body</FormLabel>
              <FormControl>
                <Input placeholder="Feed the dog..." {...field} />
              </FormControl>
            </FormItem>
          )}
        />

        <div className="flex flex-col">
          <TodoPeople
            className=""
            users={modifiedUsers}
            canEdit={true}
            handleDelete={handleDeleteUser}
          />
          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                aria-expanded={open}
                className="mx-auto mt-1.5 w-[200px] justify-between"
              >
                {value
                  ? allUsers!.find((user) => user.id === value)?.username
                  : "Add a user..."}
                <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[200px] p-0">
              <Command>
                <CommandInput placeholder="Search user..." />
                <CommandGroup>
                  {allUsers.map((user) => {
                    if (
                      !modifiedUsers.find(
                        (includedUser) => includedUser.id === user.id
                      )
                    ) {
                      return (
                        <CommandItem
                          key={user.id}
                          onSelect={(currentValue) => {
                            // console.log("Hi: ", currentValue, allUsers);
                            // console.log(
                            //   "Searched: ",
                            //   allUsers!.find(
                            //     (user) =>
                            //       user.username.toLowerCase() ===
                            //       currentValue.toLowerCase()
                            //   )
                            // );
                            const foundUser = allUsers.find(
                              (user) =>
                                user.username.toLowerCase() ===
                                currentValue.toLowerCase()
                            )!;
                            handleAddUser(foundUser);
                            setValue(undefined);
                            setOpen(false);
                          }}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              value === user.id ? "opacity-100" : "opacity-0"
                            )}
                          />
                          {user.username}
                        </CommandItem>
                      );
                    }
                  })}
                </CommandGroup>
              </Command>
            </PopoverContent>
          </Popover>
        </div>

        <Button type="submit" className="">
          Submit
        </Button>
      </form>
    </Form>
  );
}

export default function AddTodo() {
  // const { users, error, isLoading } = useUsers();
  const { users } = useContext(DataContext);
  const [open, setOpen] = useState(false);

  const closeForm = () => {
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" className="w-full">
          Add Todo
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Todo</DialogTitle>
          <DialogDescription>Add a new todo</DialogDescription>
        </DialogHeader>
        {/* <AddTodoForm allUsers={users!} /> */}
        <TodoEditForm
          todo={new pb.Todo(-1, "", new Date(Date.now()), "")}
          users={[]}
          allUsers={users!}
          closeForm={closeForm}
        />
      </DialogContent>
    </Dialog>
  );
}
