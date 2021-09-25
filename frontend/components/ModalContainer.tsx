import { SetStateAction } from "react";
import { Modal } from "react-bootstrap";
import { FaTimes, FaCheck } from "react-icons/fa";

import ButtonIcon from "./ButtonIcon";

interface iModalContainer {
  title: string;
  isShow: boolean;
  handleCloseModal: () => void;
  handleSubmitModal: () => void;
  children: JSX.Element;
}

const ModalContainer = ({
  title,
  isShow,
  handleCloseModal,
  handleSubmitModal,
  children,
}: iModalContainer): JSX.Element => {
  return (
    <Modal show={isShow} onHide={handleCloseModal}>
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>{children}</Modal.Body>
      <Modal.Footer>
        <ButtonIcon
          label="Cancel"
          icon={<FaTimes className="me-1" />}
          variant="light"
          onClick={handleCloseModal}
        />
        <ButtonIcon
          label="Save"
          icon={<FaCheck className="me-1" />}
          variant="secondary"
          onClick={handleSubmitModal}
          textColor="text-light"
        />
      </Modal.Footer>
    </Modal>
  );
};

export default ModalContainer;
