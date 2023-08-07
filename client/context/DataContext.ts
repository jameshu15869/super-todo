import { createContext } from "react";
import { pb } from "../types/pb";


export const DataContext = createContext<pb.PbData>({
    todos: undefined,
    users: undefined,
    combined: undefined
});