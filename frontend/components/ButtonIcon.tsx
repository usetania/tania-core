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
  <Button onClick={onClick} type={type} variant={variant}>
    <div className={`d-flex align-items-center ${textColor}`}>
      {icon}
      <span>{label}</span>
    </div>
  </Button>
);

export default ButtonIcon;
