import { Button } from "react-bootstrap";

interface iButtonIcon {
  label: string;
  icon: JSX.Element;
  type?: "button" | "submit" | "reset";
  variant: string;
  onClick: () => void;
  textColor?: string;
  isBlock?: boolean;
}

const ButtonIcon = ({
  label,
  icon,
  type = "button",
  variant,
  onClick,
  textColor = "text-dark",
  isBlock = false,
}: iButtonIcon): JSX.Element => (
  <Button
    onClick={onClick}
    type={type}
    variant={variant}
    className="text-decoration-none"
  >
    <div
      className={`d-flex align-items-center ${textColor} ${
        isBlock ? "justify-content-center" : ""
      }`}
    >
      {icon}
      <span>{label}</span>
    </div>
  </Button>
);

export default ButtonIcon;
