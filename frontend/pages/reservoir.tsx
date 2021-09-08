import { useState } from "react";
import type { NextPage } from "next";
import Link from "next/link";
import { Col, Form, Modal, Row, Table } from "react-bootstrap";
import { FaEdit, FaPlus, FaTimes, FaCheck } from "react-icons/fa";

import ButtonIcon from "../components/ButtonIcon";
import Layout from "../components/Layout";
import { reservoirData } from "../data";

const Reservoir: NextPage = () => {
  const [reservoirName, setReservoirName] = useState("");
  const [selectedSource, setSelectedSource] = useState();
  const [sourceNumber, setSourceNumber] = useState("0");
  const [nameIsEmpty, setNameIsEmpty] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const closeModal = () => setShowModal(false);
  const addReservoir = () => {
    if (!reservoirName) {
      setNameIsEmpty(true);
    } else {
      setNameIsEmpty(false);
      setShowModal(false);
    }
  };
  return (
    <Layout>
      <Row>
        <Col>
          <h3 className="pb-3">Water Reservoir</h3>
        </Col>
      </Row>
      <Row>
        <Col>
          <ButtonIcon
            label="Add Reservoir"
            icon={<FaPlus className="me-2" />}
            onClick={() => setShowModal(true)}
            variant="primary"
          />
        </Col>
      </Row>
      <Table responsive className="my-4">
        <thead>
          <tr>
            <th>Name</th>
            <th>Created On</th>
            <th>Source Type</th>
            <th>Capacity</th>
            <th>Used In</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {reservoirData &&
            reservoirData.map(
              ({ id, name, createdOn, sourceType, capacity, usedIn }) => (
                <tr key={id}>
                  <td>
                    <Link href={`/reservoir/${id}`}>{name}</Link>
                  </td>
                  <td>{createdOn}</td>
                  <td>{sourceType}</td>
                  <td>{capacity}</td>
                  <td>{usedIn}</td>
                  <td>
                    <FaEdit
                      onClick={() => setShowModal(true)}
                      className="show-pointer"
                    />
                  </td>
                </tr>
              )
            )}
        </tbody>
      </Table>
      <Modal show={showModal} onHide={closeModal}>
        <Modal.Header closeButton>
          <Modal.Title>Add New Reservoir</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <small className="text-muted">
            Reservoir is a water source for your farm. It can be a direct line
            from well, or water tank that has volume/capacity.
          </small>
          <Form className="mt-3">
            <Form.Group className="mb-3">
              <Form.Label>Reservoir Name</Form.Label>
              <Form.Control
                type="text"
                onChange={(e) => setReservoirName(e.target.value)}
              />
              {nameIsEmpty && (
                <Form.Text className="text-danger">
                  The name field is required
                </Form.Text>
              )}
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Source</Form.Label>
              <Form.Select onChange={(e) => setSelectedSource(e.target.value)}>
                <option>Please select source</option>
                <option value="1">Tap/Well</option>
                <option value="2">Water Tank/Cistern</option>
              </Form.Select>
            </Form.Group>
            {selectedSource === "2" && (
              <Form.Group className="mb-3">
                <Form.Control
                  type="text"
                  onChange={(e) => setSourceNumber(e.target.value)}
                  value={sourceNumber}
                />
              </Form.Group>
            )}
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <ButtonIcon
            label="Cancel"
            icon={<FaTimes className="me-1" />}
            variant="light"
            onClick={closeModal}
          />
          <ButtonIcon
            label="Save"
            icon={<FaCheck className="me-1" />}
            variant="secondary"
            onClick={addReservoir}
            textColor="text-light"
          />
        </Modal.Footer>
      </Modal>
    </Layout>
  );
};

export default Reservoir;
