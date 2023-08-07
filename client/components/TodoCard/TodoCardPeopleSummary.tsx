import { pb } from "@/types/pb";
import UserAvatar from "../UserDetail/UserAvatar";
import { cn } from "@/lib/utils";

interface TodoCardPeopleSummaryProps {
  users: pb.User[];
}

export default function TodoCardPeopleSummary({
  users,
  className,
  ...props
}: TodoCardPeopleSummaryProps & React.HTMLAttributes<HTMLElement>) {
  return (
    <div className={cn("flex flex-row", className)}>
      {users.map((user, index) => (
        <UserAvatar key={index} user={user} canHover={true} />
      ))}
    </div>
  );
}
