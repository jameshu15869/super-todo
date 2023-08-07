import { pb } from "@/types/pb";
import { useEffect, useReducer, useRef, useState } from "react";
import { useFetchTodos, useTodos } from "./useTodos";
import { useCombined, useFetchCombined } from "./useCombined";
import { useFetchUsers, useUsers } from "./useUsers";
import { todoReducer } from "@/reducers/reducers";

export function useSSE() {
    const sseEventSource = useRef<EventSource | null>(null);

    const {
      todos,
      setTodos,
      error: todoError,
      isLoading: todoIsLoading
    } = useFetchTodos();

    const {
      combined,
      setCombined,
      error: combinedError,
      isLoading: combinedIsLoading
    } = useFetchCombined();

    const {
      users,
      setUsers,
      error: userError,
      isLoading: userIsLoading
    } = useFetchUsers();

    useEffect(() => {
      sseEventSource.current = new EventSource("/api/sse");
      sseEventSource.current.onmessage = (e) => HandleServerMessage(e, todos, combined, users);
      return () => {
        if (sseEventSource.current) {
          sseEventSource.current.close();
        }
      };
    }, []);

    useEffect(() => {
      if (sseEventSource.current) {
        sseEventSource.current.onmessage = (e) => HandleServerMessage(e, todos, combined, users);
      }
    }, [todos, combined, users])

    const HandleServerMessage = (e: MessageEvent<any>,
      passedTodos: pb.Todo[] | undefined,
      combined: pb.Combined[] | undefined,
      users: pb.User[] | undefined) => {
        const data = JSON.parse(e.data);
        switch (data.messageType) {
            case "add-todo":
                const addTodoContent = JSON.parse(data.content);
                const newTodo = addTodoContent.todo as pb.Todo;
                if (todos) {
                    setTodos([...todos, newTodo]);
                }
                const updatedCombines = addTodoContent.combined as pb.Combined[];
                if (combined && updatedCombines) {
                  const filteredCombines = combined.filter(combined => combined.todo_id !== newTodo.id);
                  setCombined([...filteredCombines, ...updatedCombines]);
                }
                break;
            case "update-todo":
                const content = JSON.parse(data.content);
                const updatedTodo = content.todo as pb.Todo;
                if (todos && updatedTodo) {
                    const updatedTodoList = todos.map(todo => {
                      if (todo.id === updatedTodo.id) {
                        return updatedTodo;
                      }
                      return todo;
                    })
                    setTodos(updatedTodoList);
                    const updatedCombines = content.combined as pb.Combined[];
                    if (combined && updatedCombines) {
                        const filteredCombines = combined.filter(combined => combined.todo_id !== updatedTodo.id);
                        setCombined([...filteredCombines, ...updatedCombines]);
                    }
                }
                break;
            case "delete-todo":
              const deletedTodo = JSON.parse(data.content) as pb.Todo;
              if (todos && deletedTodo) {
                setTodos(todos.filter(todo => todo.id !== deletedTodo.id));
              }
              break;
            case "add-user":
              const addUserContent = JSON.parse(data.content);
              const addedUser = addUserContent.user as pb.User;
              if (users && addedUser) {
                setUsers([...users, addedUser]);
              }
              break;
        }
    }

    if (todoError) {
        throw todoError;
    }

    if (combinedError) {
        throw combinedError;
    }

    if (userError) {
        throw userError;
    }

    if (todoIsLoading || combinedIsLoading || userIsLoading) {
        return {
            todos: [],
            combined: [],
            users: [],
            isLoading: true
        }
    }

    return {
        todos : todos,
        combined: combined,
        users: users,
        isLoading: false
    }

}
