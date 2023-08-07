import { User2, X } from "lucide-react";
import { Avatar, AvatarFallback } from "../ui/avatar";
import { Button } from "../ui/button";
import { useState } from "react";
import { cn } from "@/lib/utils";
import { pb } from "@/types/pb";
import UserAvatar from "../UserDetail/UserAvatar";

interface TodoPersonProps {
  user: pb.User;
  canEdit: boolean;
  handleDelete?: (user_id: number) => void;
}

export default function TodoPerson({
  user,
  canEdit,
  handleDelete,
  ...props
}: TodoPersonProps) {
  const [isHovered, setIsHovered] = useState(false);

  const handleDeleteMouseEnter = () => setIsHovered(true);
  const handleDeleteMouseLeave = () => setIsHovered(false);

  const handleDeleteClick = () => {
    if (handleDelete) {
      handleDelete(user.id);
    }
  };

  return (
    <div className="flex items-center space-x-2 min-w-full relative">
      {/* edit means cannot hover, no edit means can hover*/}
      <UserAvatar user={user} canHover={!canEdit} />
      <p className="text-sm">{user.username}</p>
      {canEdit && (
        <Button
          variant="ghost"
          type="button"
          className={cn(
            "absolute right-0 p-1.5 w-min h-min",
            isHovered
              ? "hover:border-red-400 hover:border-1 hover:bg-transparent"
              : ""
          )}
          onMouseEnter={handleDeleteMouseEnter}
          onMouseLeave={handleDeleteMouseLeave}
          onClick={handleDeleteClick}
        >
          <X size={20} color={isHovered ? "red" : "black"} />
        </Button>
      )}
    </div>
  );
}
