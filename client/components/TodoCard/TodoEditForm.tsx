import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { pb } from "@/types/pb";
import { format } from "date-fns";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Separator } from "../ui/separator";
import TodoPeople from "./TodoPeople";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { cn } from "@/lib/utils";
import { CalendarIcon, Check, ChevronsUpDown } from "lucide-react";
import { Calendar } from "../ui/calendar";
import { useEffect, useState } from "react";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "../ui/command";
import { Textarea } from "../ui/textarea";
import Spinner from "../Spinner/Spinner";
import { useToast } from "../ui/use-toast";

const formSchema = z.object({
  title: z.string().min(1),
  date: z.date(),
  body: z.string().min(1),
  user_ids: z.array(z.number()).min(1),
});

interface TodoEditFormProps {
  todo: pb.Todo;
  users: pb.User[];
  allUsers: pb.User[];
  closeForm: () => void;
}

export default function TodoEditForm({
  todo,
  users,
  allUsers,
  closeForm,
}: TodoEditFormProps) {
  const [modifiedUsers, setModifiedUsers] = useState(users);
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState<number>();
  const [isLoading, setIsLoading] = useState(false);
  const [deleteIsLoading, setDeleteIsLoading] = useState(false);
  const [needConfirmation, setNeedConfirmation] = useState(false);
  const { toast } = useToast();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: todo.title,
      date: new Date(todo.todo_date),
      body: todo.body,
      user_ids: modifiedUsers.map((user) => user.id),
    },
  });

  const deleteTodo = () => {
    setNeedConfirmation(true);
  };

  const stopDelete = () => {
    setNeedConfirmation(false);
  };

  const confirmDelete = () => {
    setDeleteIsLoading(true);
    setNeedConfirmation(false);
    fetch(`/api/todos/${todo.id}/delete`, {
      method: "POST",
    })
      .then((res) => {
        toast({
          description: "Todo deleted successfully!",
        });
        setDeleteIsLoading(false);
        closeForm();
      })
      .catch((err) => {
        toast({
          variant: "destructive",
          title: "Uh oh! Something went wrong.",
          description: "Could not delete",
        });
        setDeleteIsLoading(false);
        closeForm();
      });
  };

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    setIsLoading(true);
    if (todo.id == -1) {
      fetch("/api/todos/add", {
        method: "POST",
        body: JSON.stringify({
          title: values.title,
          date: values.date,
          body: values.body,
          user_ids: modifiedUsers.map((user) => user.id),
        }),
      })
        .then((res) => {
          setIsLoading(false);
          if (res.status === 400) {
            console.log(res);
            throw new Error("Add post failed");
          }

          setOpen(false);
          closeForm();
          toast({
            description: "Todo successfully added!",
          });
        })
        .catch((err) => {
          setIsLoading(false);
          throw new Error(err);
        });
    } else {
      fetch(`/api/todos/${todo.id}/update`, {
        method: "POST",
        body: JSON.stringify({
          title: values.title,
          todo_date: values.date,
          body: values.body,
          user_ids: modifiedUsers.map((user) => user.id),
        }),
      })
        .then((res) => {
          setIsLoading(false);
          if (res.status === 400) {
            console.log(res);
            throw new Error("Update POST failed");
          }

          setOpen(false);
          closeForm();

          toast({
            description: "Todo successfully updated!",
          });
        })
        .catch((err) => {
          setIsLoading(false);
          throw new Error(err);
        });
    }
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

  const renderDeleteFlow = () => {
    if (needConfirmation) {
      return (
        <div className="flex flex-row gap-0.5">
          <Button type="button" variant="outline" onClick={stopDelete}>
            No
          </Button>
          <Button type="button" variant="destructive" onClick={confirmDelete}>
            Yes, delete todo.
          </Button>
        </div>
      );
    } else {
      return (
        <div className="flex flex-row gap-2">
          <Button type="button" variant="destructive" onClick={deleteTodo}>
            Delete Todo
          </Button>
          <Spinner isLoading={deleteIsLoading} />
        </div>
      );
    }
  };

  const renderSubmitButton = () => {
    if (todo.id === -1) {
    }
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col space-y-8"
        autoComplete="off"
      >
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Title</FormLabel>
              <FormControl>
                <Input placeholder="Feed the Dog" {...field} />
              </FormControl>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="date"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Date</FormLabel>
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
                <Textarea
                  className="h-60"
                  placeholder="Buy snacks..."
                  {...field}
                />
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

        <div className="flex flex-row-reverse justify-between">
          <div className="flex flex-row gap-2">
            <Spinner isLoading={isLoading} />
            <Button type="submit" className="">
              Submit
            </Button>
          </div>

          {todo.id !== -1 && renderDeleteFlow()}
        </div>
      </form>
    </Form>
  );
}
