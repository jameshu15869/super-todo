import { pb } from "@/types/pb";
import { Avatar, AvatarFallback } from "../ui/avatar";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "../ui/tooltip";
import { cn } from "@/lib/utils";
import UserDetail from "./UserDetail";

interface TodoAvatarRequiredProps {
  user: pb.User;
  canHover?: boolean;
  canSeeDetail?: boolean;
}

export default function UserAvatar({
  user,
  canHover = false,
  canSeeDetail = false,
}: TodoAvatarRequiredProps) {
  const renderAvatar = () => {
    return (
      <Avatar className={cn(canHover ? "hover:cursor-pointer" : "")}>
        <AvatarFallback className="unselectable">
          {user.username.charAt(0).toUpperCase()}
        </AvatarFallback>
      </Avatar>
    );
  };

  if (canHover) {
    return (
      <UserDetail user={user}>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>{renderAvatar()}</TooltipTrigger>
            <TooltipContent>
              <p>{user.username}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </UserDetail>
    );
  }
  return <>{renderAvatar()}</>;
}
