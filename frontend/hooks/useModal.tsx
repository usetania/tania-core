import { useState } from "react";

const useModal = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const showModal = () => {
    setModalOpen(true);
  };
  const closeModal = () => {
    setModalOpen(false);
  };
  const toggleModal = () => {
    setModalOpen(!modalOpen);
  };

  return { modalOpen, setModalOpen, showModal, closeModal, toggleModal };
};

export default useModal;
