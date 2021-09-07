import { Button } from "react-bootstrap";

interface iButtonIcon {
  label: string;
  icon: JSX.Element;
  type?: "button" | "submit" | "reset";
  variant: string;
  onClick: () => void;
  textColor?: string;
}

const ButtonIcon = ({
  label,
  icon,
  type = "button",
  variant,
  onClick,
  textColor = "text-dark",
}: iButtonIcon): JSX.Element => (
  <Button
    onClick={onClick}
    className={`d-flex align-items-center ${textColor}`}
    type={type}
    variant={variant}
  >
    {icon}
    <span>{label}</span>
  </Button>
);

export default ButtonIcon;
